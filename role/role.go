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

//For adding a new role
func AddNewRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	role := RoleModel{}
	//Parse req body to json
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `INSERT INTO roles(role_name, description) VALUES($1, $2)`
	_, err := db.Exec(stmt, role.RoleName, role.Description)
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		//Check for duplicate data
		if err.Code == "42701" || err.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			res := middleware.Response{
				Error:   true,
				Message: "Role already exists",
			}
			json.NewEncoder(w).Encode(res)
		} else {
			//For all other errors
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error:   true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "Role added",
	}
	json.NewEncoder(w).Encode(res)
}

//For deleting a single role
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Get role id from req params
	params := mux.Vars(r)
	roleId, err := strconv.Atoi(params["role_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}

	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `DELETE FROM roles WHERE id = $1`
	result, err := db.Exec(stmt, roleId)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected in the delete operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error:   true,
			Message: "Error counting rows affected in the delete operation",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error:   true,
				Message: "Role probrably does not exist",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything was fine, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error:   false,
		Message: "Role deleted successfully",
	}
	json.NewEncoder(w).Encode(res)

}

//For fetching a role
func GetRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	role := RoleModel{}
	//Get role id from req params
	params := mux.Vars(r)
	roleId, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}

	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT FROM roles WHERE role_id = $1`
	row := db.QueryRow(stmt, roleId)
	err = row.Scan(&role.RoleId, &role.RoleName, &role.Description)
	//Checking for no data found error and other errors
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error:   true,
			Message: "Role not found",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything went fine, return the role object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)

}

//For fetching all the roles in the db
func GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := []RoleModel{}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM roles`
	rows, err := db.Query(stmt)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusNotFound)
		if err.Code == "P0002" || err.Code == "02000" {
			res := middleware.Response{
				Error:   true,
				Message: "Role schema not yet populated",
			}
			json.NewEncoder(w).Encode(res)
		} else {
			//For all other errors
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error:   true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	defer rows.Close()
	for rows.Next() {
		roles := RoleModel{}
		err := rows.Scan(&roles.RoleId, &roles.RoleName, &roles.Description)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			res := middleware.Response{
				Error:   true,
				Message: "Error scanning result set",
			}
			json.NewEncoder(w).Encode(res)
		}
		data = append(data, roles)
	}
	//If everything went fine, return array role objects
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

//For updating a role
func EditRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	role := RoleModel{}
	//Parse req body to json
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parse req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Get role id from req params
	params := mux.Vars(r)
	roleId, err := strconv.Atoi(params["role_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `UPDATE roles SET role_name = $2, description = $3 WHERE id = $1`
	result, err := db.Exec(stmt, uint64(roleId), role.RoleName, role.Description)
	//Check for  errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected in the update operation
	if count, err := result.RowsAffected(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error:   true,
			Message: "Error returning rows affected in the update operation",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotModified)
			res := middleware.Response{
				Error:   true,
				Message: "Role probably doesnt exist",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything was fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "Role name modified successfully",
	}
	json.NewEncoder(w).Encode(res)
}
