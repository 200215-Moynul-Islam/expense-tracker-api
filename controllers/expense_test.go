package controllers

import (
	"testing"

	"expense-tracker-api/models"
)

// ---- validateCreateExpenseRequest ----

func TestValidateCreateExpenseRequest(t *testing.T) {
	tests := []struct {
		name       string
		input      CreateExpenseRequest
		wantErrMsg string
	}{
		{
			name:  "valid request",
			input: CreateExpenseRequest{Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"},
		},
		{
			name:       "missing title",
			input:      CreateExpenseRequest{Title: "", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"},
			wantErrMsg: "Title is required.",
		},
		{
			name:       "missing category",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: 350.50, Category: "", ExpenseDate: "2025-06-10"},
			wantErrMsg: "Category is required.",
		},
		{
			name:       "invalid category",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: 350.50, Category: "Snacks", ExpenseDate: "2025-06-10"},
			wantErrMsg: "Invalid category.",
		},
		{
			name:       "missing expense date",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: ""},
			wantErrMsg: "Expense date is required.",
		},
		{
			name:       "invalid date format",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "10-06-2025"},
			wantErrMsg: "Expense date must be in YYYY-MM-DD format.",
		},
		{
			name:       "zero amount",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: 0, Category: "Food", ExpenseDate: "2025-06-10"},
			wantErrMsg: "Amount must be greater than 0.",
		},
		{
			name:       "negative amount",
			input:      CreateExpenseRequest{Title: "Lunch", Amount: -50, Category: "Food", ExpenseDate: "2025-06-10"},
			wantErrMsg: "Amount must be greater than 0.",
		},
		{
			name:  "note is optional",
			input: CreateExpenseRequest{Title: "Lunch", Amount: 100, Category: "Food", Note: "", ExpenseDate: "2025-06-10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := validateCreateExpenseRequest(tt.input)
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

// ---- validateExpenseListQueryParams ----

func TestValidateExpenseListQueryParams(t *testing.T) {
	tests := []struct {
		name       string
		input      ExpenseListQueryParams
		wantErrMsg string
	}{
		{name: "all empty params are valid", input: ExpenseListQueryParams{}},
		{name: "valid category", input: ExpenseListQueryParams{Category: "Food"}},
		{
			name:       "invalid category",
			input:      ExpenseListQueryParams{Category: "Junk"},
			wantErrMsg: "Invalid category.",
		},
		{
			name:  "valid date range",
			input: ExpenseListQueryParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"},
		},
		{
			name:       "invalid date_from format",
			input:      ExpenseListQueryParams{DateFrom: "01/06/2025"},
			wantErrMsg: "Invalid date_from format. Use YYYY-MM-DD.",
		},
		{
			name:       "invalid date_to format",
			input:      ExpenseListQueryParams{DateTo: "2025-13-01"},
			wantErrMsg: "Invalid date_to format. Use YYYY-MM-DD.",
		},
		{
			name:       "date_to before date_from",
			input:      ExpenseListQueryParams{DateFrom: "2025-06-30", DateTo: "2025-06-01"},
			wantErrMsg: "date_to cannot be earlier than date_from.",
		},
		{name: "sort by amount asc", input: ExpenseListQueryParams{SortBy: "amount", SortOrder: "asc"}},
		{name: "sort by expense_date desc", input: ExpenseListQueryParams{SortBy: "expense_date", SortOrder: "desc"}},
		{
			name:       "invalid sort_by",
			input:      ExpenseListQueryParams{SortBy: "title"},
			wantErrMsg: "Invalid sort_by value.",
		},
		{
			name:       "invalid sort_order",
			input:      ExpenseListQueryParams{SortOrder: "random"},
			wantErrMsg: "Invalid sort_order value.",
		},
		{name: "valid limit", input: ExpenseListQueryParams{Limit: "10"}},
		{
			name:       "zero limit",
			input:      ExpenseListQueryParams{Limit: "0"},
			wantErrMsg: "Limit must be greater than 0.",
		},
		{
			name:       "negative limit",
			input:      ExpenseListQueryParams{Limit: "-5"},
			wantErrMsg: "Limit must be greater than 0.",
		},
		{
			name:       "non-numeric limit",
			input:      ExpenseListQueryParams{Limit: "ten"},
			wantErrMsg: "Limit must be greater than 0.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := validateExpenseListQueryParams(tt.input)
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

// ---- validateExpenseSummaryQueryParams ----

func TestValidateExpenseSummaryQueryParams(t *testing.T) {
	tests := []struct {
		name       string
		input      ExpenseSummaryQueryParams
		wantErrMsg string
	}{
		{
			name:  "valid params",
			input: ExpenseSummaryQueryParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"},
		},
		{
			name:       "missing date_from",
			input:      ExpenseSummaryQueryParams{DateTo: "2025-06-30"},
			wantErrMsg: "date_from is required.",
		},
		{
			name:       "missing date_to",
			input:      ExpenseSummaryQueryParams{DateFrom: "2025-06-01"},
			wantErrMsg: "date_to is required.",
		},
		{
			name:       "date_to before date_from",
			input:      ExpenseSummaryQueryParams{DateFrom: "2025-06-30", DateTo: "2025-06-01"},
			wantErrMsg: "date_to cannot be earlier than date_from.",
		},
		{
			name:       "invalid date_from format",
			input:      ExpenseSummaryQueryParams{DateFrom: "June 1 2025", DateTo: "2025-06-30"},
			wantErrMsg: "Invalid date_from format. Use YYYY-MM-DD.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := validateExpenseSummaryQueryParams(tt.input)
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

// ---- filterExpenses ----

func TestFilterExpenses(t *testing.T) {
	expenses := []models.Expense{
		{ID: 1, Amount: 350, Category: "Food", ExpenseDate: "2025-06-10"},
		{ID: 2, Amount: 50, Category: "Transport", ExpenseDate: "2025-06-11"},
		{ID: 3, Amount: 600, Category: "Food", ExpenseDate: "2025-06-20"},
		{ID: 4, Amount: 300, Category: "Entertainment", ExpenseDate: "2025-06-25"},
	}

	tests := []struct {
		name      string
		params    ExpenseListQueryParams
		wantIDs   []int
	}{
		{
			name:    "no filter returns all",
			params:  ExpenseListQueryParams{},
			wantIDs: []int{1, 2, 3, 4},
		},
		{
			name:    "filter by category Food",
			params:  ExpenseListQueryParams{Category: "Food"},
			wantIDs: []int{1, 3},
		},
		{
			name:    "filter by date_from",
			params:  ExpenseListQueryParams{DateFrom: "2025-06-20"},
			wantIDs: []int{3, 4},
		},
		{
			name:    "filter by date_to",
			params:  ExpenseListQueryParams{DateTo: "2025-06-11"},
			wantIDs: []int{1, 2},
		},
		{
			name:    "filter by date range",
			params:  ExpenseListQueryParams{DateFrom: "2025-06-11", DateTo: "2025-06-20"},
			wantIDs: []int{2, 3},
		},
		{
			name:    "category and date range combined",
			params:  ExpenseListQueryParams{Category: "Food", DateFrom: "2025-06-15"},
			wantIDs: []int{3},
		},
		{
			name:    "no match returns empty",
			params:  ExpenseListQueryParams{Category: "Healthcare"},
			wantIDs: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterExpenses(expenses, tt.params)

			if len(result) != len(tt.wantIDs) {
				t.Fatalf("got %d results, want %d", len(result), len(tt.wantIDs))
			}
			for i, id := range tt.wantIDs {
				if result[i].ID != id {
					t.Errorf("index %d: got id=%d, want id=%d", i, result[i].ID, id)
				}
			}
		})
	}
}

// ---- sortExpenses ----

func TestSortExpenses(t *testing.T) {
	base := []models.Expense{
		{ID: 1, Amount: 300, ExpenseDate: "2025-06-10"},
		{ID: 2, Amount: 100, ExpenseDate: "2025-06-20"},
		{ID: 3, Amount: 500, ExpenseDate: "2025-06-05"},
	}

	tests := []struct {
		name      string
		params    ExpenseListQueryParams
		wantFirst int
		wantLast  int
	}{
		{
			name:      "no sort_by keeps original order",
			params:    ExpenseListQueryParams{},
			wantFirst: 1,
			wantLast:  3,
		},
		{
			name:      "sort by amount asc",
			params:    ExpenseListQueryParams{SortBy: "amount", SortOrder: "asc"},
			wantFirst: 2, // 100
			wantLast:  3, // 500
		},
		{
			name:      "sort by amount desc",
			params:    ExpenseListQueryParams{SortBy: "amount", SortOrder: "desc"},
			wantFirst: 3, // 500
			wantLast:  2, // 100
		},
		{
			name:      "sort by expense_date asc",
			params:    ExpenseListQueryParams{SortBy: "expense_date", SortOrder: "asc"},
			wantFirst: 3, // 2025-06-05
			wantLast:  2, // 2025-06-20
		},
		{
			name:      "sort by expense_date desc",
			params:    ExpenseListQueryParams{SortBy: "expense_date", SortOrder: "desc"},
			wantFirst: 2, // 2025-06-20
			wantLast:  3, // 2025-06-05
		},
		{
			name:      "empty sort_order defaults to desc",
			params:    ExpenseListQueryParams{SortBy: "amount", SortOrder: ""},
			wantFirst: 3, // 500
			wantLast:  2, // 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := make([]models.Expense, len(base))
			copy(cp, base)

			sortExpenses(cp, tt.params)

			if cp[0].ID != tt.wantFirst {
				t.Errorf("first: got id=%d, want id=%d", cp[0].ID, tt.wantFirst)
			}
			if cp[len(cp)-1].ID != tt.wantLast {
				t.Errorf("last: got id=%d, want id=%d", cp[len(cp)-1].ID, tt.wantLast)
			}
		})
	}
}

// ---- generateExpenseSummary ----

func TestGenerateExpenseSummary(t *testing.T) {
	expenses := []models.Expense{
		{ID: 1, Amount: 350, Category: "Food", ExpenseDate: "2025-06-10"},
		{ID: 2, Amount: 200, Category: "Food", ExpenseDate: "2025-06-15"},
		{ID: 3, Amount: 100, Category: "Transport", ExpenseDate: "2025-06-20"},
		{ID: 4, Amount: 5000, Category: "Housing", ExpenseDate: "2025-05-01"}, // outside range
	}

	params := ExpenseSummaryQueryParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"}
	summary := generateExpenseSummary(expenses, params)

	if summary.TotalCount != 3 {
		t.Errorf("TotalCount: got %d, want 3", summary.TotalCount)
	}
	if summary.TotalAmount != 650.0 {
		t.Errorf("TotalAmount: got %v, want 650.0", summary.TotalAmount)
	}
	if summary.DateFrom != "2025-06-01" {
		t.Errorf("DateFrom: got %q, want %q", summary.DateFrom, "2025-06-01")
	}
	if len(summary.ByCategory) != 2 {
		t.Fatalf("ByCategory: got %d categories, want 2", len(summary.ByCategory))
	}

	var foodTotal float64
	var foodCount int
	for _, cs := range summary.ByCategory {
		if cs.Category == "Food" {
			foodTotal = cs.Total
			foodCount = cs.Count
		}
	}
	if foodTotal != 550.0 {
		t.Errorf("Food total: got %v, want 550.0", foodTotal)
	}
	if foodCount != 2 {
		t.Errorf("Food count: got %d, want 2", foodCount)
	}
}

func TestGenerateExpenseSummary_EmptyRange(t *testing.T) {
	expenses := []models.Expense{
		{ID: 1, Amount: 500, Category: "Food", ExpenseDate: "2025-05-01"},
	}

	params := ExpenseSummaryQueryParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"}
	summary := generateExpenseSummary(expenses, params)

	if summary.TotalCount != 0 {
		t.Errorf("TotalCount: got %d, want 0", summary.TotalCount)
	}
	if summary.TotalAmount != 0.0 {
		t.Errorf("TotalAmount: got %v, want 0.0", summary.TotalAmount)
	}
	if len(summary.ByCategory) != 0 {
		t.Errorf("ByCategory: got %d, want 0", len(summary.ByCategory))
	}
}

// ---- limitExpenses ----

func TestLimitExpenses(t *testing.T) {
	expenses := []models.Expense{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
	}

	tests := []struct {
		name      string
		limit     string
		wantCount int
		wantErr   bool
	}{
		{"no limit returns all", "", 5, false},
		{"limit 3", "3", 3, false},
		{"limit larger than slice", "100", 5, false},
		{"limit equals slice length", "5", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := make([]models.Expense, len(expenses))
			copy(cp, expenses)

			result, err := limitExpenses(cp, ExpenseListQueryParams{Limit: tt.limit})
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Errorf("got %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}