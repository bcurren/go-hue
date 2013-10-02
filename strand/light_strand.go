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
	if lg.Lights.GetKey(socketId) != "" {
		return true
	}
	return false
}

// An interactive way of mapping all unmapped light bulbs on the hue bridge. This
// function does the following:
//
// 1. Turn all lights white
// 2. For each unmapped light
//   a. Turn the bulb red
//   b. Call socketToLightFunc - The implementation should return the socket id for
//      the unmapped light. If 'x' returned, skip mapping the bulb.
//   c. Map the bulb to the socket id
//   d. Turn the white bulb and continue
//
// This should be used to interactively prompt a person to map a light to a position
// in the strand.
func (lg *LightStrand) MapUnmappedLights(socketToLightFunc func(string) string) error {
	lights, err := lg.API.GetLights()
	if err != nil {
		return err
	}

	lg.cleanInvalidMappedLightIds(lights)
	unmappedLightIds := lg.getUnmappedLightIds(lights)
	white := createWhiteLightState()
	red := createRedLightState()

	for _, unmappedLightId := range unmappedLightIds {
		// Turn new unmapped light red
		err = lg.API.SetLightState(unmappedLightId, red)
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

		// Turn newly mapped light white
		err = lg.API.SetLightState(unmappedLightId, white)
		if err != nil {
			return err
		}
	}

	return nil
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

func createWhiteLightState() *hue.LightState {
	white := &hue.LightState{}
	white.On = new(bool)
	*white.On = true
	white.ColorTemp = new(uint16)
	*white.ColorTemp = 1800

	return white
}

func createRedLightState() *hue.LightState {
	red := &hue.LightState{}
	red.On = new(bool)
	*red.On = true
	red.Brightness = new(uint8)
	*red.Brightness = 255
	red.Hue = new(uint16)
	*red.Hue = 65535
	red.Saturation = new(uint8)
	*red.Saturation = 255

	return red
}
