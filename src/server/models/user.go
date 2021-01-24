package models

import (
	"log"

	"../dbconn"
)

//User exported
type User struct {
	ID      int
	Email   string
	Hash    string
	Name    string
	Surname string
}

//InsertUser exported
func InsertUser(username, pass string) error {
	db, err := dbconn.NewDB()
	sqlStr := "INSERT INTO users(user_name, user_pass) VALUES(?,?)"
	insertQuery, err := db.Prepare(sqlStr)
	_, err = insertQuery.Exec(username, pass)
	return err
}

//GetUsers exported
func GetUsers() []User {
	db, err := dbconn.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	selDB, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}

	user := User{}
	users := []User{}

	for selDB.Next() {

		var userID int
		var userHash, userName string

		err = selDB.Scan(&userID, &userName, &userHash)
		if err != nil {
			panic(err.Error())
		}

		user.ID = userID
		user.Hash = userHash
		user.Name = userName
		users = append(users, user)
	}
	return users
}

//GetUser exported
func GetUser(username string) (User, error) {
	db, err := dbconn.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	selDB, err := db.Query("SELECT * FROM users WHERE user_name=?", username)
	if err != nil {
		panic(err.Error())
	}

	user := User{}

	for selDB.Next() {

		var userID int
		var userName, userHash string

		err = selDB.Scan(&userID, &userName, &userHash)
		if err != nil {
			panic(err.Error())
		}

		user.ID = userID
		user.Hash = userHash
		user.Name = userName
	}
	return user, nil
}
