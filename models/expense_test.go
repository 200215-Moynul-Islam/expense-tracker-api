package models

import (
	"os"
	"path/filepath"
	"testing"
)

// setupExpenseCSV points ExpenseCSVPath at a temp file and returns a cleanup func.
func setupExpenseCSV(t *testing.T, content string) func() {
	t.Helper()

	dir := t.TempDir()
	tmpPath := filepath.Join(dir, "expenses.csv")

	if content != "" {
		if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write temp expenses.csv: %v", err)
		}
	}

	orig := ExpenseCSVPath
	ExpenseCSVPath = tmpPath

	return func() { ExpenseCSVPath = orig }
}

var testExpenses = []Expense{
	{ID: 1, UserID: 1, Title: "Lunch", Amount: 350.50, Category: "Food", Note: "team lunch", ExpenseDate: "2025-06-10", CreatedAt: "2025-06-10T14:00:00Z"},
	{ID: 2, UserID: 1, Title: "Bus", Amount: 50.00, Category: "Transport", Note: "", ExpenseDate: "2025-06-11", CreatedAt: "2025-06-11T08:00:00Z"},
	{ID: 3, UserID: 2, Title: "Rent", Amount: 15000.00, Category: "Housing", Note: "monthly", ExpenseDate: "2025-06-01", CreatedAt: "2025-06-01T00:00:00Z"},
}

func seedTestExpenses(t *testing.T, expenses []Expense) {
	t.Helper()
	for _, e := range expenses {
		if err := CreateExpense(e); err != nil {
			t.Fatalf("seed CreateExpense id=%d: %v", e.ID, err)
		}
	}
}

// ---- GetAllExpenses ----

func TestGetAllExpenses_Empty(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()

	expenses, err := GetAllExpenses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(expenses) != 0 {
		t.Errorf("expected 0 expenses, got %d", len(expenses))
	}
}

func TestGetAllExpenses_SkipsMalformedRows(t *testing.T) {
	cleanup := setupExpenseCSV(t,
		"id,user_id,title,amount,category,note,expense_date,created_at\n"+
			"not-id,1,Lunch,350.50,Food,,2025-06-10,2025-06-10T14:00:00Z\n"+ // bad id
			"2,not-uid,Bus,50,Transport,,2025-06-11,2025-06-11T08:00:00Z\n"+ // bad user_id
			"3,1,Dinner,not-amount,Food,,2025-06-12,2025-06-12T20:00:00Z\n"+ // bad amount
			"4,1,Coffee,20.00,Food,,2025-06-13,2025-06-13T09:00:00Z\n"+      // valid
			"5,1\n", // too few columns
	)
	defer cleanup()

	expenses, err := GetAllExpenses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(expenses) != 1 {
		t.Fatalf("expected 1 valid expense, got %d", len(expenses))
	}
	if expenses[0].ID != 4 {
		t.Errorf("expected id=4, got %d", expenses[0].ID)
	}
}

// ---- GetExpensesByUserID ----

func TestGetExpensesByUserID(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	tests := []struct {
		name      string
		userID    int
		wantCount int
	}{
		{"user 1 has 2 expenses", 1, 2},
		{"user 2 has 1 expense", 2, 1},
		{"unknown user has no expenses", 99, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses, err := GetExpensesByUserID(tt.userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(expenses) != tt.wantCount {
				t.Errorf("got %d, want %d", len(expenses), tt.wantCount)
			}
		})
	}
}

// ---- GetExpenseByID ----

func TestGetExpenseByID(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	tests := []struct {
		name      string
		id        int
		userID    int
		wantFound bool
		wantTitle string
	}{
		{"exists and belongs to user", 1, 1, true, "Lunch"},
		{"expense belongs to different user", 3, 1, false, ""},
		{"non-existent id", 999, 1, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := GetExpenseByID(tt.id, tt.userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantFound {
				if e == nil {
					t.Fatalf("expected expense id=%d, got nil", tt.id)
				}
				if e.Title != tt.wantTitle {
					t.Errorf("got title %q, want %q", e.Title, tt.wantTitle)
				}
			} else if e != nil {
				t.Errorf("expected nil, got %+v", e)
			}
		})
	}
}

// ---- CreateExpense ----

