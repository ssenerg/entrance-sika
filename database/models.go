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
	Street  string    `json:"street"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	ZipCode string    `json:"zip_code"`
	Country string    `json:"country"`
	UserID  uuid.UUID `json:"user_id"`
}

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
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

func (db *Database) CreateUser(user User) error {
	
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	for i := range user.Addresses {
		user.Addresses[i].UserID = user.ID
	}
	result := db.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *Database) CreateAddress(address Address) error {
	result := db.DB.Create(&address)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *Database) CreateUsers(users []User) {
	ch := make(chan User, writeConcurrent)
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user User) {
			ch <- user
		}(user)
	}
	go func(users <-chan User, wg *sync.WaitGroup) {
		for user := range users {
			err := db.CreateUser(user)
			if err != nil {
				slog.Error("Error creating user", "error", err)
			}
			wg.Done()
		}
	}(ch, &wg)
	wg.Wait()
}

func (db *Database) GetUserByID(id uuid.UUID) (User, error) {
    var user User
    result := db.DB.Preload("Addresses").First(&user, "id = ?", id)
    if result.Error != nil {
        return user, result.Error
    }
    return user, nil
}
