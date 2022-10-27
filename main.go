package main

import (
	"fmt"
	"net/http"
	"waysfood/database"
	"waysfood/pkg/mysql"
	"waysfood/routes"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	mysql.DataBaseInit()

	database.RunMigration()

	r := mux.NewRouter()

	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	var port = "5000"
	fmt.Println("server running localhost:" + port)

	http.ListenAndServe("localhost:5000", r)
}
