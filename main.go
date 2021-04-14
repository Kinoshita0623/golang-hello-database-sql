package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB
var sc = bufio.NewScanner(os.Stdin)

type Message struct {
	id   uint64
	text string
}

func init() {
	var err error
	DbConnection, err = sql.Open("sqlite3", "./database.sql")
	if err != nil {
		fmt.Println("open error", err)
	}
	_, err = DbConnection.Query("CREATE TABLE IF NOT EXISTS message(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL)")
	if err != nil {
		fmt.Println("create table error:", err)
	}

}
func main() {

	fmt.Println("start")

	defer DbConnection.Close()

	sc.Scan()
	message := sc.Text()
	stmt, err := DbConnection.Prepare("INSERT INTO message(text) VALUES(?)")
	if err != nil {
		fmt.Println("prepare error", err)
	}
	_, stmtError := stmt.Exec(message)
	if stmtError != nil {
		fmt.Println("insert execute error:", stmtError)
	}

	result, _ := DbConnection.Query("SELECT * FROM message")
	for result.Next() {
		message := Message{}
		result.Scan(&message.id, &message.text)
		fmt.Println("読み出したメッセージ：", message)
	}

}
