package codes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassItCodes_Defined(t *testing.T) {
	// Test that all expected codes are defined
	codes := []int{
		UserCreatedSuccessfully,
		UserUpdatedSuccessfully,
		UserLoggedInSuccessfully,
		JobsRetrievedSuccessfully,
		GetJobBadRequest,
		JobIdNotFound,
	}

	for _, code := range codes {
		assert.Greater(t, code, 0, "Code should be positive")
		assert.Less(t, code, 10000, "Code should be reasonable range")
	}
}

func TestPassItCodes_Uniqueness(t *testing.T) {
	// Test that codes are unique
	codes := map[string]int{
		"UserCreatedSuccessfully":   UserCreatedSuccessfully,
		"UserUpdatedSuccessfully":   UserUpdatedSuccessfully,
		"UserLoggedInSuccessfully":  UserLoggedInSuccessfully,
		"JobsRetrievedSuccessfully": JobsRetrievedSuccessfully,
		"GetJobBadRequest":          GetJobBadRequest,
		"JobIdNotFound":             JobIdNotFound,
	}

	seenCodes := make(map[int]string)
	for name, code := range codes {
		if existingName, exists := seenCodes[code]; exists {
			t.Errorf("Duplicate code %d found: %s and %s", code, name, existingName)
		}
		seenCodes[code] = name
	}
}

func TestPassItCodes_Values(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"UserCreatedSuccessfully", UserCreatedSuccessfully, 201},
		{"UserUpdatedSuccessfully", UserUpdatedSuccessfully, 200},
		{"UserLoggedInSuccessfully", UserLoggedInSuccessfully, 202},
		{"JobsRetrievedSuccessfully", JobsRetrievedSuccessfully, 205},
		{"GetJobBadRequest", GetJobBadRequest, 400},
		{"JobIdNotFound", JobIdNotFound, 405},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code)
		})
	}
}
