package handlers

import (
	"fmt"
	"net/http"

	"exel-to-postgres/utils"
)

func UploadExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

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

	err = utils.InsertIntoPosgres("uploaded_data", rows)
	if err != nil {
		http.Error(w, "Failed to insert into DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(w, "Excel file parced and inserted into PostgreSQL!")
}
