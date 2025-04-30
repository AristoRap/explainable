package explainable

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aristorap/explainable/test"
)

func TestExplainSingleItem(t *testing.T) {
	// Load the test input and expected data
	var input test.TestStructA
	var expected map[string]any
	err := test.SetupSingle(&input, &expected)
	if err != nil {
		t.Fatalf("Failed to setup data: %v", err)
	}

	// Call the function
	result := Explain(input)

	// Compare the result with the expected struct
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
func TestExplainListItem(t *testing.T) {
	// Load the test input and expected data
	var input []test.TestStructB
	var expected map[string]any
	err := test.SetupList(&input, &expected)
	if err != nil {
		t.Fatalf("Failed to setup data: %v", err)
	}

	// Call the function
	result := Explain(input)

	// Compare the result with the expected struct
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result:\n%v\n but got:\n%v", expected, result)
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

	var parsed map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := parsed["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'data' field in response")
	}
	nameField := data["name"].(map[string]any)
	emailField := data["email"].(map[string]any)

	if nameField["description"] != "The name of the user" {
		t.Errorf("expected description for 'name', got %v", nameField["description"])
	}
	if emailField["description"] != "The email address" {
		t.Errorf("expected description for 'email', got %v", emailField["description"])
	}
}
