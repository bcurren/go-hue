package strand

import (
	"github.com/bcurren/go-hue"
	"strconv"
)

type LightStrand struct {
	api    hue.API
	Length int
	Lights *TwoWayMap
}

func NewLightStrand(length int, api hue.API) *LightStrand {
	var lightStrand LightStrand
	lightStrand.api = api
	lightStrand.Length = length
	lightStrand.Lights = NewTwoWayMap()

	return &lightStrand
}

func (lg *LightStrand) SetDelegateAPI(api hue.API) {
	lg.api = api
}

func (lg *LightStrand) MapUnmappedLights(socketToLightFunc func() string) error {
	unmappedLightIds, err := lg.GetUnmappedLightIds()
	if err != nil {
		return err
	}

	white := createWhiteLightState()
	red := createRedLightState()

	for _, unmappedLightId := range unmappedLightIds {
		// Turn new unmapped light red
		err = lg.api.SetLightState(unmappedLightId, red)
		if err != nil {
			return err
		}

		socketId := socketToLightFunc()
		lg.setSocketIdToLightId(socketId, unmappedLightId)

		// Turn newly mapped light white
		err = lg.api.SetLightState(unmappedLightId, white)
		if err != nil {
			return err
		}
	}

	return nil
}

func (lg *LightStrand) GetUnmappedLightIds() ([]string, error) {
	allHueLights, err := lg.api.GetLights()
	if err != nil {
		return nil, err
	}

	allMappedLightIds := lg.GetMappedLightIds()

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

	return unmappedLights, nil
}

func (lg *LightStrand) GetMappedLightIds() []string {
	return lg.Lights.GetValues()
}

func (lg *LightStrand) setSocketIdToLightId(socketId, lightId string) {
	if !lg.validSocketId(socketId) {
		panic("Invalid socket id.")
	}
	lg.Lights.Set(socketId, lightId)
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

func createWhiteLightState() hue.LightState {
	white := hue.LightState{}
	white.ColorTemp = new(uint16)
	*white.ColorTemp = 1800

	return white
}

func createRedLightState() hue.LightState {
	red := hue.LightState{}
	red.Brightness = new(uint8)
	*red.Brightness = 255
	red.Hue = new(uint16)
	*red.Hue = 65535
	red.Saturation = new(uint8)
	*red.Saturation = 255

	return red
}