func TestCreateExpense(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()

	e := Expense{
		ID:          1,
		UserID:      1,
		Title:       "Coffee",
		Amount:      120.00,
		Category:    "Food",
		Note:        "morning coffee",
		ExpenseDate: "2025-06-15",
		CreatedAt:   "2025-06-15T08:00:00Z",
	}

	if err := CreateExpense(e); err != nil {
		t.Fatalf("CreateExpense error: %v", err)
	}

	all, err := GetAllExpenses()
	if err != nil {
		t.Fatalf("GetAllExpenses error: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 expense, got %d", len(all))
	}
	if all[0].Amount != 120.00 {
		t.Errorf("amount: got %v, want 120.00", all[0].Amount)
	}
	if all[0].Note != "morning coffee" {
		t.Errorf("note: got %q, want %q", all[0].Note, "morning coffee")
	}
}

// ---- UpdateExpense ----

func TestUpdateExpense(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	updated := Expense{
		ID:          1,
		UserID:      1,
		Title:       "Updated Lunch",
		Amount:      500.00,
		Category:    "Food",
		Note:        "updated note",
		ExpenseDate: "2025-06-10",
		CreatedAt:   "2025-06-10T14:00:00Z",
	}

	if err := UpdateExpense(updated); err != nil {
		t.Fatalf("UpdateExpense error: %v", err)
	}

	e, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("GetExpenseByID error: %v", err)
	}
	if e == nil {
		t.Fatal("expense not found after update")
	}
	if e.Title != "Updated Lunch" {
		t.Errorf("title: got %q, want %q", e.Title, "Updated Lunch")
	}
	if e.Amount != 500.00 {
		t.Errorf("amount: got %v, want 500.00", e.Amount)
	}

	all, _ := GetAllExpenses()
	if len(all) != 3 {
		t.Errorf("total count should stay 3, got %d", len(all))
	}
}

func TestUpdateExpense_NotFound(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	ghost := Expense{ID: 999, UserID: 1, Title: "Ghost", Amount: 1, Category: "Food", ExpenseDate: "2025-06-01", CreatedAt: "2025-06-01T00:00:00Z"}
	if err := UpdateExpense(ghost); err == nil {
		t.Error("expected error for non-existent expense, got nil")
	}
}

// ---- DeleteExpense ----

func TestDeleteExpense(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	if err := DeleteExpense(1, 1); err != nil {
		t.Fatalf("DeleteExpense error: %v", err)
	}

	all, err := GetAllExpenses()
	if err != nil {
		t.Fatalf("GetAllExpenses error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 expenses after delete, got %d", len(all))
	}

	e, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("GetExpenseByID error: %v", err)
	}
	if e != nil {
		t.Error("deleted expense still found")
	}
}

func TestDeleteExpense_WrongUser(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	// expense 1 belongs to user 1 — user 2 must not be able to delete it
	if err := DeleteExpense(1, 2); err == nil {
		t.Error("expected error when deleting another user's expense")
	}

	all, _ := GetAllExpenses()
	if len(all) != 3 {
		t.Errorf("count should be unchanged at 3, got %d", len(all))
	}
}

func TestDeleteExpense_NotFound(t *testing.T) {
	cleanup := setupExpenseCSV(t, "")
	defer cleanup()
	seedTestExpenses(t, testExpenses)

	if err := DeleteExpense(999, 1); err == nil {
		t.Error("expected error for non-existent expense id")
	}
}

// ---- GetNextExpenseID ----

func TestGetNextExpenseID(t *testing.T) {
	tests := []struct {
		name     string
		seed     []Expense
		wantNext int
	}{
		{"no expenses returns 1", nil, 1},
		{"returns max id + 1", testExpenses, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupExpenseCSV(t, "")
			defer cleanup()
			seedTestExpenses(t, tt.seed)

			next, err := GetNextExpenseID()
			if err != nil {
				t.Fatalf("GetNextExpenseID error: %v", err)
			}
			if next != tt.wantNext {
				t.Errorf("got %d, want %d", next, tt.wantNext)
			}
		})
	}
}

// ---- IsValidCategory ----

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		category string
		want     bool
	}{
		{"Food", true},
		{"Transport", true},
		{"Housing", true},
		{"Entertainment", true},
		{"Shopping", true},
		{"Healthcare", true},
		{"Education", true},
		{"Utilities", true},
		{"Other", true},
		{"food", false},
		{"FOOD", false},
		{"Grocery", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			got := IsValidCategory(tt.category)
			if got != tt.want {
				t.Errorf("IsValidCategory(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}