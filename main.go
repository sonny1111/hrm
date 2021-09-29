package main

import (
	"fmt"
	"hrm/db"
	"hrm/router"
	"log"
	"net/http"
)

func main() {
	//Bringing in all the routes
	r := router.Router()
	//DB connection
	db.ConnectDB() 
	
 	log.Fatal(http.ListenAndServe(":9000", r))
    fmt.Printf("Running")
	
}