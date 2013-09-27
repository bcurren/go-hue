# go-hue/strand

The strand package allows you to map hue light ids to a physical socket position on a large strand of lights. You are then able to work the the hue.API via socket position rather than lightID.

Please see [godoc.org](http://godoc.org/github.com/bcurren/go-hue/strand) for detailed API
description.

## Example

See the [examples](example/) directory for a full program using this package. Below is an inline exerpt from that directory.

```Go
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
```
