# explainable

`explainable` is a Go package that helps describe structs and their fields with metadata. It uses reflection to provide an easy-to-understand explanation of the structure of structs, slices, and other types in your Go code.

The package automatically generates a description for each field, allowing developers to understand the structure and metadata associated with the fields of their Go types.

## Installation

To install the `explainable` package, run the following Go command:

```bash
go get github.com/aristorap/explainable
```

## Usage

### Example 1: Explaining a Struct

```go
package main

import (
	"fmt"
	"github.com/aristorap/explainable"
)

type User struct {
	ID    int    `json:"id" explain:"The unique identifier"`
	Name  string `json:"name" explain:"The name of the user"`
	Email string `json:"email" explain:"The email address"`
}

func main() {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	// Get the explanation of the User struct
	explanation := explainable.Explain(user)
	fmt.Printf("%+v", explanation)  // Print the explanation
}
```

### Example 2: Explaining a Slice of Structs

```go
package main

import (
	"fmt"
	"github.com/aristorap/explainable"
)

type Product struct {
	ID   int    `json:"id" explain:"Product identifier"`
	Name string `json:"name" explain:"Product name"`
}

func main() {
	products := []Product{
		{ID: 1, Name: "Laptop"},
		{ID: 2, Name: "Phone"},
	}

	// Get the explanation of the slice of Products
	explanation := explainable.Explain(products)
	fmt.Printf("%+v", explanation)  // Print the explanation
}
```

## Features

- **Structs**: Describes the fields of the struct along with their JSON and explanation tags.
- **Slices**: Describes slices and their elements. If the elements are structs, a representative element is explained.
- **Pointers**: The package handles struct pointers, explaining their underlying values.
- **Custom Tags**: Allows custom explanations via the `explain` tag, making it easy to provide context for each field.
- **Recursion**: Nested structs are recursively explained to give a comprehensive view of the data structure.

## Tags

The package relies on the following tags:

- `json`: Used to indicate the JSON field name.
- `explain`: Used to provide a description of the field.

### Example of Tags

```go
type Example struct {
	Field1 string `json:"field1" explain:"First field description"`
	Field2 int    `json:"field2" explain:"Second field description"`
}
```

## Error Handling

The package uses Go's `reflect` package, so in case of unexpected values or types, it recovers from panics and returns an empty explanation instead of crashing your program.

## Contribution

Contributions are welcome! Feel free to open issues or submit pull requests to improve the package. Ensure that tests are included for new features or bug fixes.

## License

`explainable` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
