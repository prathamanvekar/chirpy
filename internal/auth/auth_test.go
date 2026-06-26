package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWTFlow(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "prathamaaaaaaaaaaaa"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed unexpectedly: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed to parse a valid token: %v", err)
	}

	if parsedID != userID {
		t.Errorf("ValidateJWT returned UUID %s; want %s", parsedID, userID)
	}
}

// Testing an Expired Token
func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "prathamaaaaaaaaaaaa"

	token, _ := MakeJWT(userID, tokenSecret, -time.Hour)

	_, err := ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected an error for an expired token, but got nil")
	}
}

// Testing the Wrong Secret Key
func TestWrongSecretKey(t *testing.T) {
	userID := uuid.New()

	token, _ := MakeJWT(userID, "aaaa", time.Minute)

	_, err := ValidateJWT(token, "bbbb")

	if err == nil {
		t.Errorf("Expected signature verification to fail with wrong secret, but got nil")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed unexpectedly: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		}, {
			name:        "Wrong Secret",
			tokenString: validToken,
			tokenSecret: "bomboclat",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wannaCheckID, err := ValidateJWT(test.tokenString, test.tokenSecret)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, test.wantErr)
			}
			if wannaCheckID != userID {
				t.Errorf("Parsed ID in not what was expected, ValidateJWT failed, ValidateJWT() gotUserID = %v, want %v", wannaCheckID, test.wantUserID)
			}
		})
	}

}
