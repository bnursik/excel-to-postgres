package main

import (
	"log"
	"net/http"
	"os"

	"excel-to-postgres/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/upload", handlers.UploadExcelHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
