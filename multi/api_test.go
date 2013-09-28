package multi

import (
	"errors"
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
	"testing"
)

func Test_GetLights(t *testing.T) {
	multi, api1, api2 := newStubMultiAPI()
	api1.GetLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}
	api2.GetLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}

	lights, err := multi.GetLights()
	if err != nil {
		t.Fatal("Error was returned.", err)
	}

	if len(lights) != 2 {
		t.Error("Should merge GetLights and have 2 results.")
	}
}

func Test_GetNewLights(t *testing.T) {
	multi, api1, api2 := newStubMultiAPI()
	api1.GetNewLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}
	api2.GetNewLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}

	lights, _, err := multi.GetNewLights()
	if err != nil {
		t.Fatal("Error was returned.", err)
	}

	if len(lights) != 2 {
		t.Error("Should merge GetLights and have 2 results.")
	}
}

func Test_SearchForNewLights(t *testing.T) {
	multi, _, _ := newStubMultiAPI()

	err := multi.SearchForNewLights()
	if err != nil {
		t.Fatal("Error was returned.", err)
	}
}

func newStubMultiAPI() (*MultiAPI, *huetest.StubAPI, *huetest.StubAPI) {
	api1 := &huetest.StubAPI{}
	api2 := &huetest.StubAPI{}

	multi := NewMultiAPI()
	multi.AddAPI(api1)
	multi.AddAPI(api2)

	return multi, api1, api2
}

func Test_mergeErrors(t *testing.T) {
	myErrors := make([]error, 0, 3)

	// Handle nil and empty cases
	if mergeErrors(nil) != nil {
		t.Error("A nil slice of errors should return nil.")
	}
	if mergeErrors(myErrors) != nil {
		t.Error("An empty slice of errors should return nil.")
	}

	// Merge hue.APIErrors together
	apiError := createResourceNotAvailableAPIError("light1", "/lights/light1")
	myErrors = append(myErrors, apiError)
	myErrors = append(myErrors, apiError)
	mergedError := mergeErrors(myErrors)
	mergedAPIError, ok := mergedError.(*hue.APIError)
	if !ok {
		t.Fatalf("Should have return an APIError but return %T.", mergedError)
	}
	if len(mergedAPIError.Errors) != 2 {
		t.Error("Should be one api error detail.")
	}
	if mergedAPIError.Errors[0].Type != hue.ResourceNotAvailableErrorType {
		t.Error("First should be no resource.")
	}
	if mergedAPIError.Errors[1].Type != hue.ResourceNotAvailableErrorType {
		t.Error("Second should be no resource.")
	}

	// Return first non API error if one exists
	nonAPIError := errors.New("Some error.")
	myErrors = append(myErrors, nonAPIError)
	if nonAPIError != mergeErrors(myErrors) {
		t.Error("Non API Error should be returned of APIError.")
	}
}

func Test_MultiBridgeImplementsAPIInterface(t *testing.T) {
	multi := NewMultiAPI()
	funcThatTakesAPIAsParameter(multi)
}

func funcThatTakesAPIAsParameter(api hue.API) {
	// noop
}
