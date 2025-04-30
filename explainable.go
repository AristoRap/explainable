package explainable

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func Explain(data any) map[string]any {
	explanations := map[string]any{}
	val := reflect.ValueOf(data)
	defer func() {
		if r := recover(); r != nil {
			log.Println("explainable: recovered from panic:", r)
			explanations = map[string]any{}
		}
	}()

	// Dereference pointer if necessary
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle top-level slice
	if val.Kind() == reflect.Slice {
		if val.Len() > 0 {
			// Explain first element only as a representative sample
			sliceExplanation := map[string]any{
				"description": "List of results",
				"results":     []any{Explain(val.Index(0).Interface())},
			}
			return sliceExplanation
		}
		return explanations
	}

	// Handle struct
	if val.Kind() == reflect.Struct {
		typ := val.Type()

		for i := range val.NumField() {
			field := typ.Field(i)
			fieldVal := val.Field(i)

			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			// Get the 'json' and 'explain' tags
			jsonTag := field.Tag.Get("json")
			explainTag := field.Tag.Get("explain")

			// Skip fields without a valid 'json' tag
			if jsonTag == "" || jsonTag == "-" {
				continue
			}

			// Handle pointers
			actualVal := fieldVal
			fieldKind := field.Type.Kind()
			if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
				actualVal = fieldVal.Elem()
				fieldKind = actualVal.Kind()
			}

			switch fieldKind {
			case reflect.Struct:
				// Nested struct: Recursively describe the structure inside
				nestedExplanation := Explain(actualVal.Interface())
				explanations[jsonTag] = map[string]any{
					"description": explainTag,
					"fields":      nestedExplanation,
				}

			case reflect.Slice:
				// Slice: Handle each element of the slice
				elemType := field.Type.Elem()
				if elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}

				// Prepare slice explanation structure
				sliceExplanation := map[string]any{
					"description": explainTag,
					"items":       []any{},
				}

				// Handle slice of structs
				if elemType.Kind() == reflect.Struct {
					// Describe just one representative element
					zeroElem := reflect.New(elemType).Elem().Interface()
					elemExplanation := Explain(zeroElem)
					sliceExplanation["items"] = []any{elemExplanation}
				} else {
					// Non-struct: describe its type
					sliceExplanation["items"] = []any{
						map[string]any{
							"type":        elemType.String(),
							"description": explainTag,
						},
					}
				}

				explanations[jsonTag] = sliceExplanation

			default:
				// Handle primitive types
				explanations[jsonTag] = map[string]any{
					"type":        field.Type.String(),
					"description": explainTag,
				}
			}
		}
	}

	return explanations
}

func Respond(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Get("explain") == "true" {
		json.NewEncoder(w).Encode(map[string]any{"data": Explain(data)})
	} else {
		json.NewEncoder(w).Encode(data)
	}
}
