package explainable

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

type TestNested struct {
	FieldA string `json:"fieldA" explain:"Explanation for field A"`
	FieldB int    `json:"fieldB" explain:"Explanation for field B"`
}

type TestStruct struct {
	SimpleField string       `json:"simpleField" explain:"A simple string field"`
	NumberField int          `json:"numberField" explain:"A numeric field"`
	Nested      TestNested   `json:"nested" explain:"Nested struct"`
	Pointer     *TestNested  `json:"pointer" explain:"Pointer to nested"`
	Slice       []TestNested `json:"slice" explain:"Slice of nested structs"`
}

// Function to load JSON data from a file into a Go structure
func loadTestData(filename string, target interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func TestExplanation(t *testing.T) {
	// Load the test input data from the testdata/input.json
	var input TestStruct
	err := loadTestData("testData/input.json", &input)
	if err != nil {
		t.Fatalf("Failed to load input data: %v", err)
	}

	// Load the expected output data from the testdata/expected.json
	var expected map[string]interface{}
	err = loadTestData("testData/expected.json", &expected)
	if err != nil {
		t.Fatalf("Failed to load expected data: %v", err)
	}

	// Call the function
	result := Explain(input)

	// Compare the result with the expected structure using reflect.DeepEqual
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result:\n%v\n but got:\n%v", expected, result)
	}

	// Check that unexpected keys do not exist in the result
	for key := range result {
		if _, ok := expected[key]; !ok {
			t.Errorf("Unexpected key in result: %q with value %q", key, result[key])
		}
	}
}

func TestRespondWithoutExplain(t *testing.T) {
	type Sample struct {
		Name  string `json:"name" explain:"The name"`
		Email string `json:"email" explain:"The email address"`
	}
	data := Sample{Name: "John", Email: "john@example.com"}

	req := httptest.NewRequest("GET", "/", nil) // No "explain=true" query param here
	w := httptest.NewRecorder()

	Respond(w, req, data)

	resp := w.Result()
	defer resp.Body.Close()

	// Verify content-type
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected content-type application/json, got %s", contentType)
	}

	// Decode the response body directly into the struct
	var parsed Sample
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	// Check if the returned data matches the input data
	if parsed.Name != data.Name || parsed.Email != data.Email {
		t.Errorf("unexpected data in response: got %+v", parsed)
	}
}

func TestRespondWithExplain(t *testing.T) {
	type User struct {
		Name  string `json:"name" explain:"The name of the user"`
		Email string `json:"email" explain:"The email address"`
	}
	user := User{Name: "Alice", Email: "alice@example.com"}

	req := httptest.NewRequest("GET", "/?explain=true", nil)
	w := httptest.NewRecorder()

	Respond(w, req, user)

	resp := w.Result()
	defer resp.Body.Close()

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'data' field in response")
	}
	nameField := data["name"].(map[string]interface{})
	emailField := data["email"].(map[string]interface{})

	if nameField["description"] != "The name of the user" {
		t.Errorf("expected description for 'name', got %v", nameField["description"])
	}
	if emailField["description"] != "The email address" {
		t.Errorf("expected description for 'email', got %v", emailField["description"])
	}
}
