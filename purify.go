package purify

import (
	"fmt"
	"reflect"
	"strings"
)

// Определение типа функции для валидатора
type ValidatorFunc func(fieldValue string, param string) string

// Карта для хранения валидаторов
var validators = make(map[string]ValidatorFunc)

// Структура для хранения ошибок валидации
type ValidateError struct {
    Errors   map[string][]string `json:"errors"`
    Message  string              `json:"message"`
}

// Функция для регистрации валидаторов
func RegisterValidator(name string, fn ValidatorFunc) {
    validators[name] = fn
}

// Основная функция для валидации на основе зарегистрированных валидаторов
func ValidateStruct(s interface{}) *ValidateError {
    v := reflect.ValueOf(s)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return &ValidateError{
            Errors:  map[string][]string{"": {"expected a struct"}},
            Message: "expected a struct",
        }
    }

    t := v.Type()
    validationErrors := make(map[string][]string)
    var firstErrorMessage string

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" || jsonTag == "-" {
            jsonTag = fieldType.Name
        }

        gformTag := fieldType.Tag.Get("purify")
        if gformTag == "" {
            continue
        }

        fieldName := jsonTag
        rules := strings.Split(gformTag, "|")
        fieldValueStr := fmt.Sprintf("%v", field.Interface())

        for _, rule := range rules {
            ruleName, param := parseRule(rule)

            if validator, exists := validators[ruleName]; exists {
                if errMsg := validator(fieldValueStr, param); errMsg != "" {
                    validationErrors[fieldName] = append(validationErrors[fieldName], errMsg)
                    if firstErrorMessage == "" {
                        firstErrorMessage = errMsg
                    }
                }
            }
        }
    }

    if len(validationErrors) > 0 {
        return &ValidateError{
            Errors:  validationErrors,
            Message: firstErrorMessage,
        }
    }

    return nil
}

func parseRule(rule string) (string, string) {
    idx := strings.Index(rule, "(")
    if idx == -1 {
        return rule, ""
    }
    ruleName := rule[:idx]
    param := strings.TrimRight(rule[idx+1:], ")")
    return ruleName, param
}

