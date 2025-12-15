package repository

import (
	"Steril-App/model"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func GetDSN() (string, error) {
	config := model.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"), // Tambahkan ini di .env: DB_SSLMODE=disable
	}

	if config.Host == "" || config.User == "" {
		return "", fmt.Errorf("variabel lingkungan database tidak lengkap")
	}

	// Format DSN PostgreSQL: "host=... port=... user=... password=... dbname=... sslmode=..."
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	return dsn, nil
}

func ConnectDB() (*sql.DB, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan DSN: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal melakukan ping database: %w", err)
	}

	log.Println("Koneksi ke PostgreSQL berhasil!")

	return db, nil
}
