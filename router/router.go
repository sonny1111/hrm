package router

import (
	"hrm/user"
	"hrm/role"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	user.HandleUserRoutes(r)
	role.HandleRoleRoutes(r)
	return r
}
