package models

import (
	"log"

	"../dbconn"
)

//Address exported
type Address struct {
	ID      int
	Il      string
	Ilce    string
	Mahalle string
	Sokak   string
}

//InsertAddress exported
func InsertAddress(il, ilce, mahalle, sokak string) error {
	db, err := dbconn.NewDB()
	sqlStr := "INSERT INTO addresses(il, ilce, mah, sok) VALUES(?,?,?,?)"
	insertQuery, err := db.Prepare(sqlStr)
	_, err = insertQuery.Exec(il, ilce, mahalle, sokak)
	return err
}

//DeleteAddress exported
func DeleteAddress(id string) error {
	db, err := dbconn.NewDB()
	sqlStr := "DELETE FROM `addresses` WHERE id = ?"
	insertQuery, err := db.Prepare(sqlStr)
	_, err = insertQuery.Exec(id)
	return err
}

//UpdateAddress exported
func UpdateAddress(id, il, ilce, mahalle, sokak string) error {
	db, err := dbconn.NewDB()
	sqlStr := "UPDATE `addresses` SET il = ?, ilce = ?, mah = ?, sok = ? WHERE id = ?"
	insertQuery, err := db.Prepare(sqlStr)
	_, err = insertQuery.Exec(il, ilce, mahalle, sokak, id)
	return err
}

//GetAddress exported
func GetAddress(id string) (Address, error) {
	db, err := dbconn.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	selDB, err := db.Query("SELECT * FROM addresses WHERE id=?", id)
	if err != nil {
		panic(err.Error())
	}

	address := Address{}

	for selDB.Next() {
		var il, ilce, mahalle, sokak string
		var id int

		err = selDB.Scan(&id, &il, &ilce, &mahalle, &sokak)
		if err != nil {
			panic(err.Error())
		}

		address.ID = id
		address.Il = il
		address.Ilce = ilce
		address.Mahalle = mahalle
		address.Sokak = sokak
	}
	return address, nil
}

//GetAddresses exported
func GetAddresses() []Address {
	db, err := dbconn.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	selDB, err := db.Query("SELECT * FROM addresses")
	if err != nil {
		panic(err.Error())
	}

	address := Address{}
	addresses := []Address{}

	for selDB.Next() {
		var il, ilce, mahalle, sokak string
		var id int

		err = selDB.Scan(&id, &il, &ilce, &mahalle, &sokak)
		if err != nil {
			panic(err.Error())
		}

		address.ID = id
		address.Il = il
		address.Ilce = ilce
		address.Mahalle = mahalle
		address.Sokak = sokak
		addresses = append(addresses, address)
	}
	return addresses
}
