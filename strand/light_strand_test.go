package strand

import (
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
	"strings"
	"testing"
)

func Test_NewLightStrand(t *testing.T) {
	stubHueAPI := &huetest.StubAPI{}
	lightStrand := NewLightStrand(3, stubHueAPI)

	if lightStrand.Length != 3 {
		t.Error("Light strand should have length given in constructor.")
	}
	if len(lightStrand.Lights) != 0 {
		t.Error("Light strand should have Lights initialized to empty.")
	}
}

func Test_MapUnmappedLights(t *testing.T) {
	hueLights := make([]hue.Light, 1, 1)
	hueLights[0].Id = "3"

	stubHueAPI := &huetest.StubAPI{}
	stubHueAPI.GetLightsError = nil
	stubHueAPI.GetLightsReturn = hueLights

	lightStrand := NewLightStrand(3, stubHueAPI)

	countTimesCalled := 0
	err := lightStrand.MapUnmappedLights(func() string {
		countTimesCalled += 1

		// Should have set hue light to red
		if stubHueAPI.SetLightStateParamLightId != "3" {
			t.Error("Should have set light 3.")
		}
		if stubHueAPI.SetLightStateParamLightState.Hue == nil {
			t.Error("Should have set the state.")
		}

		return "1"
	})
	if err != nil {
		t.Fatal("Error returned when mapping")
	}

	if countTimesCalled != 1 {
		t.Error("Map function called more than 1 time.")
	}
	if lightStrand.Lights["1"] != "3" {
		t.Error("Didn't map to the correct id.")
	}

	// Should have set hue light to white
	if stubHueAPI.SetLightStateParamLightId != "3" {
		t.Error("Should have set light 3.")
	}
	if stubHueAPI.SetLightStateParamLightState.ColorTemp == nil {
		t.Error("Should have set the state.")
	}
}

func Test_GetMappedLightIds(t *testing.T) {
	lightStrand := NewLightStrand(3, nil)
	lightStrand.mapLightToSocket("3", "1")
	lightStrand.mapLightToSocket("2", "2")
	lightStrand.mapLightToSocket("1", "3")

	expected := []string{"3", "2", "1"}
	actual := lightStrand.GetMappedLightIds()
	if !stringSlicesEqual(expected, actual) {
		t.Errorf("Expected a slice of all mapped light ids. Expected %v but received %v.\n", expected, actual)
	}
}

func Test_GetUnmappedLightIds(t *testing.T) {
	hueLights := make([]hue.Light, 4, 4)
	hueLights[0].Id = "3"
	hueLights[1].Id = "1"
	hueLights[2].Id = "5"
	hueLights[3].Id = "2"

	stubHueAPI := &huetest.StubAPI{}
	stubHueAPI.GetLightsError = nil
	stubHueAPI.GetLightsReturn = hueLights

	lightStrand := NewLightStrand(3, stubHueAPI)
	lightStrand.mapLightToSocket("3", "1")
	lightStrand.mapLightToSocket("2", "2")
	lightStrand.mapLightToSocket("1", "3")

	expected := []string{"5"}
	actual, _ := lightStrand.GetUnmappedLightIds()
	if !stringSlicesEqual(expected, actual) {
		t.Errorf("Expected a slice of all unmapped light ids. Expected %v but received %v.\n", expected, actual)
	}
}

// A function to test if two string slices are equal. If you know a better way, please
// update this function. Seems like there should be a beeter way but I couldn't find one.
func stringSlicesEqual(slice1 []string, slice2 []string) bool {
	// Check both nil or both not nil
	if slice1 == nil && slice2 == nil {
		return true
	} else if slice1 == nil && slice2 != nil {
		return false
	} else if slice1 != nil && slice2 == nil {
		return false
	}

	// Length must be the same
	if len(slice1) != len(slice2) {
		return false
	}

	// Contents must be the same
	sep := "|||"
	slices1String := strings.Join(slice1, sep)
	slices2String := strings.Join(slice2, sep)
	if slices1String != slices2String {
		return false
	}

	return true
}

func Test_validSocketId(t *testing.T) {
	lightStrand := NewLightStrand(3, nil)
	if lightStrand.validSocketId("0") {
		t.Error("Socket id 0 should be invalid.")
	}
	if !lightStrand.validSocketId("1") {
		t.Error("Socket id 1 should be valid.")
	}
	if !lightStrand.validSocketId("3") {
		t.Error("Socket id 3 should be valid.")
	}
	if lightStrand.validSocketId("4") {
		t.Error("Socket id 4 should be invalid.")
	}
	if lightStrand.validSocketId("notint") {
		t.Error("Socket id that is not an int should be invalid.")
	}
}
