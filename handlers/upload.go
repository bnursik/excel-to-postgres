package handlers

import (
	"log"
	"net/http"

	"excel-to-postgres/utils"
)

func UploadExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the tableName
	tableName := r.URL.Query().Get("table")
	if tableName == "" {
		http.Error(w, "Missing 'table' query parameter", http.StatusBadRequest)
		return
	}

	// Checks drop=true flag
	dropFlag := r.URL.Query().Get("drop") == "true"

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	rows, err := utils.ParseExcel(file)
	if err != nil {
		http.Error(w, "Failed to parse Excel"+err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.InsertIntoPosgres(tableName, rows, dropFlag)
	if err != nil {
		http.Error(w, "Failed to insert into DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Excel file parced and inserted into PostgreSQL!")
}
