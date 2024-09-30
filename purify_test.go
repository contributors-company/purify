package purify

import (
	"testing"
)

type TestStruct struct {
	Name string `json:"name" purify:"email"`
}

func TestValidator(t *testing.T) {
	s := TestStruct{Name: "alexganbert@gmail.com"};

	// Валидация структуры
	err := ValidateStruct(s)
	if err != nil {
		t.Errorf("expected nil, got %v", err);
	}

}