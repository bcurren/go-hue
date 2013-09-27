package main

import (
	"fmt"
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/strand"
)

const HueUsername = "strandexample"
const HueDeviceType = "huestrand"

func main() {
	// Find the hue bridge
	bridges, err := hue.FindBridges()
	if err != nil {
		panic(err)
	}
	if len(bridges) < 1 {
		panic("Couldn't find a Hue bridge on your network.")
	}
	bridge := bridges[0]

	// Create user or find existing
	isValid, err := bridge.IsValidUser(HueUsername)
	if err != nil {
		panic(err)
	}

	var user *hue.User
	if isValid {
		user = hue.NewUserWithBridge(HueUsername, bridge)
	} else {
		user, err = bridge.CreateUser(HueUsername, HueDeviceType)
		if err != nil {
			panic(err)
		}
	}

	// Create a strand and map socketid to lightid
	lights, err := user.GetLights()
	if err != nil {
		panic(err)
	}

	lightStrand := strand.NewLightStrand(len(lights), user)
	lightStrand.MapUnmappedLights(func() string {
		fmt.Print("Enter the socket id of the red bulb: ")
		var socketId string
		fmt.Scanln(&socketId)
		return socketId
	})

	// Use lightStrand just like hue.API with socket id instead of light id
	fmt.Println("Here are all the lights on the strand. Notice the ids match their location.")
	lights, err = lightStrand.GetLights()
	for _, light := range lights {
		fmt.Printf("Id: %s, Name: %s\n", light.Id, light.Name)
	}
}
