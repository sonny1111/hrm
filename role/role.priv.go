package role

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

//For assigning privilege to a role
func AddPrivRole(w http.ResponseWriter, r *http.Request){
	privilegeName := RoleModel{}
	// roleName := RoleModel{}
	//Get role id from req params
	params := mux.Vars(r)
	roleId, err := strconv.Atoi(params["id"])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Use privilege name to get privilege id
	if err = json.NewDecoder(r.Body).Decode(&privilegeName); err != nil{
		w.WriteHeader(http.StatusBadRequest)
		res:= middleware.Response{
			Error: true,
			Message: "Unable to parse json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	var privId uint64
	stmt := `SELECT privilege_id FROM privileges WHERE privilege_name = $1`
	row := db.QueryRow(stmt, privilegeName)
	err = row.Scan(&privId)
	if err == sql.ErrNoRows {
		res := middleware.Response{
			Error: true,
			Message: "Privilege not found",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Add role id and privilege id to role_priv
	stmts := `INSERT INTO role_privileges(privilege_id, role_id) VALUES ($1, $2)`
	_, err = db.Exec(stmts, uint64(privId), roleId)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything is fine then return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "Privilege granted successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//For revoking privilege assigned to a role.
func RevokePrivRole(w http.ResponseWriter, r *http.Request){
	//Extract role_id and privilege_id from req params
	params := mux.Vars(r)
	roleId, err := strconv.Atoi(params["role_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	privilegeId, err := strconv.Atoi(params["privilege_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `DELETE FROM role_privileges WHERE privilege_id = $1 AND role_id = $2`
	result, err := db.Exec(stmt, uint64(privilegeId), uint64(roleId))
	if err, ok := err.(*pq.Error); ok {
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if the delete operation was successful
	if count, err := result.RowsAffected(); err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message:  "Error returning rows affected by revoke operation",
		}
		json.NewEncoder(w).Encode(res)
	}else{
		if count == 0{
			result := middleware.Response{
				Error: true,
				Message: "Error! Failed to revoke privilege from role.",
			}
			json.NewEncoder(w).Encode(result)
		}
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error: false,
		Message: "Role privilege revoked successfully",
	}

	json.NewEncoder(w).Encode(res)
}
