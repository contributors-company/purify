# Purify

**Purify** is a Go library for struct validation based on tags. It allows you to easily validate data within structs using `purify` tags, making your code cleaner and easier to maintain.

## Table of Contents

- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Requirements](#requirements)
- [Usage](#usage)
  - [Example](#example)
  - [Available Validation Tags](#available-validation-tags)
- [Custom Validators](#custom-validators)
  - [Creating a Custom Validator](#creating-a-custom-validator)
  - [Password Validation Example](#password-validation-example)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Author](#author)
- [Acknowledgments](#acknowledgments)

## Getting Started

### Installation

Install the library using `go get`:

```bash
go get github.com/contributors-company/purify
```

Import the package into your project:

```go
import "github.com/contributors-company/purify"
```

### Requirements

- Go 1.13 or higher

## Usage

Purify allows you to validate struct fields using `purify` tags in your struct definitions.

### Example

```go
package main

import (
    "fmt"
    "github.com/contributors-company/purify"
)

type User struct {
    Email    string `json:"email" purify:"required|email"`
    Username string `json:"username" purify:"required|min(3)|max(20)"`
    Age      int    `json:"age" purify:"min(18)|max(99)"`
}

func main() {
    user := User{
        Email:    "example@example.com",
        Username: "example_user",
        Age:      25,
    }

    // Validate the struct
    err := purify.ValidateStruct(user)
    if err != nil {
        fmt.Println("Validation error:", err)
    } else {
        fmt.Println("Validation successful")
    }
}
```

In this example:

- The `Email` field must be filled and contain a valid email address.
- The `Username` field must be filled, have a length of at least 3 and at most 20 characters.
- The `Age` field must be at least 18 and at most 99.

### Available Validation Tags

- `required` — checks that the field is not empty.
- `email` — checks that the field contains a valid email address.
- `min(n)` — checks that the field's value is not less than `n`. For strings, this is the length of the string; for numbers, the value itself.
- `max(n)` — checks that the field's value is not greater than `n`. For strings, this is the length of the string; for numbers, the value itself.

**Example of tag usage:**

```go
type Product struct {
    Name  string  `json:"name" purify:"required|min(3)|max(50)"`
    Price float64 `json:"price" purify:"required|min(0.01)"`
}
```

## Custom Validators

Purify provides the ability to create your own validators to extend the library's functionality to meet your project's specific requirements.

### Creating a Custom Validator

To create a custom validator, define a function that matches the `ValidatorFunc` type and register it using `RegisterValidator`.

**`ValidatorFunc` type:**

```go
type ValidatorFunc func(fieldValue string, param string) string
```

- `fieldValue` — the value of the field to be validated.
- `param` — the parameter passed in the validation tag (if any).
- The function should return an error message string or an empty string if the validation passes.

**Example of custom validators:**

```go
func Min() ValidatorFunc {
    return func(fieldValue string, param string) string {
        minLength, _ := strconv.Atoi(param)
        if len(fieldValue) < minLength {
            return fmt.Sprintf("minimum length is %s characters", param)
        }
        return ""
    }
}

func Max() ValidatorFunc {
    return func(fieldValue string, param string) string {
        maxLength, _ := strconv.Atoi(param)
        if len(fieldValue) > maxLength {
            return fmt.Sprintf("maximum length is %s characters", param)
        }
        return ""
    }
}

func Required() ValidatorFunc {
    return func(fieldValue string, _ string) string {
        if fieldValue == "" {
            return "required field"
        }
        return ""
    }
}

func Email() ValidatorFunc {
    return func(fieldValue string, _ string) string {
        emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
        matched, _ := regexp.MatchString(emailRegex, fieldValue)
        if !matched {
            return "invalid email"
        }
        return ""
    }
}

// Register the validators
func init() {
    RegisterValidator("min", Min())
    RegisterValidator("max", Max())
    RegisterValidator("required", Required())
    RegisterValidator("email", Email())
}
```

### Password Validation Example

Suppose you want to add password validation to meet certain requirements, such as:

- Minimum length of 8 characters.
- At least one digit.
- At least one uppercase letter.

**Creating a custom validator for password:**

```go
func Password() ValidatorFunc {
    return func(fieldValue string, _ string) string {
        if len(fieldValue) < 8 {
            return "password must be at least 8 characters long"
        }
        if !regexp.MustCompile(`[0-9]`).MatchString(fieldValue) {
            return "password must contain at least one digit"
        }
        if !regexp.MustCompile(`[A-Z]`).MatchString(fieldValue) {
            return "password must contain at least one uppercase letter"
        }
        return ""
    }
}

// Register the validator
func init() {
    RegisterValidator("password", Password())
}
```

**Using the custom validator in a struct:**

```go
type Credentials struct {
    Password string `json:"password" purify:"required|password"`
}
```

**Usage example:**

```go
package main

import (
    "fmt"
    "github.com/contributors-company/purify"
)

type Credentials struct {
    Password string `json:"password" purify:"required|password"`
}

func main() {
    creds := Credentials{
        Password: "Passw1", // Invalid password
    }

    // Validate the struct
    err := purify.ValidateStruct(creds)
    if err != nil {
        validationErrors := err.(purify.ValidationError)
        for field, errors := range validationErrors.Errors {
            fmt.Printf("Field '%s' has the following errors:\n", field)
            for _, errMsg := range errors {
                fmt.Printf("- %s\n", errMsg)
            }
        }
    } else {
        fmt.Println("Validation successful")
    }
}
```

**Expected output:**

```
Field 'Password' has the following errors:
- password must be at least 8 characters long
- password must contain at least one uppercase letter
```

## Error Handling

The `ValidateStruct` function returns `nil` if there are no validation errors or an error object with the following structure:

```json
{
    "errors": {
        "field_name": ["error_message1", "error_message2"],
        "another_field": ["error_message"]
    },
    "message": "General error message"
}
```

**Example of error handling:**

```go
err := purify.ValidateStruct(user)
if err != nil {
    validationErrors := err.(purify.ValidationError)
    for field, errors := range validationErrors.Errors {
        fmt.Printf("Field '%s' has the following errors:\n", field)
        for _, errMsg := range errors {
            fmt.Printf("- %s\n", errMsg)
        }
    }
} else {
    fmt.Println("Validation successful")
}
```

## Testing

You can run tests using the command:

```bash
go test
```

Sample test from the library:

```go
package purify

import (
    "testing"
)

type TestStruct struct {
    Name     string `json:"name" purify:"required|min(3)|max(20)"`
    Email    string `json:"email" purify:"required|email"`
    Password string `json:"password" purify:"required|password"`
}

func TestValidator(t *testing.T) {
    s := TestStruct{
        Name:     "Al",
        Email:    "invalid_email",
        Password: "pass",
    }

    // Validate the struct
    err := ValidateStruct(s)
    if err == nil {
        t.Errorf("Expected validation error, got nil")
    } else {
        validationErrors := err.(ValidationError)
        if len(validationErrors.Errors) != 3 {
            t.Errorf("Expected 3 validation errors, got %d", len(validationErrors.Errors))
        }
    }
}
```
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.