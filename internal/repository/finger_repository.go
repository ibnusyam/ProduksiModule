package repository

import (
	"Steril-App/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type FingerRepository struct {
	DB *sql.DB
}

func NewFingerRepository(db *sql.DB) *FingerRepository {
	return &FingerRepository{DB: db}
}

func (repo *FingerRepository) FindEmptyFingerSlot() (string, error) {
	var emptySlotID int

	// Query SQL menggunakan GENERATE_SERIES dan EXCEPT untuk efisiensi
	query := `
        SELECT * FROM GENERATE_SERIES(1, 127)
        EXCEPT
        -- Karena finger_id di DB adalah VARCHAR, kita harus mengkonversinya ke INT untuk perbandingan
        SELECT CAST(finger_id AS INT) FROM fingerid
        ORDER BY 1
        LIMIT 1;
    `

	// QueryRow dan Scan untuk mendapatkan slot pertama
	err := repo.DB.QueryRow(query).Scan(&emptySlotID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Jika tidak ada error dan tidak ada baris, berarti semua slot (1-127) sudah terisi
			return "", fmt.Errorf("semua slot finger ID (1-127) sudah terisi")
		}
		// Error database lainnya
		return "", fmt.Errorf("gagal mencari slot kosong: %w", err)
	}

	// Mengkonversi ID slot (int) yang ditemukan menjadi string sebelum dikembalikan
	return strconv.Itoa(emptySlotID), nil
}

func (repo *FingerRepository) AddFingerData(nik string, fingerIDValue string) error {
	query := `INSERT INTO fingerid (nik, finger_id) VALUES ($1, $2)`
	result, err := repo.DB.Exec(query, nik, fingerIDValue)
	if err != nil {
		// Log error SQL yang mungkin terjadi (misal: NIK sudah ada/duplikat)
		fmt.Println("gagal")
		return fmt.Errorf("gagal insert data user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("gagal ini")
		log.Printf("Gagal mendapatkan jumlah baris terpengaruh: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Println("0")
		return fmt.Errorf("insert user gagal: 0 baris terpengaruh")
	}
	return nil
}

func (repo *FingerRepository) FindNikByID(id int) (string, error) {
	var nik string
	query := `SELECT nik FROM fingerid WHERE finger_id = $1`
	err := repo.DB.QueryRow(query, id).Scan(&nik)
	if err != nil {
		return "", fmt.Errorf("gagal menjalankan query :%w", err)
	}
	return nik, nil
}

func (repo *FingerRepository) GetDataFingerUser(nik string) error {
	query := `SELECT finger_id FROM fingerid where nik = $1`
	result, err := repo.DB.Query(query, nik)
	if err != nil {
		return fmt.Errorf("gagal melakukan pencarian data :%w", err)
	}

	var fingers []model.FingerID

	for result.Next() {
		var finger = model.FingerID{}
		result.Scan(&finger.FingerID)

		fingers = append(fingers, finger)
	}

	fmt.Println(fingers)
	return nil
}
