package main

import (
	"os"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"database/sql"
)

var selectStatement = `
	SELECT firstname, lastname FROM contacts limit 10
`

var initialStatement = `
	CREATE TABLE IF NOT EXISTS contacts(id SERIAL PRIMARY KEY, firstname VARCHAR(200), lastname VARCHAR(200));
`

func main() {
	http.HandleFunc("/", database)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func database(res http.ResponseWriter, req *http.Request) {

//	fmt.Fprintln(res, "testing...\n")

	var db *sql.DB
	var err error
	dokku_db := os.Getenv("DATABASE_URL")

	db, err = sql.Open("postgres", dokku_db)
	if err != nil {
		fmt.Printf("sql.Open error: %v\n", err)
		return
	}
	defer db.Close()

	err = doInitialize(db)
	if err != nil {
		fmt.Print("initialize error: %v\n", err)
		return
	}

	err = doSelect(res, db)
	if err != nil {
		fmt.Print("select error: %v\n", err)
		return
	}
}

func doInitialize(db *sql.DB) error {
	var stmt *sql.Stmt
	var err error

	stmt, err = db.Prepare(initialStatement)
	if err != nil {
		fmt.Printf("db.Prepare initializing error: %v\n", err)
		return err
	}


	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("stmt.Exec error: %v\n", err)
		return err
	}

	defer stmt.Close()

	stmt, err = db.Prepare("INSERT INTO contacts(firstname,lastname) VALUES($1, $2)")
	if err != nil {
		fmt.Printf("stmt.Prepare error: %v\n", err)
	}

	_, err = stmt.Exec("John", "Smith")
	if err != nil {
		fmt.Printf("stmt.Exec error: %v\n", err)
	}


	return nil
	}


func doSelect(res http.ResponseWriter, db *sql.DB) error {
	var stmt *sql.Stmt
	var err error

	stmt, err = db.Prepare(selectStatement)
	if err != nil {
		fmt.Printf("db.Prepare error: %v\n", err)
		return err
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		fmt.Printf("stmt.Query error: %v\n", err)
		return err
	}


	defer stmt.Close()

	for rows.Next() {
		var firstname string
		var lastname string

		err = rows.Scan(&firstname, &lastname)
		if err != nil {
			fmt.Printf("rows.Scan error: %v\n", err)
			return err
		}

		fmt.Fprintln(res, "firstname: "+firstname+"    lastname: "+lastname+"\n");
		}

		return nil
	}
