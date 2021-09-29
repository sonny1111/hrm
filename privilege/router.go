package privilege

import (
	"hrm/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleUserRoutes(r *mux.Router) {
	//Endpoint for adding a new privilege
	http.HandleFunc("/newpriv", middleware.JwtVerify(middleware.IsAuthorize("add_privilege", AddNewPrivilege)))

	//Endpoint for fetching a single privilege by id
	http.HandleFunc("/privs/:privilege_id",
		middleware.JwtVerify(middleware.IsAuthorize("read_one_priv", GetPrivilege)))

	//Endpoint for fetching all privileges
	http.HandleFunc("/privs",
		middleware.JwtVerify(middleware.IsAuthorize("read_all_privs", GetPrivileges)))

	//Endpoint for editing a single privilege by id
	http.HandleFunc("/privs/:privilege_id",
		middleware.JwtVerify(middleware.IsAuthorize("modify_priv", EditPrivilege)))
}
