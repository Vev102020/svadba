package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)


const (
	uploadDir = "./uploads"
	questDir  = "./quest"
	siGameDir  = "./sigame"
	port      = ":8080"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Type    string    `json:"type"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

func main() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(questDir, 0755); err != nil {
		log.Fatal(err)
	}
	
	if err := os.MkdirAll(siGameDir, 0755); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Статика (собранный Vite)
	mux.Handle("/static/", spaHandler("static"))
	// Загруженные файлы
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))
	mux.Handle("/quest/", http.StripPrefix("/quest/", http.FileServer(http.Dir(questDir))))

	// API
	mux.HandleFunc("/sigame/", sigameHandler)
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/files", filesHandler)
	mux.HandleFunc("/api/upload", uploadHandler)
	mux.HandleFunc("/api/quest/files", questFilesHandler)
	mux.HandleFunc("/api/quest/upload", questUploadHandler)
	mux.HandleFunc("/api/sigame/", sigameAPIHandler)


	log.Println("╔════════════════════════════════════════════════════╗")
	log.Printf("║  Сервер: http://0.0.0.0%s                      ║", port)
	log.Println("║  Откройте этот URL на другом устройстве по Wi-Fi  ║")
	log.Println("╚════════════════════════════════════════════════════╝")
	log.Fatal(http.ListenAndServe("0.0.0.0"+port, mux))
}


func spaHandler(staticDir string) http.Handler {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.TrimPrefix(r.URL.Path, "/static/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		fullPath := filepath.Join(staticDir, cleanPath)
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})
}


func sigameHandler(w http.ResponseWriter, r *http.Request) {
   	path := r.URL.Path

	// 1. Файлы из rounds (JSON, картинки, видео) — отдаем как статику
	if strings.HasPrefix(path, "/sigame/rounds/") {
		fs := http.StripPrefix("/sigame/", http.FileServer(http.Dir(siGameDir)))
		fs.ServeHTTP(w, r)
		return
	}

	// 2. Любой другой путь /sigame, /sigame/edit — отдаем index.html для Vue Router
	http.ServeFile(w, r, "./static/index.html")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./static/index.html")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		name := info.Name()
		ext := strings.ToLower(filepath.Ext(name))
		fileType := "other"
		if isImage(ext) {
			fileType = "image"
		} else if isVideo(ext) {
			fileType = "video"
		}

		files = append(files, FileInfo{
			Name:    name,
			Size:    info.Size(),
			Type:    fileType,
			URL:     "/uploads/" + name,
			Created: info.ModTime(),
		})
	}

	// Новые файлы первыми
	sort.Slice(files, func(i, j int) bool {
		return files[i].Created.After(files[j].Created)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму. Файлы > 32MB пишутся во временный файл на диск, а не в RAM
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Ошибка при разборе формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "Файлы не загружены", http.StatusBadRequest)
		return
	}

	var uploaded []FileInfo

	for _, fileHeader := range files {
		func() {
			// Открываем входящий файл (поток)
			file, err := fileHeader.Open()
			if err != nil {
				return
			}
			defer file.Close()

			// Генерируем безопасное уникальное имя
			name := filepath.Base(fileHeader.Filename)
			ext := filepath.Ext(name)
			base := strings.TrimSuffix(name, ext)
			safeName := fmt.Sprintf("%s_%d%s", base, time.Now().UnixMilli(), ext)
			safeName = strings.ReplaceAll(safeName, " ", "_")

			dstPath := filepath.Join(uploadDir, safeName)

			// Создаём файл на диске и сразу туда пишем через io.Copy
			dst, err := os.Create(dstPath)
			if err != nil {
				return
			}
			defer dst.Close()

			if _, err = io.Copy(dst, file); err != nil {
				os.Remove(dstPath)
				return
			}

			info, _ := os.Stat(dstPath)
			fileType := "other"
			lowerExt := strings.ToLower(ext)
			if isImage(lowerExt) {
				fileType = "image"
			} else if isVideo(lowerExt) {
				fileType = "video"
			}

			uploaded = append(uploaded, FileInfo{
				Name:    safeName,
				Size:    info.Size(),
				Type:    fileType,
				URL:     "/uploads/" + safeName,
				Created: info.ModTime(),
			})
		}()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploaded)
}

func isImage(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	}
	return false
}

func isVideo(ext string) bool {
	switch ext {
	case ".mp4", ".webm", ".mov", ".avi", ".mkv", ".m4v":
		return true
	}
	return false
}


func questFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	entries, err := os.ReadDir(questDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := []FileInfo{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		name := info.Name()
		ext := strings.ToLower(filepath.Ext(name))
		fileType := "other"
		if isImage(ext) {
			fileType = "image"
		} else if isVideo(ext) {
			fileType = "video"
		}
		files = append(files, FileInfo{
			Name:    name,
			Size:    info.Size(),
			Type:    fileType,
			URL:     "/quest/" + name,
			Created: info.ModTime(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func questUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	questID := r.URL.Query().Get("id")
	if questID == "" {
		http.Error(w, "Отсутствует id", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		http.Error(w, "Нет файлов", http.StatusBadRequest)
		return
	}
	fh := fileHeaders[0]
	file, err := fh.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	
	safeName := questID + ext
	dstPath := filepath.Join(questDir, safeName)

	// Удаляем старый файл, если есть
	os.Remove(dstPath)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, _ := os.Stat(dstPath)
	fileType := "other"
	if isImage(ext) {
		fileType = "image"
	} else if isVideo(ext) {
		fileType = "video"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FileInfo{
		Name:    safeName,
		Size:    info.Size(),
		Type:    fileType,
		URL:     "/quest/" + safeName,
		Created: info.ModTime(),
	})
}

// SiGame
type PackMeta struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func sigamePacksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	roundsDir := filepath.Join(siGameDir, "rounds")
	os.MkdirAll(roundsDir, 0755)

	entries, err := os.ReadDir(roundsDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	packs := []PackMeta{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		packID := entry.Name()
		dataPath := filepath.Join(roundsDir, packID, "data.json")
		data, err := os.ReadFile(dataPath)
		if err != nil {
			// Папка есть, но data.json ещё не создан — добавляем с пустым именем
			packs = append(packs, PackMeta{ID: packID, Name: ""})
			continue
		}

		var pack struct {
			Name string `json:"name"`
		}
		json.Unmarshal(data, &pack)
		packs = append(packs, PackMeta{ID: packID, Name: pack.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packs)
}

// /api/sigame/pack/{id} — GET (загрузить), POST (создать/обновить), DELETE (удалить)
func sigamePackHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем id из пути: /api/sigame/pack/pack_abc123
	path := strings.TrimPrefix(r.URL.Path, "/api/sigame/pack/")
	packID := filepath.Base(path)
	if packID == "" || packID == "." || packID == "/" {
		http.Error(w, "Не указан ID пака", http.StatusBadRequest)
		return
	}

	packDir := filepath.Join(siGameDir, "rounds", packID)

	switch r.Method {
	case http.MethodGet:
		dataPath := filepath.Join(packDir, "data.json")
		data, err := os.ReadFile(dataPath)
		if err != nil {
			http.Error(w, "Пак не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case http.MethodPost:
		// Читаем JSON из тела запроса и сохраняем в data.json
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		os.MkdirAll(packDir, 0755)
		dataPath := filepath.Join(packDir, "data.json")
		if err := os.WriteFile(dataPath, body, 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))

	case http.MethodDelete:
		os.RemoveAll(packDir)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))

	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// POST /api/sigame/upload?id=pack_xxx — загрузка медиафайла в папку пака
func sigameUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	packID := r.URL.Query().Get("id")
	if packID == "" {
		http.Error(w, "Отсутствует id пака", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		http.Error(w, "Нет файлов", http.StatusBadRequest)
		return
	}

	fh := fileHeaders[0]
	file, err := fh.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	packDir := filepath.Join(siGameDir, "rounds", packID)
	os.MkdirAll(packDir, 0755)

	safeName := filepath.Base(fh.Filename)
	safeName = strings.ReplaceAll(safeName, " ", "_")
	dstPath := filepath.Join(packDir, safeName)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url":  "/sigame/rounds/" + packID + "/" + safeName,
		"name": safeName,
	})
}

func sigameAPIHandler(w http.ResponseWriter, r *http.Request) {
	// CORS (обязательно для Vue + Go на разных портах)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/sigame")

	switch path {
	case "/packs":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		sigamePacksHandler(w, r)

	case "/upload":
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		sigameUploadHandler(w, r)

	default:
		// Обрабатываем /pack/{id}
		if strings.HasPrefix(path, "/pack/") {
			sigamePackHandler(w, r) // внутри уже извлекает ID из path
			return
		}

		http.NotFound(w, r)
	}
}