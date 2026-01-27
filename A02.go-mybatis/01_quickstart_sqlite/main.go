package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zhuxiujia/GoMybatis"
)

func main() {
	dbPath := "./test.db"
	os.Remove(dbPath)
	db, _ := sql.Open("sqlite3", dbPath)
	db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT,
		status INTEGER,
		create_time TEXT
	);`)
	db.Close()

	engine := GoMybatis.GoMybatisEngine{}.New()
	_, err := engine.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	session, err := engine.NewSession("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	fmt.Println(">>> Step 1: Insert user via Session.Exec")
	username := "Gopher"
	email := "gopher@golang.org"
	status := 1
	createTime := time.Now().Format("2006-01-02 15:04:05")

	sqlInsert := fmt.Sprintf("INSERT INTO users(username, email, status, create_time) VALUES('%s', '%s', %d, '%s')", 
		username, email, status, createTime)

	affected, err := session.Exec(sqlInsert)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Success: Rows Affected = %d\n", affected)

	fmt.Println("\n>>> Step 2: Select user via Session.Query")
	rows, err := session.Query("SELECT * FROM users LIMIT 1")
	if err != nil {
		panic(err)
	}
	
	if len(rows) > 0 {
		row := rows[0]
		fmt.Println("Success: Retrieved User:")
		for k, v := range row {
			fmt.Printf(" - %s: %s\n", k, string(v))
		}
	}
}
