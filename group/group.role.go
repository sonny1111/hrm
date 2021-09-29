package group

import (
	"database/sql"
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"hrm/role"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

//More than one role can be assigned to a group
func AddRoleToGroup(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//Use role name to get the role_id
	role := role.RoleModel{}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT role_id FROM group_roles WHERE role_name = $1`
	row := db.QueryRow(stmt, role.RoleName)
	err := row.Scan(&role.RoleId)
	//Check if role exists
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error: true,
			Message: "Role not found!!!",
		}
		json.NewEncoder(w).Encode(res)
	}
	//check for other errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Extract group_id from req params and convert to int
	params := mux.Vars(r)
	groupId, err := strconv.Atoi(params["group_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if role has been assigned to the group already
	stmt = `SELECT role_id FROM group_roles WHERE role_id = $1 AND group_id = $2`
	row = db.QueryRow(stmt, role.RoleId, groupId)
	var myRoleId uint64
	err = row.Scan(&myRoleId)
		//Check for other errors
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error: true,
				Message: "Internal server error" + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	//Check if any row returned
	if myRoleId != 0 {
		w.WriteHeader(http.StatusFound)
		res := middleware.Response{
			Error: true,
			Message: "Duplicate data!!! Role already assigned to group",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Now, assigned role to group
	stmt = `INSERT INTO group_roles(group_id, role_id) VALUES ($1, $2)`
	_, err = db.Exec(stmt, uint64(groupId), uint64(role.RoleId)) 
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything went fine, then return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "Role assigned to group successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func RemoveRoleFromGroup(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//Use role name to get role id
	role := role.RoleModel{}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT role_id FROM roles WHERE role_name = $1`
	row := db.QueryRow(stmt, role.RoleName)
	err := row.Scan(&role.RoleId)
	//Check if role exists
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error: true,
			Message: "Role not found!",
		}
		json.NewEncoder(w).Encode(res)
	}
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error: true,
			Message: "Internal server error" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	
}