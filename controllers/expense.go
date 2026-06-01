package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
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

type UpdateExpenseRequest struct {
	Title string `json:"title" valid:"Required"`
	Amount float64 `json:"amount"`
	Category string `json:"category" valid:"Required"`
	Note string `json:"note"`
	ExpenseDate string `json:"expense_date" valid:"Required"`
}

type ExpenseListQueryParams struct {
	Category string
	DateFrom string
	DateTo string
	SortBy string
	SortOrder string
	Limit string
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

func (c *ExpenseController) GetAll() {
	userID := c.Ctx.Input.GetData("userID").(int)

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to get expenses:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	queryParams := ExpenseListQueryParams{
		Category: c.GetString("category"),
		DateFrom: c.GetString("date_from"),
		DateTo: c.GetString("date_to"),
		SortBy: c.GetString("sort_by"),
		SortOrder: c.GetString("sort_order"),
		Limit: c.GetString("limit"),
	}

	normalizeExpenseListQueryParams(&queryParams)

	errMessage, err := validateExpenseListQueryParams(queryParams)
	if err != nil {
		logs.Error("failed to validate expense list query params:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if errMessage != "" {
		c.Error(http.StatusBadRequest, errMessage)
		return
	}

	expenses = filterExpenses(expenses, queryParams)

	sortExpenses(expenses, queryParams)

	expenses, err = limitExpenses(expenses, queryParams)
	if err != nil {
		logs.Error("failed to limit expenses:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	c.Success(http.StatusOK, "Expenses retrieved successfully.", expenses)
}

func (c *ExpenseController) GetByID() {
	userID := c.Ctx.Input.GetData("userID").(int)

	idStr := c.Ctx.Input.Param(":id")
	if idStr == "" {
		c.Error(http.StatusBadRequest, "Expense id is required.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.Error(http.StatusBadRequest, "Invalid expense id.")
		return
	}

	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		logs.Error("failed to get expense by id:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if expense == nil {
		c.Error(http.StatusNotFound, "Expense not found.")
		return
	}

	c.Success(http.StatusOK, "Expense retrieved successfully.", expense)
}

func (c *ExpenseController) Update() {
	userID := c.Ctx.Input.GetData("userID").(int)

	idStr := c.Ctx.Input.Param(":id")
	if idStr == "" {
		c.Error(http.StatusBadRequest, "Expense id is required.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.Error(http.StatusBadRequest, "Invalid expense id.")
		return
	}

	existingExpense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		logs.Error("failed to get expense by id:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if existingExpense == nil {
		c.Error(http.StatusNotFound, "Expense not found.")
		return
	}

	var req UpdateExpenseRequest

	err = json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Warn("failed to parse update expense request body:", err)

		c.Error(http.StatusBadRequest, "Invalid request body.")
		return
	}

	normalizeUpdateExpenseRequest(&req)

	errMessage, err := validateUpdateExpenseRequest(req)
	if err != nil {
		logs.Error("failed to validate update expense request:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if errMessage != "" {
		c.Error(http.StatusBadRequest, errMessage)
		return
	}

	updatedExpense := models.Expense{
		ID: existingExpense.ID,
		UserID: existingExpense.UserID,
		Title: req.Title,
		Amount: req.Amount,
		Category: req.Category,
		Note: req.Note,
		ExpenseDate: req.ExpenseDate,
		CreatedAt: existingExpense.CreatedAt,
	}

	err = models.UpdateExpense(updatedExpense)
	if err != nil {
		logs.Error("failed to update expense:", err)

		if err.Error() == "expense not found" {
			c.Error(http.StatusNotFound, "Expense not found.")
			return
		}

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	c.Success(http.StatusOK, "Expense updated successfully.", updatedExpense)
}

func (c *ExpenseController) Delete() {
	userID := c.Ctx.Input.GetData("userID").(int)

	idStr := c.Ctx.Input.Param(":id")
	if idStr == "" {
		c.Error(http.StatusBadRequest, "Expense id is required.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.Error(http.StatusBadRequest, "Invalid expense id.")
		return
	}

	err = models.DeleteExpense(id, userID)
	if err != nil {
		logs.Error("failed to delete expense:", err)

		if err.Error() == "expense not found" {
			c.Error(http.StatusNotFound, "Expense not found.")
			return
		}

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	c.Success(http.StatusOK, "Expense deleted successfully.", nil)
}

// Private helper functions

func normalizeCreateExpenseRequest(req *CreateExpenseRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.Note = strings.TrimSpace(req.Note)
	req.ExpenseDate = strings.TrimSpace(req.ExpenseDate)
}

func normalizeUpdateExpenseRequest(req *UpdateExpenseRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.Note = strings.TrimSpace(req.Note)
	req.ExpenseDate = strings.TrimSpace(req.ExpenseDate)
}

func normalizeExpenseListQueryParams(params *ExpenseListQueryParams) {
	params.Category = strings.TrimSpace(params.Category)
	params.DateFrom = strings.TrimSpace(params.DateFrom)
	params.DateTo = strings.TrimSpace(params.DateTo)
	params.SortBy = strings.TrimSpace(params.SortBy)
	params.SortOrder = strings.TrimSpace(params.SortOrder)
	params.Limit = strings.TrimSpace(params.Limit)
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

func validateUpdateExpenseRequest(request UpdateExpenseRequest) (string, error) {
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

func validateExpenseListQueryParams(params ExpenseListQueryParams) (string, error) {
	var validationEngine validation.Validation

	_, err := validationEngine.Valid(&params)
	if err != nil {
		return "", err
	}

	if params.Category != "" {
		if !models.IsValidCategory(params.Category) {
			return "Invalid category.", nil
		}
	}

	if params.DateFrom != "" {
		_, err := time.Parse("2006-01-02", params.DateFrom)
		if err != nil {
			return "Invalid date_from format. Use YYYY-MM-DD.", nil
		}
	}

	if params.DateTo != "" {
		_, err := time.Parse("2006-01-02", params.DateTo)
		if err != nil {
			return "Invalid date_to format. Use YYYY-MM-DD.", nil
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		if params.DateTo < params.DateFrom {
			return "date_to cannot be earlier than date_from.", nil
		}
	}

	if params.SortBy != "" && params.SortBy != "amount" && params.SortBy != "expense_date" {
		return "Invalid sort_by value.", nil
	}

	if params.SortOrder != "" && params.SortOrder != "asc" && params.SortOrder != "desc" {
		return "Invalid sort_order value.", nil
	}

	if params.Limit != "" {
		limit, err := strconv.Atoi(params.Limit)
		if err != nil || limit <= 0 {
			return "Limit must be greater than 0.", nil
		}
	}

	return "", nil
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

func filterExpenses(expenses []models.Expense, params ExpenseListQueryParams,) []models.Expense {
	filteredExpenses := expenses

	// Filter by category
	if params.Category != "" {
		result := make([]models.Expense, 0)

		for _, expense := range filteredExpenses {
			if expense.Category == params.Category {
				result = append(result, expense)
			}
		}

		filteredExpenses = result
	}

	// Filter by date from
	if params.DateFrom != "" {
		result := make([]models.Expense, 0)

		for _, expense := range filteredExpenses {
			if expense.ExpenseDate >= params.DateFrom {
				result = append(result, expense)
			}
		}

		filteredExpenses = result
	}

	// Filter by date to
	if params.DateTo != "" {
		result := make([]models.Expense, 0)

		for _, expense := range filteredExpenses {
			if expense.ExpenseDate <= params.DateTo {
				result = append(result, expense)
			}
		}

		filteredExpenses = result
	}

	return filteredExpenses
}

func sortExpenses(expenses []models.Expense, params ExpenseListQueryParams) {
	sortOrder := params.SortOrder
	
	if params.SortBy == "" {
		return
	}

	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.Slice(expenses, func(i, j int) bool {
		switch params.SortBy {
		case "amount":
			if sortOrder == "asc" {
				return expenses[i].Amount < expenses[j].Amount
			}

			return expenses[i].Amount > expenses[j].Amount

		case "expense_date":
			if sortOrder == "asc" {
				return expenses[i].ExpenseDate < expenses[j].ExpenseDate
			}

			return expenses[i].ExpenseDate > expenses[j].ExpenseDate
		}

		return false
	})
}

func limitExpenses(expenses []models.Expense, params ExpenseListQueryParams) ([]models.Expense, error) {
	if params.Limit == "" {
		return expenses, nil
	}

	limit, err := strconv.Atoi(params.Limit)
	if err != nil {
		return nil, err
	}

	if limit < len(expenses) {
		expenses = expenses[:limit]
	}

	return expenses, nil
}
