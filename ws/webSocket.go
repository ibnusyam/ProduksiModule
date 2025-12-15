package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync" // PENTING: Import ini untuk Mutex

	"Steril-App/internal/repository" // Sesuaikan import path
	"Steril-App/model"               // Sesuaikan import path

	// "Steril-App/ws"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	RepoFingerSocket *repository.AddFingerRepository
	RepoFinger       *repository.FingerRepository
	RepoLogFinger    *repository.FingerLogRepository
}

func NewWebSocketHandler(repoFingerSocket *repository.AddFingerRepository, repoFinger *repository.FingerRepository, repoLogFinger *repository.FingerLogRepository) *WebSocketHandler {
	return &WebSocketHandler{
		RepoFingerSocket: repoFingerSocket,
		RepoFinger:       repoFinger,
		RepoLogFinger:    repoLogFinger,
	}
}

// --- BAGIAN INI DIPERBAIKI ---

// Variabel global dengan pengamanan Mutex
var (
	nodeMcuConn *websocket.Conn
	connMutex   sync.Mutex // Ini adalah "Gembok/Kunci" pengaman
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket: Endpoint untuk NodeMCU connect (/ws)
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("❌ Error upgrade:", err)
		return err
	}

	// --- PENGAMANAN KONEKSI BARU ---
	connMutex.Lock() // Kunci pintu dulu
	if nodeMcuConn != nil {
		// Jika ada koneksi lama yang nyangkut, tutup paksa biar gak bentrok
		nodeMcuConn.Close()
	}
	nodeMcuConn = conn // Simpan koneksi baru
	connMutex.Unlock() // Buka kunci
	// -------------------------------

	log.Println("✅ NodeMCU Terhubung ke Backend!")

	// Pastikan koneksi ditutup bersih saat fungsi selesai
	defer func() {
		connMutex.Lock()
		if nodeMcuConn == conn { // Cek apakah yang ditutup adalah koneksi ini
			nodeMcuConn = nil
		}
		connMutex.Unlock()
		conn.Close()
		log.Println("⚠️ Koneksi WebSocket Ditutup/Dibersihkan")
	}()

	// Loop mendengarkan pesan
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Error umum saat putus adalah "websocket: close 1006 (abnormal closure)"
			log.Println("❌ NodeMCU Terputus/Read Error:", err)
			break // Keluar dari loop -> akan memicu defer di atas
		}

		log.Printf("📩 Pesan Masuk: %s\n", message)

		// Proses JSON
		var response model.SensorResponse // Pastikan pakai struct yang benar
		if err := json.Unmarshal(message, &response); err != nil {
			log.Println("⚠️ Gagal parsing JSON:", err)
			continue
		}

		// Logika Bisnis
		if response.Action == "ABSENSI" {
			// Panggil Repository
			// Pastikan method di repo kamu benar menerima parameter yg sesuai
			// Contoh: h.Repo.AddFingerByID(response.ID)
			fmt.Println("masuk coy")
			// fmt.Println(response)
			nik, err := h.RepoFinger.FindNikByID(response.ID)
			if err != nil {
				log.Println("❌ gagal cari NIK:", err)
				continue
			}

			err = h.RepoLogFinger.AddFingerLog(nik)
			if err != nil {
				log.Println("❌ gagal tambah log absensi:", err)
				continue
			}

			fmt.Println("✅ Data Absensi Diterima untuk ID:", response.ID)

		} else {
			fmt.Println("ℹ️ Action lain:", response.Action)
		}
	}

	return nil
}

// SendCommand: Fungsi bantuan untuk mengirim data JSON ke NodeMCU
func SendCommand(data interface{}) error {
	// --- PENGAMANAN SAAT MENGIRIM ---
	connMutex.Lock()
	defer connMutex.Unlock() // Pastikan kunci dibuka setelah fungsi selesai

	if nodeMcuConn == nil {
		// Ini terjadi kalau NodeMCU mati/putus
		return echo.NewHTTPError(http.StatusServiceUnavailable, "NodeMCU belum terhubung!")
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Gagal marshal JSON: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memproses data JSON")
	}

	log.Printf("🔥 Mengirim ke NodeMCU: %s", string(jsonBytes))

	// Kirim pesan
	err = nodeMcuConn.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		log.Printf("❌ Gagal kirim pesan: %v", err)
		// Jika gagal kirim, kita anggap koneksi rusak
		nodeMcuConn.Close()
		nodeMcuConn = nil
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal kirim perintah (Koneksi Putus)")
	}

	return nil
}
