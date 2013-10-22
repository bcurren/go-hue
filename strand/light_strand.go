// Package with a LightStrand that implementes the hue.API interface. It maps
// light ids to socket locations on a strand of lights.
package strand

import (
	"errors"
	"fmt"
	"github.com/bcurren/go-hue"
	"strconv"
	"strings"
)

// Structure that holds the mapping from socket id to light id. This implements
// the hue.API interface so it can be used as a drop in replacement.
type LightStrand struct {
	API    hue.API
	Length int
	Lights *TwoWayMap
}

// Create a new light strand with the given length and hue.API to delegate to.
func NewLightStrand(length int, api hue.API) *LightStrand {
	var lightStrand LightStrand
	lightStrand.API = api
	lightStrand.Length = length
	lightStrand.Lights = NewTwoWayMap()

	return &lightStrand
}

// Create a new light strand with the given length and hue.API to delegate to.
func NewLightStrandWithMap(length int, api hue.API, initMap map[string]string) *LightStrand {
	var lightStrand LightStrand
	lightStrand.API = api
	lightStrand.Length = length
	lightStrand.Lights = LoadTwoWayMap(initMap)

	return &lightStrand
}

func (lg *LightStrand) IsMappedSocketId(socketId string) bool {
	if lg.Lights.GetValue(socketId) != "" {
		return true
	}
	return false
}

func (lg *LightStrand) ChangeLength(newLength int) {
	oldLength := lg.Length
	lg.Length = newLength

	for i := newLength + 1; i <= oldLength; i++ {
		socketId := strconv.Itoa(i)
		lg.Lights.DeleteWithKey(socketId)
	}
}

func (lg *LightStrand) CleanInvalidLightIds() error {
	lights, err := lg.API.GetLights()
	if err != nil {
		return err
	}

	lg.cleanInvalidMappedLightIds(lights)
	return nil
}

// An interactive way of mapping all unmapped light bulbs on the hue bridge. This
// function does the following:
//
// 1. Turn all lights normal state
// 2. For each unmapped light
//   a. Turn the bulb selectedState
//   b. Call socketToLightFunc - The implementation should return the socket id for
//      the unmapped light. If 'x' returned, skip mapping the bulb.
//   c. Map the bulb to the socket id
//   d. Turn the buld normal state and continue
//
// This should be used to interactively prompt a person to map a light to a position
// in the strand.
func (lg *LightStrand) MapUnmappedLights(normalState, selectedState *hue.LightState,
	socketToLightFunc func(string) string) error {

	unmappedLightIds, err := lg.GetUnmappedLightIds()
	if err != nil {
		return err
	}

	for _, unmappedLightId := range unmappedLightIds {
		// Turn new unmapped light selected state
		err = lg.API.SetLightState(unmappedLightId, selectedState)
		if err != nil {
			return err
		}

		// Update the map. Skip socket ids 'x'.
		socketId := socketToLightFunc(unmappedLightId)
		if "X" != strings.ToUpper(socketId) || socketId == "" {
			if !lg.validSocketId(socketId) {
				return errors.New(fmt.Sprintf("Invalid socket id provided '%s'.", socketId))
			}
			lg.Lights.Set(socketId, unmappedLightId)
		}

		// Turn newly mapped light to normal color
		err = lg.API.SetLightState(unmappedLightId, normalState)
		if err != nil {
			return err
		}
	}

	err = lg.API.SetGroupState(hue.AllLightsGroupId, normalState)
	if err != nil {
		return err
	}

	return nil
}

// Get a list of the unmapped light ids.
func (lg *LightStrand) GetUnmappedLightIds() ([]string, error) {
	lights, err := lg.API.GetLights()
	if err != nil {
		return nil, err
	}

	lg.cleanInvalidMappedLightIds(lights)
	unmappedLightIds := lg.getUnmappedLightIds(lights)

	return unmappedLightIds, nil
}

func (lg *LightStrand) GetMap() map[string]string {
	return lg.Lights.Normal
}

func (lg *LightStrand) cleanInvalidMappedLightIds(allHueLights []hue.Light) {
	allMappedLightIds := lg.Lights.GetValues()

	for _, mappedLightId := range allMappedLightIds {
		invalidLightId := true
		for _, hueLight := range allHueLights {
			if hueLight.Id == mappedLightId {
				invalidLightId = false
				break
			}
		}
		if invalidLightId {
			lg.Lights.DeleteWithValue(mappedLightId)
		}
	}
}

func (lg *LightStrand) getUnmappedLightIds(allHueLights []hue.Light) []string {
	allMappedLightIds := lg.Lights.GetValues()

	unmappedLights := make([]string, 0, 5)
	for _, hueLight := range allHueLights {
		alreadyMapped := false
		for _, mappedLightId := range allMappedLightIds {
			if hueLight.Id == mappedLightId {
				alreadyMapped = true
				break
			}
		}
		if !alreadyMapped {
			unmappedLights = append(unmappedLights, hueLight.Id)
		}
	}

	return unmappedLights
}

func (lg *LightStrand) validSocketId(socketId string) bool {
	socketIdAsInt, err := strconv.Atoi(socketId)
	if err != nil {
		return false
	}

	if socketIdAsInt <= 0 || socketIdAsInt > lg.Length {
		return false
	}

	return true
}
