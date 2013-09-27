package strand

import (
	"github.com/bcurren/go-hue"
	"time"
)

func (lg *LightStrand) GetLights() ([]hue.Light, error) {
	lights, err := lg.api.GetLights()
	return lg.mapHueLightIdToSocketId(lights), err
}

func (lg *LightStrand) GetNewLights() ([]hue.Light, time.Time, error) {
	lights, lastUpdated, err := lg.api.GetNewLights()
	return lg.mapHueLightIdToSocketId(lights), lastUpdated, err
}

func (lg *LightStrand) SearchForNewLights() error {
	return lg.api.SearchForNewLights()
}

func (lg *LightStrand) GetLightAttributes(socketId string) (*hue.LightAttributes, error) {
	return lg.api.GetLightAttributes(lg.getLightIdFromSocketId(socketId))
}

func (lg *LightStrand) SetLightName(socketId string, name string) error {
	return lg.api.SetLightName(lg.getLightIdFromSocketId(socketId), name)
}

func (lg *LightStrand) SetLightState(socketId string, state hue.LightState) error {
	return lg.api.SetLightState(lg.getLightIdFromSocketId(socketId), state)
}

func (lg *LightStrand) GetConfiguration() (*hue.Configuration, error) {
	return lg.api.GetConfiguration()
}
