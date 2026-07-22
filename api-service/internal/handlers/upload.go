package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		SendError(w, http.StatusBadRequest, "File too large")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	// Определяем папку назначения (по умолчанию blogs)
	folder := r.FormValue("folder")
	if folder == "" {
		folder = "blogs"
	}

	// Create storage path
	basePath := os.Getenv("STORAGE_PATH")
	if basePath == "" {
		basePath = "./storage"
	}
	targetDir := filepath.Join(basePath, folder)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		SendError(w, http.StatusInternalServerError, "Could not create directory: "+err.Error())
		return
	}
	
	// Create a unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(targetDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Could not create file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		SendError(w, http.StatusInternalServerError, "Could not save file")
		return
	}

	// Return the static URL path
	urlPath := fmt.Sprintf("/static/%s/%s", folder, filename)
	SendJSON(w, http.StatusOK, map[string]string{
		"url": urlPath,
	})
}
