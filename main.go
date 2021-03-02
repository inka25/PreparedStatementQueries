package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

const (
	// db config
	MysqlUser = "root"
	MysqlPassword = ""
	MysqlHost = "localhost"
	MysqlPort = "3306"
	MysqlDB = "bvgc"

	insertOrder = `INSERT IGNORE INTO orders(
		id, 
		redirect_url, user_id, 
	    transaction_id, country, 
		status, expired_at,
		total_price,
		coin_amount, details
	) VALUES (
	  	?,
		?, ?,
		?, ?,
		?, ?,
		?,
		?, ?
	)`

	insertSettings = `INSERT IGNORE INTO settings(
		type, country, settings
	) VALUES (
	  	?, ?, ?
	)`

	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	executeRow = 100000
)

func main(){
	rand.Seed(time.Now().UnixNano())

	db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		MysqlUser,
		MysqlPassword,
		MysqlHost,
		MysqlPort,
		MysqlDB,
	))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(200)
	//fmt.Println("connection db ok")


	ctx := context.Background()

	insertSettingsStmt, err := db.PreparexContext(ctx, insertSettings)
	if err != nil{
		panic(err)
	}

	_, err = db.ExecContext(ctx,`truncate settings`)
	if err != nil{
		panic(err)
	}

	start := time.Now()
	for i := 0; i < executeRow; i++{

		typeInput := i % 100
		countryInput := RandStringBytes(2)
		contentInput := RandStringBytes(10)
		contentBytes,_:=json.Marshal(contentInput)

		_,err = insertSettingsStmt.ExecContext(ctx,typeInput,countryInput,contentBytes)
		if err != nil{
			panic(err)
		}


	}
	fmt.Println("prepared stmt execute time ",time.Since(start).Milliseconds(), " milliseconds")

	//------------------------------------------

	//how things are done in repos currently

	_, err = db.ExecContext(ctx,`truncate settings`)
	if err != nil{
		panic(err)
	}
	start = time.Now()
	for i := 0; i < executeRow; i++{

		typeInput := i % 100
		countryInput := RandStringBytes(2)
		contentInput := RandStringBytes(10)
		contentBytes,_:=json.Marshal(contentInput)

		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		_, err = tx.ExecContext(ctx, insertSettings, typeInput,countryInput, contentBytes)
		if err != nil{
			panic(err)
		}

		if err := tx.Commit(); err != nil{
			panic(err)
		}
	}

	fmt.Println("begin commit execute time ",time.Since(start).Milliseconds(), " milliseconds")

	_, err = db.ExecContext(ctx,`truncate settings`)
	if err != nil{
		panic(err)
	}

}


func RandStringBytes(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}