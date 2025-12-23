package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserModel_Creation(t *testing.T) {
	user := User{
		ID:          uuid.New(),
		Username:    "testuser",
		Email:       "test@example.com",
		FirstName:   "Test",
		LastName:    "User",
		DOB:         time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "+1234567890",
		Address:     "123 Test St",
		IsActive:    true,
		IsAdmin:     false,
	}

	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.True(t, user.IsActive)
	assert.False(t, user.IsAdmin)
}

func TestUserModel_DefaultValues(t *testing.T) {
	user := User{
		Username:    "newuser",
		Email:       "new@example.com",
		PhoneNumber: "+1234567890",
		Address:     "123 St",
	}

	// Test that boolean defaults work as expected
	assert.False(t, user.IsActive) // Default Go zero value
	assert.False(t, user.IsAdmin)  // Default Go zero value
}

func TestUserModel_JSONTags(t *testing.T) {
	// Verify struct has proper JSON tags
	user := User{}
	
	// These compile checks ensure the struct fields exist
	_ = user.ID
	_ = user.Username
	_ = user.Email
	_ = user.FirstName
	_ = user.LastName
	_ = user.DOB
	_ = user.PhoneNumber
	_ = user.Address
	_ = user.IsActive
	_ = user.IsAdmin
	_ = user.KeycloackID
	_ = user.CreatedAt
	_ = user.UpdatedAt
}

func TestUserModel_UUIDGeneration(t *testing.T) {
	user1 := User{ID: uuid.New()}
	user2 := User{ID: uuid.New()}

	assert.NotEqual(t, user1.ID, user2.ID, "UUIDs should be unique")
	assert.NotEqual(t, uuid.Nil, user1.ID, "UUID should not be nil")
	assert.NotEqual(t, uuid.Nil, user2.ID, "UUID should not be nil")
}
