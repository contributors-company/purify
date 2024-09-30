package purify

// Removed the import statement for the non-existent "purify" package

import (
	"fmt"
	"regexp"
	"strconv"
)

func Min() ValidatorFunc {
    return func(fieldValue string, param string) string {
        maxLength, _ := strconv.Atoi(param)

        if len(fieldValue) < maxLength {
            return fmt.Sprintf("min length is %s", param)
        }
        return "";
    }
}

func Max() ValidatorFunc {
	return func(fieldValue string, param string) string {
		maxLength, _ := strconv.Atoi(param)

		if len(fieldValue) > maxLength {
			return fmt.Sprintf("max length is %s", param)
		}
		return "";
	}
}

func Required() ValidatorFunc {
	return func(fieldValue string, param string) string {
		if fieldValue == "" {
			return "required"
		}
		return "";
	}
}

func Email() ValidatorFunc {
	return func(fieldValue string, param string) string {
		// Simple regex for email validation
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, fieldValue)
		if !matched {
			return "invalid email"
		}
		return ""
	}
}

func init() {
    RegisterValidator("min", Min())
	RegisterValidator("max", Max())
	RegisterValidator("required", Required())
	RegisterValidator("email", Email())
}