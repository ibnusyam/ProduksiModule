package repository

import (
	"Steril-App/model"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(data *model.CreateUserRequest) error {
	query := "INSERT INTO users (nik, full_name) VALUES ($1, $2)"
	fmt.Println(data.NIK)
	result, err := repo.DB.Exec(query, data.NIK, data.FullName)
	if err != nil {
		// Log error SQL di sini.
		log.Printf("ERROR SQL: Gagal insert user %s: %v", data.NIK, err)

		// Kembalikan error yang lebih spesifik jika perlu (misal: duplikasi key)
		return fmt.Errorf("gagal menambahkan user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ERROR DB Driver: Gagal mendapatkan jumlah baris terpengaruh setelah INSERT: %v", err)
		return fmt.Errorf("verifikasi insert gagal: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("WARNING: Insert user %s berhasil dieksekusi, tetapi 0 baris terpengaruh.", data.NIK)
		return fmt.Errorf("insert user gagal: 0 baris terpengaruh")
	}

	log.Printf("User %s berhasil ditambahkan. (%d baris terpengaruh)", data.NIK, rowsAffected)
	return nil
}

func (repo *UserRepository) IsUserExist(data *model.CreateUserRequest) (string, error) {
	var nik string
	query := "SELECT nik FROM users where nik = $1"
	err := repo.DB.QueryRow(query, data.NIK).Scan(&nik)
	// fmt.Println(nik)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return "exist", nil
}

func (repo *UserRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := repo.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal menjalankan query DELETE: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user dengan ID %s tidak ditemukan", id)
	}

	return nil
}

func (repo *UserRepository) GetAllUser() ([]model.UserResponse, error) {
	query := `SELECT id, nik, full_name FROM users`
	result, err := repo.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data dari database :%w", err)
	}

	var rows []model.UserResponse

	for result.Next() {
		var row = model.UserResponse{}
		result.Scan(&row.ID, &row.NIK, &row.FullName)

		rows = append(rows, row)
	}

	return rows, nil

}
