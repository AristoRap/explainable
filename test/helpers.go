package test

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestNested struct {
	FieldA string `json:"fieldA" explain:"Explanation for field A"`
	FieldB int    `json:"fieldB" explain:"Explanation for field B"`
}
type TestNestedTwice struct {
	FieldC TestNested `json:"fieldC" explain:"Explanation for field C"`
}

type TestStructA struct {
	SimpleField string          `json:"simpleField" explain:"A simple string field"`
	NumberField int             `json:"numberField" explain:"A numeric field"`
	Nested      TestNested      `json:"nested" explain:"Nested struct"`
	NestedTwice TestNestedTwice `json:"nestedTwice" explain:"Twice nested struct"`
	Pointer     *TestNested     `json:"pointer" explain:"Pointer to nested"`
	Slice       []TestNested    `json:"slice" explain:"Slice of nested structs"`
}

type TestStructB struct {
	SimpleField string       `json:"simpleField" explain:"A simple string field"`
	NumberField int          `json:"numberField" explain:"A numeric field"`
	Nested      TestNested   `json:"nested" explain:"Nested struct"`
	Pointer     *TestNested  `json:"pointer" explain:"Pointer to nested"`
	Slice       []TestNested `json:"slice" explain:"Slice of nested structs"`
}

func LoadTestData(filename string, target any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func SetupSingle(input *TestStructA, expected *map[string]any) error {
	testInput := fmt.Sprintf("test/data/input%s.json", "SingleItem")
	testExpected := fmt.Sprintf("test/data/expected%s.json", "SingleItem")

	err := LoadTestData(testInput, &input)
	if err != nil {
		return err
	}
	err = LoadTestData(testExpected, &expected)
	if err != nil {
		return err
	}
	return nil
}

func SetupList(input *[]TestStructB, expected *map[string]any) error {
	testInput := fmt.Sprintf("test/data/input%s.json", "ListItem")
	testExpected := fmt.Sprintf("test/data/expected%s.json", "ListItem")

	err := LoadTestData(testInput, &input)
	if err != nil {
		return err
	}
	err = LoadTestData(testExpected, &expected)
	if err != nil {
		return err
	}
	return nil
}
