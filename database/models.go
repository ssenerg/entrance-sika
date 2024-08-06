package database

import (
	"encoding/json"
	"os"
	"sync"

	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const writeConcurrent int = 10

type Address struct {
	gorm.Model
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
	UserID  string `json:"user_id"` // Foreign key to associate with User
}

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid" json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Addresses   []Address `json:"addresses" gorm:"foreignKey:UserID"`
}

func ReadUserFromJson(filename string) ([]User, error) {

	var users []User
	file, err := os.Open(filename)
	if err != nil {
		return users, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (db *Database) CreateUser(user User) {
	if result := db.DB.Create(&user).Error; result != nil {
		slog.Error("Error while creating user", "error", result, "user", user)
	}
}

func (db *Database) CreateUsers(users []User) {
	ch := make(chan User, writeConcurrent)
	var wg *sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user User) {
			ch <- user
		}(user)
	}
	go func(users <-chan User, wg *sync.WaitGroup) {
		for user := range users {
			db.CreateUser(user)
			wg.Done()
		}
	}(ch, wg)
	wg.Wait()
}

func (db *Database) GetUserByID(id uuid.UUID) (User, error) {
	var user User
	result := db.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
