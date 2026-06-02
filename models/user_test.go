package models

import (
	"os"
	"path/filepath"
	"testing"
)

// setupUserCSV points UserCSVPath at a temp file and returns a cleanup func.
func setupUserCSV(t *testing.T, content string) func() {
	t.Helper()

	dir := t.TempDir()
	tmpPath := filepath.Join(dir, "users.csv")

	if content != "" {
		if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write temp users.csv: %v", err)
		}
	}

	orig := UserCSVPath
	UserCSVPath = tmpPath

	return func() { UserCSVPath = orig }
}

// ---- GetAllUsers ----

func TestGetAllUsers_Empty(t *testing.T) {
	cleanup := setupUserCSV(t, "")
	defer cleanup()

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUsers_WithData(t *testing.T) {
	cleanup := setupUserCSV(t,
		"id,name,email,password,created_at\n"+
			"1,Alice,alice@example.com,pass123,2025-01-01T00:00:00Z\n"+
			"2,Bob,bob@example.com,pass456,2025-01-02T00:00:00Z\n",
	)
	defer cleanup()

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
	if users[1].Email != "bob@example.com" {
		t.Errorf("expected bob@example.com, got %s", users[1].Email)
	}
}

func TestGetAllUsers_SkipsMalformedRows(t *testing.T) {
	cleanup := setupUserCSV(t,
		"id,name,email,password,created_at\n"+
			"not-an-id,Bad,bad@example.com,pass,2025-01-01T00:00:00Z\n"+ // bad id
			"2,Good,good@example.com,pass,2025-01-01T00:00:00Z\n"+
			"3,Short\n", // too few columns
	)
	defer cleanup()

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 valid user, got %d", len(users))
	}
	if users[0].ID != 2 {
		t.Errorf("expected id=2, got %d", users[0].ID)
	}
}

// ---- GetUserByEmail ----

func TestGetUserByEmail(t *testing.T) {
	cleanup := setupUserCSV(t,
		"id,name,email,password,created_at\n"+
			"1,Alice,alice@example.com,pass123,2025-01-01T00:00:00Z\n",
	)
	defer cleanup()

	tests := []struct {
		name      string
		email     string
		wantFound bool
	}{
		{"existing email", "alice@example.com", true},
		{"non-existent email", "ghost@example.com", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserByEmail(tt.email)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantFound && user == nil {
				t.Errorf("expected user for email %q, got nil", tt.email)
			}
			if !tt.wantFound && user != nil {
				t.Errorf("expected nil for email %q, got %+v", tt.email, user)
			}
			if tt.wantFound && user != nil && user.Email != tt.email {
				t.Errorf("got email %q, want %q", user.Email, tt.email)
			}
		})
	}
}

// ---- CreateUser ----

func TestCreateUser(t *testing.T) {
	cleanup := setupUserCSV(t, "")
	defer cleanup()

	user := User{
		ID:        1,
		Name:      "Charlie",
		Email:     "charlie@example.com",
		Password:  "secret99",
		CreatedAt: "2025-06-01T10:00:00Z",
	}

	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Charlie" {
		t.Errorf("expected Charlie, got %s", users[0].Name)
	}
	if users[0].Email != "charlie@example.com" {
		t.Errorf("expected charlie@example.com, got %s", users[0].Email)
	}
}

func TestCreateUser_Multiple(t *testing.T) {
	cleanup := setupUserCSV(t, "")
	defer cleanup()

	toCreate := []User{
		{ID: 1, Name: "Alice", Email: "a@a.com", Password: "pass", CreatedAt: "2025-01-01T00:00:00Z"},
		{ID: 2, Name: "Bob", Email: "b@b.com", Password: "pass", CreatedAt: "2025-01-02T00:00:00Z"},
	}
	for _, u := range toCreate {
		if err := CreateUser(u); err != nil {
			t.Fatalf("CreateUser error: %v", err)
		}
	}

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

// ---- GetNextID ----

func TestGetNextID(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantID  int
	}{
		{
			name:    "empty file returns 1",
			content: "id,name,email,password,created_at\n",
			wantID:  1,
		},
		{
			name: "returns max id + 1",
			content: "id,name,email,password,created_at\n" +
				"1,A,a@a.com,p,2025-01-01T00:00:00Z\n" +
				"3,B,b@b.com,p,2025-01-01T00:00:00Z\n",
			wantID: 4,
		},
		{
			name: "single user",
			content: "id,name,email,password,created_at\n" +
				"5,E,e@e.com,p,2025-01-01T00:00:00Z\n",
			wantID: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupUserCSV(t, tt.content)
			defer cleanup()

			nextID, err := GetNextID()
			if err != nil {
				t.Fatalf("GetNextID error: %v", err)
			}
			if nextID != tt.wantID {
				t.Errorf("got %d, want %d", nextID, tt.wantID)
			}
		})
	}
}