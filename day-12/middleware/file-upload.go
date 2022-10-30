package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// HANDLING UPLOAD IMAGE
func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("upload-image")
		if handler == nil {
			ctx := context.WithValue(r.Context(), "dataFile", "empty")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error retrieving the File")
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)

		tempFile, err := os.CreateTemp("uploads", "image-*"+handler.Filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error.")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		filename := data[8:]

		ctx := context.WithValue(r.Context(), "dataFile", filename)
		fmt.Println(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
