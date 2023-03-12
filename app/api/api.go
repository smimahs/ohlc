package api

import (
	"database/sql"
	"fmt"
	"os"

	"app/config"
	pq "github.com/lib/pq"

	"bufio"
	"encoding/csv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"strconv"
	"strings"
	//"net/http"
)

func Connect() (*sql.DB, error) {
	// Get database connection string from environment variable
	config.Load()
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Open database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected!")
	return db, nil
}

func Update_API(app *iris.Application, db *sql.DB) {
	// Define a handler for POST requests to "/update"
	app.Handle("POST", "/update", func(ctx context.Context) {
		// Define a struct to hold the JSON parameters
		var params struct {
			URL string `json:"url"`
		}
		// Parse the JSON request body into the params struct
		err := ctx.ReadJSON(&params)
		if err != nil {
			// If there was an error parsing the request, return an error response
			ctx.JSON(context.Map{"response": err.Error()})
		} else {
			// // Create the OHLCData table if it doesn't already exist
			_, err = db.Exec(`
            CREATE TABLE IF NOT EXISTS OHLCData (
                id SERIAL PRIMARY KEY,
                UNIX VARCHAR(100) NOT NULL,
                SYMBOL VARCHAR(50) NOT NULL,
                OPEN VARCHAR(50) NOT NULL,
                HIGH VARCHAR(50) NOT NULL,
                LOW VARCHAR(50) NOT NULL,
                CLOSE VARCHAR(50) NOT NULL
            )
            `)
			if err != nil {
				// If there was an error parsing the request, return an error response
				ctx.JSON(context.Map{"response": err.Error()})
			}
			// Begin a database transaction
			tx, err := db.Begin()
			if err != nil {
				// If there was an error parsing the request, return an error response
				ctx.JSON(context.Map{"response": err.Error()})
			}
			defer tx.Rollback()
			// Prepare the bulk insert statement
			// Prepare a bulk insert statement using the "COPY FROM STDIN" syntax provided by the pq driver
			stmt, err := tx.Prepare(pq.CopyIn("ohlcdata", "unix", "symbol", "open", "high", "low", "close"))
			if err != nil {
				// If there was an error parsing the request, return an error response
				ctx.JSON(context.Map{"response": err.Error()})
			}

			// Open the CSV file specified by the "url" parameter
			// if csv get from url uncomment this part
			/*resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()
			reader := csv.NewReader(resp.Body)
			reader.Comma = ','
			*/

			// Open the CSV file specified by the "url" parameter
			// if csv get from a file name inside the project use this part
			file, err := os.Open(params.URL)
			if err != nil {
				// If there was an error parsing the request, return an error response
				ctx.JSON(context.Map{"response": err.Error()})
			}
			defer file.Close()
			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ','
			reader.Comment = '#'

			// Read the header row
			_, err = reader.Read()
			if err != nil {
				// If there was an error parsing the request, return an error response
				ctx.JSON(context.Map{"response": err.Error()})
			}

			var values []interface{}
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					ctx.JSON(context.Map{"response": err.Error()})
				}
				// Convert the row to a slice of interface{} values
				for _, col := range row {
					values = append(values, strings.TrimSpace(col))
				}

				// Write the row to the CopyInReader
				_, err = stmt.Exec(values...)
				if err != nil {
					ctx.JSON(context.Map{"response": err.Error()})
				}
				// Reset the values slice for the next row
				values = nil
			}
			// Close the statement to flush the data to the database
			_, err = stmt.Exec()
			if err != nil {
				ctx.JSON(context.Map{"response": err.Error()})
			}
			err = tx.Commit()
			if err != nil {
				ctx.JSON(context.Map{"response": err.Error()})
			}
			ctx.JSON(context.Map{"response": "Bulk insert from CSV successful"})
			fmt.Println("Bulk insert from CSV successful")

		}

	})
}

func Query_API(app *iris.Application, db *sql.DB) {
	app.Handle("GET", "/query", func(ctx context.Context) {
		// Extract the query string parameters for pagination and searching
		page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
		limit, _ := strconv.Atoi(ctx.URLParamDefault("limit", "10"))
		offset := (page - 1) * limit
		searchTerm := strings.TrimSpace(ctx.URLParamDefault("search", ""))

		// Build the SQL query string
		query := "SELECT unix, symbol, open, high, low, close FROM OHLCData"
		if searchTerm != "" {
			query += fmt.Sprintf(" WHERE symbol ILIKE '%%%s%%'", searchTerm)
		}
		query += fmt.Sprintf(" ORDER BY unix DESC LIMIT %d OFFSET %d", limit, offset)

		// Execute the query
		rows, err := db.Query(query)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}
		defer rows.Close()

		// Iterate through each row and build the JSON response
		var data []map[string]interface{}
		for rows.Next() {
			row := make(map[string]interface{})
			var unix string
			var open, high, low, close string
			var symbol string
			err := rows.Scan(&unix, &symbol, &open, &high, &low, &close)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return
			}
			row["unix"] = unix
			row["symbol"] = symbol
			row["open"] = open
			row["high"] = high
			row["low"] = low
			row["close"] = close
			data = append(data, row)
		}

		// Check for any errors that occurred during iteration
		err = rows.Err()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		// Return the JSON response
		ctx.JSON(data)
	})
}
