package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"expense-tracker-api/models"

	logs "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type ExpenseController struct {
	BaseController
}

type CreateExpenseRequest struct {
	Title string `json:"title" valid:"Required"`
	Amount float64 `json:"amount"`
	Category string `json:"category" valid:"Required"`
	Note string `json:"note"`
	ExpenseDate string `json:"expense_date" valid:"Required"`
}

func (c *ExpenseController) Create() {
	userID := c.Ctx.Input.GetData("userID").(int)

	var req CreateExpenseRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Warn("failed to parse create expense request body:", err)

		c.Error(http.StatusBadRequest, "Invalid request body.")
		return
	}

	normalizeCreateExpenseRequest(&req)

	errMessage, err := validateCreateExpenseRequest(req)
	if err != nil {
		logs.Error("failed to validate create expense request:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if errMessage != "" {
		c.Error(http.StatusBadRequest, errMessage)
		return
	}

	id, err := models.GetNextExpenseID()
	if err != nil {
		logs.Error("failed to generate next expense id:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	expense := models.Expense{
		ID: id,
		UserID: userID,
		Title: req.Title,
		Amount: req.Amount,
		Category: req.Category,
		Note: req.Note,
		ExpenseDate: req.ExpenseDate,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	err = models.CreateExpense(expense)
	if err != nil {
		logs.Error("failed to create expense:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	c.Success(http.StatusCreated, "Expense created successfully.", expense)
}

// Private helper functions

func normalizeCreateExpenseRequest(req *CreateExpenseRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.Note = strings.TrimSpace(req.Note)
	req.ExpenseDate = strings.TrimSpace(req.ExpenseDate)
}

func validateCreateExpenseRequest(request CreateExpenseRequest) (string, error) {
	var validationEngine validation.Validation

	ok, err := validationEngine.Valid(&request)
	if err != nil {
		return "", err
	}

	if ok {
		_, err = time.Parse("2006-01-02", request.ExpenseDate)
		if err != nil {
			return "Expense date must be in YYYY-MM-DD format.", nil
		}

		if !models.IsValidCategory(request.Category) {
			return "Invalid category.", nil
		}

		if request.Amount <= 0 {
			return "Amount must be greater than 0.", nil
		}

		return "", nil
	}

	firstError := validationEngine.Errors[0]

	switch firstError.Key {
	case "Title":
		return mapExpenseTitleError(firstError), nil
	case "Amount":
		return mapExpenseAmountError(firstError), nil
	case "Category":
		return mapExpenseCategoryError(firstError), nil
	case "ExpenseDate":
		return mapExpenseDateError(firstError), nil
	default:
		return "Invalid request data.", nil
	}
}

func mapExpenseTitleError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Title is required."
	default:
		return "Invalid title."
	}
}

func mapExpenseAmountError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Amount is required."
	default:
		return "Invalid amount."
	}
}

func mapExpenseCategoryError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Category is required."
	default:
		return "Invalid category."
	}
}

func mapExpenseDateError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Expense date is required."
	default:
		return "Invalid expense date."
	}
}
