package strand

import (
	"testing"
)

func Test_LoadTwoWayMap(t *testing.T) {
	initMap := map[string]string{
		"key1": "val1",
		"key2": "val1",
		"key3": "val3",
	}
	tMap := LoadTwoWayMap(initMap)

	if "key2" != tMap.GetKey("val1") {
		t.Error("Should find new key for duplicate value.")
	}
	if "" != tMap.GetKey("val2") {
		t.Error("Should not find key for invalid value.")
	}
	if "key3" != tMap.GetKey("val3") {
		t.Error("Should find key for value.")
	}
	if "" != tMap.GetValue("key1") {
		t.Error("Should not find old key for duplicate value.")
	}
	if "val1" != tMap.GetValue("key2") {
		t.Error("Should find new value for duplicate value.")
	}
	if "val3" != tMap.GetValue("key3") {
		t.Error("Should find value for key.")
	}
}

func Test_Set(t *testing.T) {
	tMap := NewTwoWayMap()
	tMap.Set("key1", "val1")

	if "val1" != tMap.GetValue("key1") {
		t.Error("Invalid value when calling GetValue().")
	}

	if "key1" != tMap.GetKey("val1") {
		t.Error("Invalid value when calling GetKey().")
	}
}

func Test_SetDuplicateKey(t *testing.T) {
	tMap := NewTwoWayMap()
	tMap.Set("key1", "val1")
	tMap.Set("key1", "newVal1")

	if "newVal1" != tMap.GetValue("key1") {
		t.Error("Should find new value from key.")
	}

	if "key1" != tMap.GetKey("newVal1") {
		t.Error("Should find key from new value.")
	}

	if "" != tMap.GetKey("val1") {
		t.Error("Should not find key from old value.")
	}
}

func Test_SetDuplicateValue(t *testing.T) {
	tMap := NewTwoWayMap()
	tMap.Set("key1", "val1")
	tMap.Set("key2", "val1")

	if "val1" != tMap.GetValue("key2") {
		t.Error("Should find value from new key.")
	}

	if "key2" != tMap.GetKey("val1") {
		t.Error("Should find new key from value.")
	}

	if "" != tMap.GetValue("key1") {
		t.Error("Should not find value from old key.")
	}
}

func Test_GetValueWithNoExistentKey(t *testing.T) {
	tMap := NewTwoWayMap()

	if "" != tMap.GetValue("nonexistent") {
		t.Error("Should get empty string when nonexistent key is used in GetValue().")
	}
}

func Test_GetKeyWithNoExistentKey(t *testing.T) {
	tMap := NewTwoWayMap()

	if "" != tMap.GetKey("nonexistent") {
		t.Error("Should get empty string when nonexistent key is used in GetKey().")
	}
}

func Test_GetValues(t *testing.T) {
	tMap := NewTwoWayMap()
	tMap.Set("key1", "val1")
	tMap.Set("key2", "val2")

	actualValues := tMap.GetValues()
	expectedValues := []string{"val1", "val2"}

	if !stringSlicesEqual(expectedValues, actualValues) {
		t.Errorf("GetValues() returned unexpected slice. Expected %v but received %v.\n",
			expectedValues, actualValues)
	}
}

func Test_GetKeys(t *testing.T) {
	tMap := NewTwoWayMap()
	tMap.Set("key1", "val1")
	tMap.Set("key2", "val2")

	actualKeys := tMap.GetKeys()
	expectedKeys := []string{"key1", "key2"}

	if !stringSlicesEqual(expectedKeys, actualKeys) {
		t.Errorf("GetKeys() returned unexpected slice. Expected %v but received %v.\n",
			expectedKeys, actualKeys)
	}
}
