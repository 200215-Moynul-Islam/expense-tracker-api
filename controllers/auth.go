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

type AuthController struct {
	BaseController
}

type RegisterRequest struct {
	Name string `json:"name" valid:"Required"`
	Email string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required;MinSize(6)"`
}

func (c *AuthController) Register() {
	var req RegisterRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Warn("failed to parse register request body:", err)

		c.Error(http.StatusBadRequest, "Invalid request body.")
		return
	}

	normalizeRegisterRequest(&req)

	errMessage, err := validateRegisterRequest(req)
	if err != nil {
		logs.Error("failed to validate register request:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if errMessage != "" {
		c.Error(http.StatusBadRequest, errMessage)
		return
	}

	existingUser, err := models.GetUserByEmail(req.Email)
	if err != nil {
		logs.Error("failed to check existing user:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	if existingUser != nil {
		c.Error(http.StatusConflict, "Email already exists.")
		return
	}

	nextID, err := models.GetNextID()
	if err != nil {
		logs.Error("failed to generate next user id:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}
	user := models.User{
		ID: nextID,
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	err = models.CreateUser(user)
	if err != nil {
		logs.Error("failed to create user:", err)

		c.Error(http.StatusInternalServerError, "Internal server error.")
		return
	}

	c.Success(http.StatusCreated, "User registered successfully.", nil)
}

func normalizeRegisterRequest(req *RegisterRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
}

func validateRegisterRequest(request RegisterRequest) (string, error) {
	var validationEngine validation.Validation

	ok, err := validationEngine.Valid(&request)
	if err != nil {
		return "", err
	}

	if ok {
		return "", nil
	}

	firstError := validationEngine.Errors[0]

	switch firstError.Key {
	case "Name":
		return mapNameError(firstError), nil
	case "Email":
		return mapEmailError(firstError), nil
	case "Password":
		return mapPasswordError(firstError), nil
	default:
		return "Invalid request data.", nil
	}
}

func mapNameError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Name is required."
	default:
		return "Invalid name."
	}
}

func mapEmailError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Email is required."
	case "Email":
		return "Invalid email format."
	default:
		return "Invalid email."
	}
}

func mapPasswordError(validationError *validation.Error) string {
	switch validationError.Name {
	case "Required":
		return "Password is required."
	case "MinSize":
		return "Password must be at least 6 characters."
	default:
		return "Invalid password."
	}
}
