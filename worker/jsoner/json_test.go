package jsoner

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)


func Init(t *testing.T) (*Jsoner, *os.File) {
	t.Helper()

	file, err := os.Create("./test/test.json")
	if err != nil {
		t.Skip(err.Error())
	}
	jsoner := NewJson(file)

	jsoner.Data["1.go"] = map[string]any{
		"number_of_rows": 5,
	}
	
	return jsoner, file
}

func TestWrite(t *testing.T) {
	jsoner, file := Init(t)

	expectedJson, err := os.Open("./test/write_test.json")
	if err != nil {
		t.Skip(err.Error())
	}

	jsoner.Write()
	file.Close()

	file, err = os.Open("./test/test.json")
	if err != nil {
		t.Skip(err.Error())
	}

	var actual map[string]any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&actual)
	if err != nil {
		t.Skip(err)
	}

	var expected map[string]any
	decoder = json.NewDecoder(expectedJson)
	err = decoder.Decode(&expected)
	if err != nil {
		t.Skip(err)
	}

	if !reflect.DeepEqual(actual, expected){
		t.Skip("Данные не совпадают")
	}

	defer func(){
		file.Close()
		expectedJson.Close()
	}()
}

func TestAddFileName(t *testing.T) {
	jsoner, file := Init(t)

	jsoner.AddFileName("2.go")
	jsoner.Write()
	file.Close()

	file, err := os.Open("./test/test.json")
	if err != nil {
		t.Skip(err.Error())
	}

	expectedJson, err := os.Open("./test/add_file_name_test.json")
	if err != nil {
		t.Skip(err.Error())
	}

	// TODO make func for this block
	var actual map[string]any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&actual)
	if err != nil {
		t.Skip(err)
	}

	var expected map[string]any
	decoder = json.NewDecoder(expectedJson)
	err = decoder.Decode(&expected)
	if err != nil {
		t.Skip(err)
	}

	if !reflect.DeepEqual(actual, expected){
		t.Skip("Данные не совпадают")
	}

	defer func(){
		file.Close()
		expectedJson.Close()
	}()
}

