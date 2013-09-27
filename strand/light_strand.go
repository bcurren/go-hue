package strand

import (
	"github.com/bcurren/go-hue"
	"strconv"
)

type LightStrand struct {
	api    hue.API
	Length int
	Lights map[string]string
}

func NewLightStrand(length int, api hue.API) *LightStrand {
	var lightStrand LightStrand
	lightStrand.api = api
	lightStrand.Length = length
	lightStrand.Lights = make(map[string]string)

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
		lg.mapLightToSocket(unmappedLightId, socketId)

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
	lightIds := make([]string, 0, len(lg.Lights))
	for _, lightId := range lg.Lights {
		lightIds = append(lightIds, lightId)
	}
	return lightIds
}

func (lg *LightStrand) mapLightToSocket(lightId, socketId string) {
	if !lg.validSocketId(socketId) {
		panic("Invalid socket id.")
	}
	lg.Lights[socketId] = lightId
}

func (lg *LightStrand) getLightIdFromSocketId(socketId string) string {
	if !lg.validSocketId(socketId) {
		panic("Invalid socket id.")
	}
	return lg.Lights[socketId]
}

func (lg *LightStrand) mapHueLightIdToSocketId(hueLights []hue.Light) []hue.Light {
	if hueLights == nil {
		return nil
	}

	lightsWithSocketId := make([]hue.Light, 0, len(hueLights))

	for _, light := range hueLights {
		socketLightId := lg.Lights[light.Id]

		// Skip any lights that haven't been registered with strand
		if socketLightId != "" {
			light.Id = socketLightId
			lightsWithSocketId = append(lightsWithSocketId, light)
		}
	}

	return lightsWithSocketId
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
