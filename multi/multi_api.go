package multi

import (
	"github.com/bcurren/go-hue"
)

type MultiAPI struct {
	APIs []hue.API
}

func NewMultiAPI() *MultiAPI {
	multi := &MultiAPI{}
	multi.APIs = make([]hue.API, 0, 2)
	return multi
}

func (m *MultiAPI) AddAPI(api hue.API) {
	m.APIs = append(m.APIs, api)
}
