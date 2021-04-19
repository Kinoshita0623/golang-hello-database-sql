package main

import (
	"fmt"
)

type Message struct {
	Id   uint64 `json:"id"`
	Text string `json:"text"`
}

type MessageRepository interface {
	Add(message Message) (Message, error)
	Remove(messageId int64) error
	Find(messageId int64) (Message, error)
	FindAll() ([]Message, error)
}

type MessageRepositoryImpl struct {
}

func CreateMessageRepository() MessageRepository {
	repository := MessageRepositoryImpl{}
	var re MessageRepository
	re = &repository
	return re
}

func (r *MessageRepositoryImpl) FindAll() ([]Message, error) {
	rows, err := DbConnection.Query("SELECT id, text FROM message")
	var messages []Message
	if err != nil {
		return messages, err
	}

	for rows.Next() {
		var msg Message
		rows.Scan(&msg.Id, &msg.Text)
		messages = append(messages, msg)
	}

	return messages, err
}

func (r *MessageRepositoryImpl) Add(message Message) (Message, error) {
	if message.Id == 0 {
		return create(message)
	} else {
		return update(message)
	}
}

func (*MessageRepositoryImpl) Find(id int64) (Message, error) {

	var message Message
	error := DbConnection.QueryRow("SELECT id, text FROM message WHERE id = ? LIMIT 1", id).Scan(&message.Id, &message.Text)
	return message, error
}

func (*MessageRepositoryImpl) Remove(id int64) error {
	_, err := DbConnection.Exec("DELETE FROM message WHERE id = ?", id)
	return err
}

func create(msg Message) (Message, error) {
	var message Message
	stmt, e := DbConnection.Prepare("INSERT INTO message(text) VALUES(?)")
	if e != nil {
		fmt.Println("エラー", e)
	}

	result, err := stmt.Exec(msg.Text)

	if err != nil {
		fmt.Println("stmtエラー", err)
		return message, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("insertId取得失敗")
		return message, err
	}
	message, err = find(insertedId)

	return message, err
}

func find(id int64) (Message, error) {
	var message Message
	err := DbConnection.QueryRow("SELECT id, text FROM message WHERE id = ? LIMIT 1", id).Scan(&message.Id, &message.Text)
	return message, err
}

func update(msg Message) (Message, error) {
	var message Message
	stmt, err := DbConnection.Prepare("UPDATE message SET text=? WHERE id=?")
	if err != nil {
		return message, err
	}
	result, err := stmt.Exec(msg.Id, msg.Text)
	if err != nil {
		return message, err
	}
	id, err := result.LastInsertId()
	return find(id)
}
