package user

import (
	"hrm/middleware"
	"github.com/gorilla/mux"
)

func HandleUserRoutes(r *mux.Router) {
	//Endpoint for authenticating user
	r.HandleFunc("/authenicate", AuthenticateUser).Methods("POST")

	//Endpoint for registering new user
	r.HandleFunc("/register",
	 middleware.JwtVerify(middleware.IsAuthorize("create_user", RegisterUser))).Methods("POST")
	
	 //Endpoint for fetching all users
	r.HandleFunc("/users", 
	middleware.JwtVerify(middleware.IsAuthorize("read_all_users", GetUsers))).Methods("GET")

	//Endpoint for fetching a single user by id
	r.HandleFunc("/users/:user_id", 
	middleware.JwtVerify(middleware.IsAuthorize("read_one_user", GetUser))).Methods("GET")

	//Endpoint for editing a single user by id
	r.HandleFunc("/users/:user_id",
	 middleware.JwtVerify(middleware.IsAuthorize("modify_user", EditUser))).Methods("PUT")

	 //Endpoint for deleting a user by id
	r.HandleFunc("/users/:user_id",
	middleware.JwtVerify(middleware.IsAuthorize("delete_user", EditUser))).Methods("DELETE")

	//Endpoint for granting role to a user
	r.HandleFunc("/users/:user_id", 
	middleware.JwtVerify(middleware.IsAuthorize("grant_user_role", AssignRoleToUser))).Methods("PUT")
	
	//For revoking roles granted to a user
	r.HandleFunc("/users/:user_id",
middleware.JwtVerify(middleware.IsAuthorize("revoke_user_role", RemoveRoleFromUser))).Methods("PUT")
}
