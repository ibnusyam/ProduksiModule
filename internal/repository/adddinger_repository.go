package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type AddFingerRepository struct {
	DB *sql.DB
}

func NewAddFingerRepository(db *sql.DB) *AddFingerRepository {
	return &AddFingerRepository{DB: db}
}

func (repo *AddFingerRepository) AddFingerByID(id int) (string, error) {
	var nik int
	getNikQuery := "SELECT nik FROM fingerid  where finger_id = $1"
	err := repo.DB.QueryRow(getNikQuery, id).Scan(&nik)
	if err != nil {
		if err == sql.ErrNoRows {
			// Ini terjadi jika finger_id tidak ditemukan
			log.Printf("Finger ID '%x' tidak ditemukan.", id)
			return "", fmt.Errorf("data tidak ditemukan")
		}
		// Ini terjadi jika ada error lain (koneksi, tipe data, dll.)
		fmt.Println(err)
		return "", fmt.Errorf("gagal menjalankan query atau scan: %w", err)
	}
	fmt.Println(nik)

	query := "INSERT INTO fingerlog (nik) VALUES ($1)"
	result, err := repo.DB.Exec(query, nik)
	if err != nil {
		// Log error SQL yang mungkin terjadi (misal: NIK sudah ada/duplikat)
		fmt.Println("gagal")
		return "", fmt.Errorf("gagal insert data user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("gagal ini")
		log.Printf("Gagal mendapatkan jumlah baris terpengaruh: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Println("0")
		return "", fmt.Errorf("insert user gagal: 0 baris terpengaruh")
	}

	fmt.Printf("User dengan NIK % berhasil ditambahkan.\n", nik)
	return "berhasil", nil
}
