package utility

import (
	"testing"
)

func TestCEPValidator(t *testing.T) {
	validCEPs := []string{
		"12345678",
		"12345-678",
		"00000-000",
		"99999-999",
	}

	invalidCEPs := []string{
		"1234-5678",
		"1234567",
		"123456789",
		"999999999",
		"12a45-678",
		"12345_678",
		"ABCDE-FFF",
	}

	for _, cep := range validCEPs {
		if _, err := CEPValidator(cep); err != nil {
			t.Errorf("Expected CEP %s to be valid, but got invalid", cep)
		}
	}

	for _, cep := range invalidCEPs {
		if _, err := CEPValidator(cep); err == nil {
			t.Errorf("Expected CEP %s to be invalid, but got valid", cep)
		}
	}
}

func TestCEPFormatter(t *testing.T) {
	tests := map[string]string{
		"12345-678": "12345678",
		"04446-160": "04446160",
		"87654321":  "87654321",
		"00000000":  "00000000",
		"99999999":  "99999999",
	}

	for input, expected := range tests {
		formatted, _ := CEPFormatter(input)
		if formatted != expected {
			t.Errorf("Expected CEPFormatter(%s) to be %s, but got %s", input, expected, formatted)
		}
	}
}
