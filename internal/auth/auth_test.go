package auth

import (
	"testing"
	"time"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func TestHash(t *testing.T) {
	cases := []struct {
		input     string
		expected1 bool
		expected2 error
		expected3 error
	}{
		{
			input:   "password",
			expected1: true,
			expected2: nil,
			expected3: nil,
		},
		{
			input:   "",
			expected1: true,
			expected2: nil,
			expected3: nil,
		},
		{
			input:   "pass word",
			expected1: true,
			expected2: nil,
			expected3: nil,
		},
	}

	for _, c := range cases {
		hash, actual2 := HashPassword(c.input)
		if !errors.Is(actual2, c.expected2) {
			t.Errorf("Error: Actual error does not match expected error")
			t.Errorf("Actual: %s   Expected: %s", actual2, c.expected2)
		}

		actual1, actual3 := CheckPasswordHash(c.input, hash)
		if actual1 != c.expected1 {
			t.Errorf("Error: Actual hash does not match expected hash")
			t.Errorf("Actual: %t   Expected: %t", actual1 , c.expected1)
		}
		if !errors.Is(actual3, c.expected3) {
			t.Errorf("Error: Actual error does not match expected error")
			t.Errorf("Actual: %s   Expected: %s", actual2, c.expected2)
		}
	}
}

func TestJWT(t *testing.T) {
	min, _ := time.ParseDuration("1m")
	nano, _ := time.ParseDuration("1ns")

	cases := []struct {
		input1	  uuid.UUID
		input2    string
		input3    time.Duration
		expected1 uuid.UUID
		expected2 error
		expected3 error
	}{
		{
			input1:   uuid.New(),
			input2:   "qwerty",
			input3:   min,
			expected2: nil,
			expected3: nil,
		},
		{
			input1:   uuid.New(),
			input2:   "Its a secret to everyone",
			input3:   min,
			expected2: nil,
			expected3: nil,
		},
		{
			input1:   uuid.Nil,
			input2:   "this should fail",
			input3:   nano,
			expected2: nil,
			expected3: jwt.ErrTokenInvalidClaims,
		},
	}

	for _, c := range cases {
		c.expected1 = c.input1
		token, actual2 := MakeJWT(c.input1, c.input2, c.input3)
		if !errors.Is(actual2, c.expected2) {
			t.Errorf("Error: Actual error does not match expected error")
			t.Errorf("Actual: %s   Expected: %s", actual2, c.expected2)
		}

		actual1, actual3 := ValidateJWT(token, c.input2)
		if actual1 != c.expected1 {
			t.Errorf("Error: Actual uuid does not match expected uuid")
			t.Errorf("Actual: %s   Expected: %s", actual1 , c.expected1)
		}
		if !errors.Is(actual3, c.expected3) {
			t.Errorf("Error: Actual error does not match expected error")
			t.Errorf("Actual: %s   Expected: %s", actual3, c.expected3)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		input1    http.Header
		input2    string
		input3    []string
		expected1 string
		wantErr   bool
	}{
		{
			input1:    http.Header{},
			input2:    "Authorization",
			input3:    []string{"bearer", "This is probably a token"},
			expected1: "This is probably a token",
			wantErr:   false,
		},
		{
			input1:     http.Header{},
			input2:    "NotAuthorization",
			input3:    []string{"bearer", "This is probably a token"},
			expected1: "",
			wantErr:   true,
		},
		{
			input1:     http.Header{},
			input2:    "Authorization",
			input3:    []string{"bearer"},
			expected1: "",
			wantErr:   true,
		},
		{
			input1:     http.Header{},
			input2:    "NotAuthorization",
			input3:    []string{"beaer", "This is probably a token"},
			expected1: "",
			wantErr:   true,
		},
	}

	for _, c := range cases {
		c.input1[c.input2] = c.input3
		actual1, err := GetBearerToken(c.input1)

		if actual1 != c.expected1 {
			t.Errorf("Error: Actual token does not match expected token")
			t.Errorf("Actual: %s   Expected: %s", actual1 , c.expected1)
		}
		if (err != nil) != c.wantErr {
			t.Errorf("GetBearerToken() error = %v, wantErr %v", err, c.wantErr)
			return
		}
	}
}