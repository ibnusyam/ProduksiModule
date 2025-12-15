package repository

import (
	"Steril-App/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type FingerLogRepository struct {
	DB *sql.DB
}

func NewFingerLogRepostory(db *sql.DB) *FingerLogRepository {
	return &FingerLogRepository{DB: db}
}

func (repo *FingerLogRepository) AddFingerLog(nik string) error {
	query := `INSERT INTO fingerlog (nik) VALUES ($1)`
	result, err := repo.DB.Exec(query, nik)
	if err != nil {
		// Log error SQL yang mungkin terjadi (misal: NIK sudah ada/duplikat)
		fmt.Println("gagal")
		return fmt.Errorf("gagal insert log finger: %w", err)
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

func (repo *FingerLogRepository) AddManualFingerLog(nik string, timestamp time.Time) error {
    // Kita insert NIK dan TIMESTAMP sesuai input
    query := `INSERT INTO fingerlog (nik, timestamp) VALUES ($1, $2)`
    
    result, err := repo.DB.Exec(query, nik, timestamp)
    if err != nil {
        log.Printf("Error insert manual: %v", err)
        return fmt.Errorf("gagal insert log finger manual: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("gagal cek rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("gagal insert: 0 baris terpengaruh")
    }
    
    return nil
}

func (repo *FingerLogRepository) GetFingerLog(date string) ([]model.FingerLogResult, error) {
	query := `SELECT
    u.nik,
    u.full_name,
    f.timestamp
    FROM
        fingerlog f
    JOIN
        users u ON f.nik = u.nik
    WHERE
        f.timestamp::date = $1
    ORDER BY
        u.nik ASC,
        f.timestamp ASC;`

	// Query ke database
	result, err := repo.DB.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer result.Close()

	// 1. Slice untuk menampung hasil akhir (Agar urutan terjaga)
	var data []model.FingerLogResult

	// 2. Map 'pembantu' untuk melacak posisi index berdasarkan NIK
	// Key: NIK, Value: Index di dalam slice 'data'
	indices := make(map[string]int)

	for result.Next() {
		var row model.RawFingerLog

		if err := result.Scan(&row.NIK, &row.FullName, &row.Timestamp); err != nil {
			return nil, fmt.Errorf("gagal scan data :%w", err)
		}

		// Cek apakah NIK ini sudah pernah kita masukkan ke slice 'data'?
		if idx, exists := indices[row.NIK]; exists {
			// KASUS: SUDAH ADA
			// Kita ambil data di slice berdasarkan index-nya, lalu append timestamp
			data[idx].Timestamps = append(data[idx].Timestamps, row.Timestamp)

			// fmt.Println("Append data ke:", row.NIK)
		} else {
			// KASUS: BELUM ADA (Data Baru)
			newEntry := model.FingerLogResult{
				NIK:        row.NIK,
				FullName:   row.FullName,
				Timestamps: []time.Time{row.Timestamp},
			}

			// Masukkan ke slice utama
			data = append(data, newEntry)

			// Catat posisi index-nya di map
			// len(data)-1 adalah index elemen yang baru saja kita masukkan
			indices[row.NIK] = len(data) - 1

			// fmt.Println("Buat entry baru:", row.NIK)
		}
	}

	// Tidak perlu looping map lagi disini, karena 'data' sudah terisi rapi & urut
	return data, nil
}

func (repo *FingerLogRepository) SaveUserNote(nik, date, note string) error {
    // Syntax PostgreSQL untuk UPSERT:
    // Jika kombinasi (nik, date) belum ada -> INSERT
    // Jika sudah ada (konflik) -> UPDATE kolom detail-nya saja
    query := `
        INSERT INTO detaillog (nik, date, detail)
        VALUES ($1, $2, $3)
        ON CONFLICT (nik, date) 
        DO UPDATE SET detail = EXCLUDED.detail;
    `

    _, err := repo.DB.Exec(query, nik, date, note)
    if err != nil {
        return fmt.Errorf("gagal menyimpan note ke database: %w", err)
    }
    return nil
}

// GetNotesByDate: Mengambil semua notes pada tanggal tertentu
func (repo *FingerLogRepository) GetNotesByDate(date string) ([]model.NoteResponse, error) {
    query := `SELECT nik, detail FROM detaillog WHERE date = $1`

    rows, err := repo.DB.Query(query, date)
    if err != nil {
        return nil, fmt.Errorf("gagal query notes: %w", err)
    }
    defer rows.Close()

    var notes []model.NoteResponse
    for rows.Next() {
        var n model.NoteResponse
        // Kita hanya butuh NIK dan isinya (detail) untuk mapping di Frontend
        if err := rows.Scan(&n.NIK, &n.Note); err != nil {
            return nil, fmt.Errorf("gagal scan row notes: %w", err)
        }
        notes = append(notes, n)
    }
    
    // Kembalikan slice kosong [] jika tidak ada data (bukan nil) agar JSON-nya "[]"
    if notes == nil {
        notes = []model.NoteResponse{}
    }

    return notes, nil
}

func (repo *FingerLogRepository) DeleteFingerLog(nik string, timestamp time.Time) error {
    // Query delete berdasarkan NIK dan Timestamp persis
    query := `DELETE FROM fingerlog WHERE nik = $1 AND timestamp = $2`

    result, err := repo.DB.Exec(query, nik, timestamp)
    if err != nil {
        log.Printf("Error deleting log: %v", err)
        return fmt.Errorf("gagal menghapus data: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("gagal cek rows affected: %v", err)
    }

    // Jika 0, berarti tidak ada data yang cocok (mungkin salah detik atau salah jam)
    if rowsAffected == 0 {
        return fmt.Errorf("data tidak ditemukan atau sudah terhapus")
    }

    return nil
}


