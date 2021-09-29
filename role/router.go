package role

import (

	"hrm/middleware"

	"github.com/gorilla/mux"
)

func HandleRoleRoutes(r *mux.Router) {
//Endpoint for creating a new role
	r.HandleFunc("/newrole",
	 middleware.JwtVerify(middleware.IsAuthorize("create_role", AddNewRole))).Methods("POST")

//Endpoint for registering new user
	r.HandleFunc("/roles", 
	middleware.JwtVerify(middleware.IsAuthorize("read_all_roles", GetRoles))).Methods("GET")

//Endpoint for fetching all users
r.HandleFunc("/roles/:role_id",
 middleware.JwtVerify(middleware.IsAuthorize("read_one_role", GetRole))).Methods("GET")

//Endpoint for fetching a single user by id
r.HandleFunc("/roles/:role_id", 
middleware.JwtVerify(middleware.IsAuthorize("delete_role", DeleteRole))).Methods("DELETE")

//Endpoint for editing a single user by id
r.HandleFunc("/roles/:role_id", 
middleware.JwtVerify(middleware.IsAuthorize("modify_role", EditRole))).Methods("PUT")
}
