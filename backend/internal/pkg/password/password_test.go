package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "simple password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "complex password",
			password: "C0mpl3x!P@ssw0rd#2024",
			wantErr:  false,
		},
		{
			name:     "chinese characters",
			password: "密码测试123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "very long password",
			password: "this_is_a_very_long_password_that_should_still_work_fine_1234567890",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash)
				assert.Contains(t, hash, "$argon2id$")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "test-password-123"

	// Hash the password
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "incorrect password",
			password: "wrong-password",
			hash:     hash,
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			expected: false,
		},
		{
			name:     "invalid hash format",
			password: password,
			hash:     "invalid-hash",
			expected: false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := VerifyPassword(tt.password, tt.hash)
			if tt.hash == "" || tt.hash == "invalid-hash" {
				// Invalid hash should return error
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestHashPassword_Uniqueness(t *testing.T) {
	password := "same-password-123"

	// Generate two hashes for the same password
	hash1, err := HashPassword(password)
	assert.NoError(t, err)

	hash2, err := HashPassword(password)
	assert.NoError(t, err)

	// Both should be valid and verify the password
	valid1, err := VerifyPassword(password, hash1)
	assert.NoError(t, err)
	assert.True(t, valid1)

	valid2, err := VerifyPassword(password, hash2)
	assert.NoError(t, err)
	assert.True(t, valid2)

	// But the hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2)
}

func TestVerifyPassword_MultiplePasswords(t *testing.T) {
	passwords := []string{
		"password1",
		"admin123",
		"test@test.com",
		"  spaces  ",
		"特殊字符!@#$%",
	}

	for _, pwd := range passwords {
		t.Run(pwd, func(t *testing.T) {
			hash, err := HashPassword(pwd)
			assert.NoError(t, err)
			valid, err := VerifyPassword(pwd, hash)
			assert.NoError(t, err)
			assert.True(t, valid)
		})
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key1 := GenerateAPIKey()
	key2 := GenerateAPIKey()

	// Keys should not be empty
	assert.NotEmpty(t, key1)
	assert.NotEmpty(t, key2)

	// Keys should have the correct prefix
	assert.Contains(t, key1, "ling_")
	assert.Contains(t, key2, "ling_")

	// Keys should be unique
	assert.NotEqual(t, key1, key2)
}

func TestErrInvalidHash(t *testing.T) {
	_, err := VerifyPassword("password", "invalid-format-hash")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidHash, err)
}