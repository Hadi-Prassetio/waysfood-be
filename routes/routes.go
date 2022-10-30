package routes

import "github.com/gorilla/mux"

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	ProductRoutes(r)
	AuthRoutes(r)
	OrderRoutes(r)
	CartRoutes(r)
	Transaction(r)
}
