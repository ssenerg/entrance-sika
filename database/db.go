package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Config
	DB *gorm.DB
}

type Config struct {
	host       string
	user       string
	password   string
	dbName     string
	port       string
	sslMode    string
	timezone   string
	gormConfig *gorm.Config
}

func NewConfig(host, user, password, dbname, port, sslmode, timezone string, gormConfig *gorm.Config) (*Config, error) {
	if host == "" {
		host = "localhost"
	}
	if user == "" {
		return nil, errors.New("user is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if dbname == "" {
		return nil, errors.New("dbname is required")
	}
	if port == "" {
		port = "5432"
	}

	switch sslmode {
	case "":
		sslmode = "disable"
	case "disable", "require":
	default:
		return nil, errors.New("sslmode must be either disable or require")
	}

	if timezone == "" {
		timezone = "UTC"
	} else if _, err := time.LoadLocation(timezone); err != nil {
		return nil, err
	}

	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}
	
	return &Config{
		host:       host,
		user:       user,
		password:   password,
		dbName:     dbname,
		port:       port,
		sslMode:    sslmode,
		timezone:   timezone,
		gormConfig: gormConfig,
	}, nil
}

func NewDatabase(config *Config) (*Database, error) {
	db, err := openDB(config)
	if err != nil {
		return nil, err
	}
	database := &Database{Config: *config, DB: db}

	if err := database.autoMigrate(); err != nil {
		return nil, err
	}
	return database, nil
}

func (c *Config) dsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", c.host, c.user, c.password, c.dbName, c.port, c.sslMode, c.timezone)
}

func openDB(config *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.dsn()), config.gormConfig)
}

func (db *Database) autoMigrate() error {
	return db.DB.AutoMigrate(&User{}, &Address{})
}
