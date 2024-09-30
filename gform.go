package main

import (
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
func validateStruct(s interface{}) *ValidateError {
    validationErrors := make(map[string][]string)
    var firstErrorMessage string

    // Используем reflect для работы с полями структуры
    v := reflect.ValueOf(s)
    t := reflect.TypeOf(s)

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        // Получаем значение тега json
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" {
            jsonTag = fieldType.Name // Если тега нет, используем имя поля
        }

        // Получаем значение тега gform для валидации
        gformTag := fieldType.Tag.Get("gform")

        // Если тег существует, обрабатываем его
        if gformTag != "" {
            fieldName := jsonTag // Используем имя из тега json
            rules := strings.Split(gformTag, "|")

            for _, rule := range rules {
                parts := strings.Split(rule, "(")
                ruleName := parts[0]
                var param string
                if len(parts) > 1 {
                    param = strings.TrimRight(parts[1], ")")
                }

                // Ищем валидатор по имени
                if validator, exists := validators[ruleName]; exists {
                    if errMsg := validator(field.String(), param); errMsg != "" {
                        validationErrors[fieldName] = append(validationErrors[fieldName], errMsg)
                        if firstErrorMessage == "" {
                            firstErrorMessage = errMsg // Запоминаем первое сообщение
                        }
                    }
                }
            }
        }
    }

    // Если есть ошибки, возвращаем ValidateError
    if len(validationErrors) > 0 {
        return &ValidateError{
            Errors:  validationErrors,
            Message: firstErrorMessage, // Первое сообщение для удобства
        }
    }

    return nil // Возвращаем nil, если ошибок нет
}

func main() {}
