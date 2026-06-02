# Expense Tracker API

A RESTful Personal Expense Tracker API built with Go and Beego v2 using CSV file storage.

Repository: https://github.com/200215-Moynul-Islam/expense-tracker-api

---

## Features

- User registration and authentication
- Expense CRUD operations
- Expense filtering and summary reporting
- CSV-based persistent storage
- Middleware-based request authentication
- Swagger API documentation
- Table-driven unit tests
- Test coverage reporting support

---

## Technology Stack

- Go
- Beego v2
- CSV File Storage
- Swagger (auto-generated via Beego)

---

## Requirements

Before running the project, make sure you have:

- Go 1.22+
- Git
- Bee CLI (optional for Swagger generation)

Install Bee CLI:

```bash
go install github.com/beego/bee/v2@latest
```

---

## Quick Setup

### 1. Clone Repository

```bash
git clone https://github.com/200215-Moynul-Islam/expense-tracker-api.git
cd expense-tracker-api
```

### 2. Create Configuration File

Copy the example configuration:

```bash
cp conf/app.conf.example conf/app.conf
```

Example configuration:

```ini
appname = expense-tracker-api
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true
```

> **Note:** This README uses the hardcoded port `8080` everywhere for quick setup and testing. If you change the `httpport` value inside `conf/app.conf`, replace `8080` in all URLs and API examples with your configured port number.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
bee run
```

Application will start at:

```text
http://localhost:8080
```

### 5. Run Tests

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

Generate detailed coverage profile:

```bash
go test ./... -coverprofile=coverage.out
```

View function-level coverage:

```bash
go tool cover -func=coverage.out
```

Open HTML coverage report:

```bash
go tool cover -html=coverage.out
```

### Current Coverage Report

Example output from the current project test suite:

```text
expense-tracker-api             coverage: 0.0% of statements
ok      expense-tracker-api/controllers (cached)        coverage: 37.0% of statements
ok      expense-tracker-api/middlewares (cached)        coverage: 100.0% of statements
ok      expense-tracker-api/models      (cached)        coverage: 85.4% of statements
expense-tracker-api/routers             coverage: 0.0% of statements
ok      expense-tracker-api/utils       (cached)        coverage: 83.3% of statements
```

---

## Swagger Documentation

Generate Swagger documentation:

```bash
bee run -gendoc=true -downdoc=true
```

Swagger UI will be available at:

```text
http://localhost:8080/swagger/
```

Note:

- Swagger files are generated locally

---

## Authentication

Protected endpoints require the following header:

```http
X-User-ID: <user-id>
```

---

## API Base URL

```text
http://localhost:8080/api/v1
```

---

## Allowed Expense Categories

- Food
- Transport
- Housing
- Entertainment
- Shopping
- Healthcare
- Education
- Utilities
- Other

---

## Project Structure

```text
.
├── conf/
│   └── app.conf.example
│
├── controllers/
│   ├── auth.go
│   ├── auth_test.go
│   ├── expense.go
│   ├── expense_test.go
│   ├── health.go
│   ├── health_test.go
│   ├── base.go
│   └── base_test.go
│
├── data/
│   └── .gitkeep
│
├── middlewares/
│   ├── auth.go
│   └── auth_test.go
│
├── models/
│   ├── user.go
│   ├── user_test.go
│   ├── expense.go
│   └── expense_test.go
│
├── routers/
│   └── router.go
│
├── utils/
│   ├── csv.go
│   └── csv_test.go
│
├── main.go
├── go.mod
├── go.sum
└── README.md
```

---

## Quick API Testing

All examples below use port `8080`.

---

### Health Check

```bash
curl --location 'http://localhost:8080/api/v1/health'
```

---

### Register User

```bash
curl --location 'http://localhost:8080/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data '{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}'
```

---

### Login User

```bash
curl --location 'http://localhost:8080/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data '{
  "email": "john@example.com",
  "password": "password123"
}'
```

---

### Create Expense

```bash
curl --location 'http://localhost:8080/api/v1/expenses' \
--header 'Content-Type: application/json' \
--header 'X-User-ID: 1' \
--data '{
  "title": "Lunch",
  "amount": 250,
  "category": "Food",
  "note": "Office lunch",
  "expense_date": "2026-06-02"
}'
```

---

### Get All Expenses

```bash
curl --location 'http://localhost:8080/api/v1/expenses' \
--header 'X-User-ID: 1'
```

---

### Get Expense By ID

```bash
curl --location 'http://localhost:8080/api/v1/expenses/1' \
--header 'X-User-ID: 1'
```

---

### Update Expense

```bash
curl --location --request PUT 'http://localhost:8080/api/v1/expenses/1' \
--header 'Content-Type: application/json' \
--header 'X-User-ID: 1' \
--data '{
  "title": "Updated Lunch",
  "amount": 320,
  "category": "Food",
  "note": "Updated note",
  "expense_date": "2026-06-02"
}'
```

---

### Delete Expense

```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/expenses/1' \
--header 'X-User-ID: 1'
```

---

### Expense Summary

```bash
curl --location 'http://localhost:8080/api/v1/expenses/summary?date_from=2026-06-01&date_to=2026-06-30' \
--header 'X-User-ID: 1'
```

---

## Notes

- CSV data files are generated automatically at runtime inside the `data/` directory
- The repository includes a `data/.gitkeep` file to preserve the directory structure
- `users.csv` and `expenses.csv` are not included in the repository
- Swagger files are generated locally and are not committed
- Designed for learning clean backend architecture with Go and Beego
