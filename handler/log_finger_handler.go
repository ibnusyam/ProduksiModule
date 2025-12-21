package handler

import (
	"Steril-App/internal/repository"
	"Steril-App/model"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type LogFingerHandler struct {
	Repo *repository.FingerLogRepository
}

func NewLogFingerHanlere(repo *repository.FingerLogRepository) *LogFingerHandler {
	return &LogFingerHandler{Repo: repo}
}

func (h *LogFingerHandler) GetFingerLog(c echo.Context) error {
	request := model.FingerLogRequest{}
	c.Bind(&request)
	result, err := h.Repo.GetFingerLog(request.Date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err,
		})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *LogFingerHandler) SaveNote(c echo.Context) error {
    // 1. Bind JSON dari Frontend
    request := model.NoteRequest{}
    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Format data tidak valid",
            "error":   err.Error(),
        })
    }

    // 2. Validasi Input
    if request.NIK == "" || request.Date == "" {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "NIK dan Tanggal wajib diisi",
        })
    }

    // 3. Panggil Repository
    err := h.Repo.SaveUserNote(request.NIK, request.Date, request.Note)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "Gagal menyimpan catatan",
            "error":   err.Error(),
        })
    }

    // 4. Sukses
    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Catatan berhasil disimpan",
        "data":    request,
    })
}

// GetNotes: Dipanggil saat load awal halaman atau saat ganti tanggal
// Contoh URL: GET /notes?date=2025-12-14
func (h *LogFingerHandler) GetNotes(c echo.Context) error {
    // 1. Ambil parameter date dari URL
    date := c.QueryParam("date")
    if date == "" {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Parameter 'date' diperlukan",
        })
    }

    // 2. Panggil Repository
    notes, err := h.Repo.GetNotesByDate(date)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "Gagal mengambil data catatan",
            "error":   err.Error(),
        })
    }

    // 3. Return JSON Array
    return c.JSON(http.StatusOK, notes)
}

func (h *LogFingerHandler) AddManualFingerLog(c echo.Context) error {
    // 1. Bind request
    request := model.AddManualFingerLogRequest{}
    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Format data tidak valid",
            "error":   err.Error(),
        })
    }
    fmt.Println(request)

    // 2. Validasi Input
    if request.NIK == "" || request.Timestamp == "" {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "NIK dan Timestamp harus diisi",
        })
    }

    // 3. Parsing String Waktu dengan LOKASI (Zona Waktu)
    layout := "2006-01-02 15:04:05"
    
    // Load lokasi Asia/Jakarta
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        // Fallback jika server tidak punya data tz (misal di docker alpine belum install tzdata)
        // Kita paksa offset +7
        loc = time.FixedZone("WIB", 7*60*60) 
    }

    // Gunakan ParseInLocation, bukan Parse biasa
    parsedTime, err := time.ParseInLocation(layout, request.Timestamp, loc)
    
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Format waktu salah. Gunakan format: YYYY-MM-DD HH:mm:ss",
            "error":   err.Error(),
        })
    }

    // 4. Panggil Repository Manual
    // Go sekarang tahu bahwa parsedTime adalah WIB. 
    // Saat dikirim ke Postgres, Go akan otomatis mengonversinya ke UTC yang benar 
    // (misal jam 08:00 WIB -> dikirim sebagai 01:00 UTC).
    err = h.Repo.AddManualFingerLog(request.NIK, parsedTime)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "Gagal menambahkan data manual",
            "error":   err.Error(),
        })
    }

    // 5. Sukses
    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Berhasil menambahkan data manual",
        "data": map[string]interface{}{
            "nik":       request.NIK,
            "timestamp": parsedTime, // Ini akan return string dengan offset +0700
        },
    })
}

func (h *LogFingerHandler) DeleteFingerLog(c echo.Context) error {
    // 1. Bind Request
    request := model.DeleteFingerLogRequest{}
    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Format data tidak valid",
            "error":   err.Error(),
        })
    }
    
    // Debug: Lihat data mentah yang dikirim frontend
    fmt.Printf("Request Delete -> NIK: %s | Timestamp String: %s\n", request.NIK, request.Timestamp)

    if request.NIK == "" || request.Timestamp == "" {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "NIK dan Timestamp harus diisi"})
    }

    // 2. PARSING WAKTU (MULTI-FORMAT)
    // Kita siapkan beberapa kemungkinan format agar fleksibel
    layouts := []string{
        "2006-01-02 15:04:05.000 -0700", // Format DETAIL (sesuai data kamu: 2025-11-27 17:18:23.838 +0700)
        "2006-01-02 15:04:05.000",       // Format detail tanpa timezone
        "2006-01-02 15:04:05",           // Format standar (detik 00)
        time.RFC3339,                    // Format ISO standar (2025-11-27T17:...)
    }

    var parsedTime time.Time
    parsedSuccess := false

    for _, layout := range layouts {
        // Coba parse
        t, err := time.Parse(layout, request.Timestamp)
        if err == nil {
            parsedTime = t
            parsedSuccess = true
            fmt.Println("Berhasil parse dengan format:", layout) // Debug
            break
        }
    }

    if !parsedSuccess {
        fmt.Println("Gagal parsing waktu. String:", request.Timestamp)
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "Format waktu tidak dikenali. Pastikan format YYYY-MM-DD HH:mm:ss.SSS +0700",
        })
    }

    // 3. Panggil Repository
    // Pastikan parsedTime ini memiliki presisi milidetik yang sama dengan DB
    err := h.Repo.DeleteFingerLog(request.NIK, parsedTime)
    if err != nil {
        fmt.Println("Repository Error:", err.Error()) // Debug di terminal
        
        // Sesuaikan pesan error ini dengan return dari repo kamu
        if err.Error() == "data tidak ditemukan atau sudah terhapus" {
            return c.JSON(http.StatusNotFound, map[string]interface{}{
                "message": "Data tidak ditemukan. Kemungkinan selisih milidetik.",
                "debug_time": parsedTime.String(),
            })
        }
        
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "Gagal menghapus data",
            "error":   err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Berhasil menghapus log finger",
    })
}