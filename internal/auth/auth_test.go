package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

const testSecret = "super-secret-key"

func TestMakeAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error creating token: %v", err)
	}

	parsedID, err := ValidateJWT(token, testSecret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if parsedID != userID {
		t.Errorf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error creating token: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error with wrong secret, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()

	// Token that expires immediately
	token, err := MakeJWT(userID, testSecret, -time.Minute)
	if err != nil {
		t.Fatalf("unexpected error creating token: %v", err)
	}

	_, err = ValidateJWT(token, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateJWT_InvalidTokenString(t *testing.T) {
	invalidToken := "this.is.not.a.valid.token"

	_, err := ValidateJWT(invalidToken, testSecret)
	if err == nil {
		t.Fatal("expected error for invalid token string, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headerValue   string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "valid bearer token",
			headerValue:   "Bearer abc123",
			expectedToken: "abc123",
			expectError:   false,
		},
		{
			name:          "valid bearer token with extra spaces",
			headerValue:   "   Bearer    abc123   ",
			expectedToken: "abc123",
			expectError:   false,
		},
		{
			name:          "case insensitive bearer",
			headerValue:   "bearer abc123",
			expectedToken: "abc123",
			expectError:   false,
		},
		{
			name:        "missing header",
			headerValue: "",
			expectError: true,
		},
		{
			name:        "missing token",
			headerValue: "Bearer",
			expectError: true,
		},
		{
			name:        "wrong scheme",
			headerValue: "Basic abc123",
			expectError: true,
		},
		{
			name:        "too many parts",
			headerValue: "Bearer abc 123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerValue != "" {
				headers.Set("Authorization", tt.headerValue)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != tt.expectedToken {
				t.Fatalf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}
