package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertIntoPosgres(tableName string, rows [][]string, dropIfExists bool) error {
	if len(rows) < 2 {
		return fmt.Errorf("excel must have at least one header row and one data row")
	}

	db, err := ConnectDB()
	if err != nil {
		return nil
	}
	defer db.Close()

	headers := rows[0]

	// Drop the table if it exists
	checkQuery := `SELECT to_regclass($1);` // regclass returns table name if its exists
	var existingTable sql.NullString
	err = db.QueryRow(checkQuery, tableName).Scan(&existingTable)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}

	if existingTable.Valid && existingTable.String == tableName {
		if dropIfExists {
			_, err := db.Exec(fmt.Sprintf(`DROP TABLE "%s";`, tableName))
			if err != nil {
				return fmt.Errorf("failed to drop existing table: %w", err)
			}
		} else {
			return fmt.Errorf("table '%s' already exists (pass ?drop=true to replace it)", tableName)
		}
	}

	// Create table
	columnDefs := make([]string, len(headers))
	for i, col := range headers {
		columnDefs[i] = fmt.Sprintf("%q TEXT", col)
	}

	createQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %q (%s);`, tableName, strings.Join(columnDefs, ", "))
	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("table creation failed: %w", err)
	}

	placeholder := make([]string, len(headers))
	for i := range headers {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}
	inserQuery := fmt.Sprintf(`INSERT INTO %q (%s) VALUES (%s);`,
		tableName,
		strings.Join(quoteSlice(headers), ", "),
		strings.Join(placeholder, ", "),
	)

	for _, row := range rows[1:] {
		values := make([]interface{}, len(headers))
		for i := range headers {
			if i < len(row) {
				values[i] = row[i]
			} else {
				values[i] = ""
			}
		}

		_, err := db.Exec(inserQuery, values...)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}

	return nil
}

func quoteSlice(cols []string) []string {
	quoted := make([]string, len(cols))
	for i, c := range cols {
		quoted[i] = fmt.Sprintf(`"%s"`, c)
	}

	return quoted
}
