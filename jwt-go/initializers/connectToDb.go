package initializers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)
var DB *sql.DB

func ConnectToDb(){
	var err error
	psqlConn:=os.Getenv("DB_URL")
	DB,err=sql.Open("postgres",psqlConn)
	if err !=nil{
		log.Fatalf(err.Error())
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("The database is connected")
}