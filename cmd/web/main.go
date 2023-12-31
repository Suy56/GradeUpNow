package main

/*
Function/varibales in different packages(e.g in theory folder) must start with capital letter
if you want to access them in other package(e.g main). And access them with package_name.FunctionName().


Setup:
go get github.com/texttheater/golang-levenshtein/levenshtein	//for pattern matching
go get github.com/gorilla/mux									//to handle http request
go get github.com/go-sql-driver/mysql							//to manage database
go mod init github.com/Suy56/GradeUpNow
go mod tidy

Run program:
go ./cmd/web

Format to enter keywords in the database: 
"word1,word2,word3..."
ENTER ALL WORDS IN LOWER CASE, SEPARATE THEM WITH ',' AND LEAVE NO WHITE SPACE IN BETWEEN.
*/

import (
	"flag"
	"database/sql"
	"log"
	"net/http"
	"os"
	
	_"github.com/go-sql-driver/mysql"
	"github.com/Suy56/GradeUpNow/internal/models"
)


// This struct to make errors more readable and also to access acts as a wrapper around database/sql package. 
type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	user *models.Model
}
//To enter dynamic data into html
type template_data struct{
	User_data []*models.User		
	Individual_user_data *models.User
	Theory_data models.Theory
	Mcq_data models.Mcq
	Model_data models.Model
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
		return db, nil
	}


func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:root@/GradeUpNow?parseTime=true", "MySQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	db,err:=openDB(*dsn)
	if err!=nil{
		errorLog.Fatal(err)
	}
	defer db.Close()
	
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		user:	&models.Model{DB:db},
	}
	

	srv:=&http.Server{
		Addr: 		*addr,
		ErrorLog:	errorLog,
		Handler: 	app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}