package models

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)


type PagingType struct {
	Page		int
	Pagesize	int
}

type OrderingType struct {
	Order 		string
}

type LimitType struct {
	Limit int
}

type OptionType struct {
	Page		int
	Pagesize	int
	Order		string
	Limit		int
}

type Where struct {
	Column		string
	Value		interface{}
	Compare		string
}


func Paging(page int , pagesize int) PagingType {
	return PagingType{Page: page, Pagesize: pagesize}
}

func Ordering(order string) OrderingType {
	return OrderingType{Order: order}
}

func Limit(limit int) LimitType {
	return LimitType{Limit: limit}
}

func GetConnection() *sql.DB {
	err := godotenv.Load()
  	if err != nil {
    	log.Fatal("Error loading .env file")
  	}
  
	r1, err := sql.Open(os.Getenv("DATABASE"),os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Database Connect Error")
		return nil
	}

	r1.SetMaxOpenConns(1000)
	r1.SetMaxIdleConns(10)
	r1.SetConnMaxIdleTime(5 * time.Minute)

	return r1
}


func NewConnection() *sql.DB {
	db := GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(100 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(500 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(1 * time.Second)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(2 * time.Second)

	db = GetConnection()

	return db
}