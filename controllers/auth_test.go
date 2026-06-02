package controllers

import (
	"testing"
)

// ---- validateRegisterRequest ----

func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name        string
		input       RegisterRequest
		wantErrMsg  string // empty means no error expected
	}{
		{
			name:  "valid registration",
			input: RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret99"},
		},
		{
			name:       "missing name",
			input:      RegisterRequest{Name: "", Email: "alice@example.com", Password: "secret99"},
			wantErrMsg: "Name is required.",
		},
		{
			name:       "missing email",
			input:      RegisterRequest{Name: "Alice", Email: "", Password: "secret99"},
			wantErrMsg: "Email is required.",
		},
		{
			name:       "invalid email format",
			input:      RegisterRequest{Name: "Alice", Email: "not-an-email", Password: "secret99"},
			wantErrMsg: "Invalid email format.",
		},
		{
			name:       "missing password",
			input:      RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: ""},
			wantErrMsg: "Password is required.",
		},
		{
			name:       "password too short",
			input:      RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "abc"},
			wantErrMsg: "Password must be at least 6 characters.",
		},
		{
			name:  "password exactly 6 chars is valid",
			input: RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "abcdef"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := validateRegisterRequest(tt.input)
			if err != nil {
				t.Fatalf("unexpected engine error: %v", err)
			}
			if tt.wantErrMsg == "" && msg != "" {
				t.Errorf("expected no error, got %q", msg)
			}
			if tt.wantErrMsg != "" && msg != tt.wantErrMsg {
				t.Errorf("got %q, want %q", msg, tt.wantErrMsg)
			}
		})
	}
}

// ---- validateLoginRequest ----

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name       string
		input      LoginRequest
		wantErrMsg string
	}{
		{
			name:  "valid login",
			input: LoginRequest{Email: "alice@example.com", Password: "secret99"},
		},
		{
			name:       "missing email",
			input:      LoginRequest{Email: "", Password: "secret99"},
			wantErrMsg: "Email is required.",
		},
		{
			name:       "invalid email format",
			input:      LoginRequest{Email: "bad-email", Password: "secret99"},
			wantErrMsg: "Invalid email format.",
		},
		{
			name:       "missing password",
			input:      LoginRequest{Email: "alice@example.com", Password: ""},
			wantErrMsg: "Password is required.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := validateLoginRequest(tt.input)
			if err != nil {
				t.Fatalf("unexpected engine error: %v", err)
			}
			if tt.wantErrMsg == "" && msg != "" {
				t.Errorf("expected no error, got %q", msg)
			}
			if tt.wantErrMsg != "" && msg != tt.wantErrMsg {
				t.Errorf("got %q, want %q", msg, tt.wantErrMsg)
			}
		})
	}
}

// ---- normalizeRegisterRequest ----

func TestNormalizeRegisterRequest(t *testing.T) {
	tests := []struct {
		name      string
		input     RegisterRequest
		wantName  string
		wantEmail string
	}{
		{
			name:      "trims whitespace",
			input:     RegisterRequest{Name: "  Alice  ", Email: "  Alice@Example.COM  ", Password: "  pass  "},
			wantName:  "Alice",
			wantEmail: "alice@example.com",
		},
		{
			name:      "lowercases email",
			input:     RegisterRequest{Name: "Bob", Email: "BOB@EXAMPLE.COM", Password: "pass"},
			wantEmail: "bob@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.input
			normalizeRegisterRequest(&req)

			if tt.wantName != "" && req.Name != tt.wantName {
				t.Errorf("Name: got %q, want %q", req.Name, tt.wantName)
			}
			if req.Email != tt.wantEmail {
				t.Errorf("Email: got %q, want %q", req.Email, tt.wantEmail)
			}
		})
	}
}