package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DbConnection *sql.DB
var sc = bufio.NewScanner(os.Stdin)

type Message struct {
	id   uint64
	text string
}

func init() {
	var err error
	DbConnection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", "test", "secret", "127.0.0.1:3308", "app-database"))
	if err != nil {
		fmt.Println("open error", err)
	}
	_, err = DbConnection.Query("CREATE TABLE IF NOT EXISTS message(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, text VARCHAR(255) NOT NULL)")
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
