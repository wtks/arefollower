package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/robfig/cron"
	"log"
	"os"
	"strings"
	"time"
)

const SQL_DATETIME_FORMAT = "2006-01-02 15:04:05"

var (
	HOSTNAME   = os.Getenv("MARIADB_HOSTNAME")
	DATABASE   = os.Getenv("MARIADB_DATABASE")
	USERNAME   = os.Getenv("MARIADB_USERNAME")
	PASSWORD   = os.Getenv("MARIADB_PASSWORD")
	DB_PORT    = 3306
	DB_ADDRESS = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?interpolateParams=true", USERNAME, PASSWORD, HOSTNAME, DB_PORT, DATABASE)
)

var conn *sql.DB

func main() {
	//init database connection
	db, err := sql.Open("mysql", DB_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	conn = db
	conn.Exec(`CREATE TABLE IF NOT EXISTS ranking (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ranked_datetime DATETIME NOT NULL,
  rank INT NOT NULL,
  video_id VARCHAR(20) NOT NULL,
  title VARCHAR(100) NOT NULL,
  upload_date DATETIME NOT NULL,
  thumb_url TEXT NOT NULL,
  length VARCHAR(10) NOT NULL,
  view INT NOT NULL,
  comment INT NOT NULL,
  mylist INT NOT NULL,
  tags TEXT
)`)

	//init & start cron
	c := cron.New()
	c.AddFunc("0 5 * * * *", crawlRanking)
	c.Start()

	//init & start http server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3000"))
}

func crawlRanking() {
	//fetch ranking
	rank, err := FetchRanking()
	for i := 0; err != nil && i < 4; i++ {
		time.Sleep(5 * time.Second)
		rank, err = FetchRanking()
	}
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := conn.Prepare("INSERT INTO ranking (ranked_datetime,rank,video_id,title,upload_date,thumb_url,length,view,comment,mylist,tags) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i, v := range rank.Videos {
		if v == nil {
			continue
		}
		_, err = stmt.Exec(
			rank.PubDate.Format(SQL_DATETIME_FORMAT),
			i+1,
			v.Id,
			v.Title,
			v.UploadDate.Format(SQL_DATETIME_FORMAT),
			v.ThumbUrl,
			v.Length,
			v.View,
			v.Comment,
			v.Mylist,
			strings.Join(v.Tags, ","),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}