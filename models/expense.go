package models

import (
	"encoding/csv"
	"os"
	"strconv"

	"expense-tracker-api/utils"

	logs "github.com/beego/beego/v2/core/logs"
)

const (
	totalExpenseFields = 8
	expenseCSVPath = "data/expenses.csv"
	expenseFilePermission = 0644
)

var expenseCSVHeader = []string{
	"id",
	"user_id",
	"title",
	"amount",
	"category",
	"note",
	"expense_date",
	"created_at",
}

var AllowedCategories = []string{
	"Food",
	"Transport",
	"Housing",
	"Entertainment",
	"Shopping",
	"Healthcare",
	"Education",
	"Utilities",
	"Other",
}

type Expense struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Title string `json:"title"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
	Note string `json:"note"`
	ExpenseDate string `json:"expense_date"`
	CreatedAt string `json:"created_at"`
}

func GetAllExpenses() ([]Expense, error) {
	if err := utils.EnsureCSVFile(expenseCSVPath, expenseCSVHeader); err != nil {
		return nil, err
	}

	file, err := os.Open(expenseCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	expenses := make([]Expense, 0, len(records))

	for i := 1; i < len(records); i++ {
		// An expense row should have at least 8 columns
		// (id, user_id, title, amount, category, note, expense_date, created_at).
		if len(records[i]) < totalExpenseFields {
			// Log an error for the invalid expense and skip this record.
			// We don't want to stop processing all expenses just because one record is bad.
			logs.Warn("invalid expense row:", records[i])
			continue
		}

		id, err := strconv.Atoi(records[i][0])
		if err != nil {
			// Log an error for the invalid expense ID and skip this record.
			logs.Warn("invalid expense id:", records[i])
			continue
		}

		userID, err := strconv.Atoi(records[i][1])
		if err != nil {
			// Log an error for the invalid user ID and skip this record.
			logs.Warn("invalid expense user id:", records[i])
			continue
		}

		amount, err := strconv.ParseFloat(records[i][3], 64)
		if err != nil {
			// Log an error for the invalid amount and skip this record.
			logs.Warn("invalid expense amount:", records[i])
			continue
		}

		expenses = append(expenses, Expense{
			ID: id,
			UserID: userID,
			Title: records[i][2],
			Amount: amount,
			Category: records[i][4],
			Note: records[i][5],
			ExpenseDate: records[i][6],
			CreatedAt: records[i][7],
		})
	}

	return expenses, nil
}

func GetExpensesByUserID(userID int) ([]Expense, error) {
	expenses, err := GetAllExpenses()
	if err != nil {
		return nil, err
	}

	result := make([]Expense, 0)

	for _, e := range expenses {
		if e.UserID == userID {
			result = append(result, e)
		}
	}

	return result, nil
}

func GetExpenseByID(id int, userID int) (*Expense, error) {
	expenses, err := GetAllExpenses()
	if err != nil {
		return nil, err
	}

	for _, e := range expenses {
		if e.ID == id && e.UserID == userID {
			return &e, nil
		}
	}

	return nil, nil
}
