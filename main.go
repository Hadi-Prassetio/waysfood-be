package main

import (
	"fmt"
	"net/http"
	"waysfood/database"
	"waysfood/pkg/mysql"
	"waysfood/routes"

	"github.com/gorilla/mux"
)

func main() {

	mysql.DataBaseInit()

	database.RunMigration()

	r := mux.NewRouter()

	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	var port = "5000"
	fmt.Println("server running localhost:" + port)

	http.ListenAndServe("localhost:5000", r)
}
