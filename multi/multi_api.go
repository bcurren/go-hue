package multi

import (
	"github.com/bcurren/go-hue"
)

type MultiAPI struct {
	apis []hue.API
}

func NewMultiAPI() *MultiAPI {
	multi := &MultiAPI{}
	multi.apis = make([]hue.API, 0, 2)
	return multi
}

func (m *MultiAPI) AddAPI(api hue.API) {
	m.apis = append(m.apis, api)
}
