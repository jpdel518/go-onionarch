package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type config struct {
	SQLDriver string
	DbName    string
	DbUser    string
	DbPass    string
}

var DB *sql.DB

const (
	tableArticle = "articles"
	tableAuthor  = "authors"
)

func init() {
	LoadConfig()
}

func LoadConfig() {
	c := config{
		SQLDriver: os.Getenv("RDB_DRIVER"),
		DbName:    os.Getenv("RDB_NAME"),
		DbUser:    os.Getenv("RDB_USER"),
		DbPass:    os.Getenv("RDB_PASSWORD"),
	}

	DB, err := sql.Open(c.SQLDriver, c.DbUser+":"+c.DbPass+"@tcp(mysql:3306)/"+c.DbName+"?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	cmdArt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
 id INTEGER PRIMARY KEY AUTO_INCREMENT,
 uuid VARCHAR(255) NOT NULL UNIQUE,
 name VARCHAR(255),
 email VARCHAR(255),
 password VARCHAR(255),
 created_at DATETIME
 )`, tableArticle)

	cmdAu := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
 id INTEGER PRIMARY KEY AUTO_INCREMENT,
 uuid VARCHAR(255) NOT NULL UNIQUE,
 name VARCHAR(255),
 email VARCHAR(255),
 password VARCHAR(255),
 created_at DATETIME
 )`, tableAuthor)

	_, err = DB.Exec(cmdArt)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = DB.Exec(cmdAu)
	if err != nil {
		log.Fatalln(err)
	}
}
