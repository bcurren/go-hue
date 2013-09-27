package strand

type TwoWayMap struct {
	Normal   map[string]string
	Inverted map[string]string
}

func NewTwoWayMap() *TwoWayMap {
	twoWayMap := &TwoWayMap{}
	twoWayMap.Normal = make(map[string]string)
	twoWayMap.Inverted = make(map[string]string)
	return twoWayMap
}

func LoadTwoWayMap(initMap map[string]string) *TwoWayMap {
	twoWayMap := NewTwoWayMap()
	for k, v := range initMap {
		twoWayMap.Set(k, v)
	}
	return twoWayMap
}

func (m *TwoWayMap) Set(key, value string) {
	if oldValue := m.Normal[key]; oldValue != "" {
		delete(m.Inverted, oldValue)
	}
	m.Normal[key] = value

	if oldKey := m.Inverted[value]; oldKey != "" {
		delete(m.Normal, oldKey)
	}
	m.Inverted[value] = key
}

func (m *TwoWayMap) GetKey(value string) string {
	return m.Inverted[value]
}

func (m *TwoWayMap) GetValue(key string) string {
	return m.Normal[key]
}

func (m *TwoWayMap) GetKeys() []string {
	keys := make([]string, 0, len(m.Normal))
	for k, _ := range m.Normal {
		keys = append(keys, k)
	}
	return keys
}

func (m *TwoWayMap) GetValues() []string {
	values := make([]string, 0, len(m.Normal))
	for _, v := range m.Normal {
		values = append(values, v)
	}
	return values
}

func (m *TwoWayMap) Length() int {
	return len(m.Normal)
}
