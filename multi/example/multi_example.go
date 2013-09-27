package example

func Test_Blah(t *testing.T) {
	multi := MultiAPI(nil)
	
	bridges, err := hue.FindBridges()
	if err != nil {
		panic(err)
	}
	
	for _, bridge := range bridges {
		multi.AddAPI(bridge)
	}
	
	lights, err := multi.GetLights()
	if err != nil {
		panic(err)
	}
}
