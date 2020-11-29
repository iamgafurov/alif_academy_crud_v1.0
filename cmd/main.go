package main

import (
	"database/sql"
	"github.com/iamgafurov/crud/cmd/app"
	"github.com/iamgafurov/crud/pkg/customers"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net"
	"net/http"
	"os"
)

func main(){
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://postgres:postgres@localhost:5432/db"

	if err := execute(host, port,dsn);err != nil{
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string)(err error){
	db, err :=sql.Open("pgx",dsn)
	if err != nil {
		return err
	}
	defer func(){
		if cerr := db.Close();cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()
	
	mux := http.NewServeMux()
	customersSvc := customers.NewService(db)
	server := app.NewServer(mux, customersSvc)
	server.Init()

	srv := &http.Server {
		Addr: net.JoinHostPort(host,port),
		Handler: server,
	}
	log.Print("Server is starting on http://",srv.Addr)
	return srv.ListenAndServe()
}

