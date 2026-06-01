package models

import (
	"encoding/csv"
	"os"
	"strconv"

	logs "github.com/beego/beego/v2/core/logs"

	"expense-tracker-api/utils"
)

const (
	totalUserFields = 5
	userCSVPath = "data/users.csv"
	userFilePermission = 0644
)

var userCSVHeader = []string{
	"id",
	"name",
	"email",
	"password",
	"created_at",
}

type User struct {
	ID int
	Name string
	Email string
	Password string
	CreatedAt string
}

func GetAllUsers() ([]User, error) {
	if err := utils.EnsureCSVFile(userCSVPath, userCSVHeader); err != nil {
		return nil, err
	}

	file, err := os.Open(userCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(records))

	for i := 1; i < len(records); i++ {
		// An user row should have at least 5 columns (id, name, email, password, created_at).
		if len(records[i]) < totalUserFields {
			// Log an error for the invalid user and skip this record.
			// We don't want to stop processing all users just because one record is bad.
			logs.Warn("invalid user row:", records[i])
			continue
		}

		id, err := strconv.Atoi(records[i][0])
		if err != nil {
			// Log an error for the invalid user ID and skip this record.
			// We don't want to stop processing all users just because one record is bad.
			logs.Warn("invalid user id:", records[i])
			continue
		}

		users = append(users, User{
			ID: id,
			Name: records[i][1],
			Email: records[i][2],
			Password: records[i][3],
			CreatedAt: records[i][4],
		})
	}

	return users, nil
}

func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, nil
}

func GetNextID() (int, error) {
	users, err := GetAllUsers()
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, u := range users {
		maxID = max(maxID, u.ID)
	}

	return maxID + 1, nil
}

func CreateUser(user User) error {
	if err := utils.EnsureCSVFile(userCSVPath, userCSVHeader); err != nil {
		return err
	}

	file, err := os.OpenFile(userCSVPath, os.O_APPEND|os.O_WRONLY, userFilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write([]string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	})

	if err != nil {
		return err
	}

	writer.Flush()

	return writer.Error()
}