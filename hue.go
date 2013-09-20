package hue

//
// func FindBridges() ([]Bridge, error) {
// }
//
//
//
// bridges, err := hue.FindBridges()
// if err != nil {
// 	return
// }
// bridge = bridges[0]
// user, err := bridge.CreateUser("devicetype")
// if err != nil {
// 	return
// }
// lights, err := user.GetLights()
// if err != nil {
// 	return
// }
//
// type jsonConn interface {
// 	Send(method string, map[string]interface{}) (map[string]interface{} error)
// }
//
// // Configuration
// type Bridge struct {
// 	BridgeAddr *net.TCPAddr
// 	api network
// }
// type User struct {
// 	Bridge Bridge
// 	Username string
// }
//
// type Configuration struct {
// 	proxyport uint16
// 	utc string
// 	name string
// 	swupdate object
// 	whitelist object
// 	swversion string
// 	proxyaddress string
// 	mac string
// 	linkbutton bool
// 	ipaddress string
// 	netmask string
// 	gateway string
// 	dhcp bool
// 	portalservices bool
// }
// type DataStore struct {
// 	Lights []Light
// 	Groups []Group
// 	Schedules []Schedule
// 	Config Configuration
// }
// type LightState struct {
// 	On bool
// 	Brightness uint8
// 	Hue uint16
// 	Saturation uint8
// 	xy [2]float
// 	ColorTemp uint16
// 	Alert string
// 	Effect string
// 	TransitionTime uint16 /* write only */
// 	ColorMode string /* read only */
// 	Reachable bool /* read only */
// }
// type Light struct {
// 	Id string
// 	Name string
// }
// type LightAttributes struct {
// 	Id string
// 	Name string
// 	State LightState
// 	Type string
// 	ModelId string
// 	SoftwareVersion string
// 	Pointsymbol object // reserved for future use
// }
//
// CreateUser(deviceType, username *optional*) (*User, error)
// (u *User) GetConfiguration() (*Configuration, error)
// (u *User) UpdateConfiguration(*Configuration) error
// (u *User) DeleteUser(*User) error
// (u *User) GetDataStore() (DataStore, error)
//
// // Lights
// (u *User) GetLights() ([]Light, error)
// (u *User) GetNewLights() ([]Light, lastScan time.Time, error)
// (u *User) StartSearchForNewLights() (error)
// (u *User) GetLightAttributes(lightId string) (LightAttributes, error)
// (u *User) SetLightName(lightId string, name string) error
// (u *User) SetLightState(lightId string, state LightState) []error
//
