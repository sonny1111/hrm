package privilege

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

//For adding a new privilege
func AddNewPrivilege(w http.ResponseWriter, r *http.Request) {
	priv := PrivilegeModel{}
	//Extract priv object from req body
	err := json.NewDecoder(r.Body).Decode(&priv)
	if err != nil {
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
	stmt := `INSERT INTO privileges(privilege_name, description) VALUES ($1, $2)`
	_, err = db.Exec(stmt)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" || err.Code == "42701" {
			w.WriteHeader(http.StatusFound)
			res := middleware.Response{
				Error:   true,
				Message: "Privilege already exists",
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
		Message: "Privilege added successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//For deleting a privilege
func DeletePrivilege(w http.ResponseWriter, r *http.Request) {
	//Extract privilege id from req params
	params := mux.Vars(r)
	privId, err := strconv.Atoi(params["id"])
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
		//Call db connection
		db := db.ConnectDB()
		defer db.Close()
		stmt := `DELETE FROM privileges WHERE privilege_id = $1`
		result, err := db.Exec(stmt, uint64(privId))
		//Check for errors
		if err, ok := err.(*pq.Error); ok {
			res := middleware.Response{
				Error:   true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
		//Check if any row is affected by the delete operation
		if count, err := result.RowsAffected(); err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			res := middleware.Response{
				Error:   true,
				Message: "Error counting rows affected by delete operation",
			}
			json.NewEncoder(w).Encode(res)
		} else {
			if count == 0 {
				w.WriteHeader(http.StatusNotImplemented)
				res := middleware.Response{
					Error:   true,
					Message: "No row affected by the delete operation",
				}
				json.NewEncoder(w).Encode(res)
			}
		}
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error:   false,
		Message: "Privilege removed successfully",
	}
	json.NewEncoder(w).Encode(res)

}

//For fetching a single privilege
func GetPrivilege(w http.ResponseWriter, r *http.Request) {
	priv := PrivilegeModel{}
	//Extract privilege id from req params
	params := mux.Vars(r)
	privId, err := strconv.Atoi(params["id"])
	if err != nil {
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM privileges WHERE privilege_id = $1`
	row := db.QueryRow(stmt, uint64(privId))
	err = row.Scan(&priv.PrivilegeId, &priv.PrivilegeName, &priv.Description)
	//Check for no data found error
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error:   true,
			Message: "Privilege not found",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check for other possible errors
	if err, ok := err.(*pq.Error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		res := middleware.Response{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(priv)
}

//For fetchng all the privileges
func GetPrivileges(w http.ResponseWriter, r *http.Request) {
	data := []PrivilegeModel{}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT * FROM privileges`
	rows, err := db.Query(stmt)
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "02000" || err.Code == "P0002" {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error:   true,
				Message: "Privilege schema not yet populated",
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
		privs := PrivilegeModel{}
		err := rows.Scan(&privs.PrivilegeId, &privs.PrivilegeName, &privs.Description)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			res := middleware.Response{
				Error:   true,
				Message: "Unable to scan privileges result set" + err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
		data = append(data, privs)
	}
	//if everything went well, return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

//For editing privilege
func EditPrivilege(w http.ResponseWriter, r *http.Request) {
	priv := PrivilegeModel{}
	//Extract privilege id from req params
	params := mux.Vars(r)
	privId, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Parse req body to json
	err = json.NewDecoder(r.Body).Decode(&priv)
	if err != nil {
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
	stmt := ` UPDATE privileges SET privilege_name = $2, description = $3 WHERE id = $1`
	result, err := db.Exec(stmt, privId, priv.PrivilegeName, priv.Description)
	//Checking for errors
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
			Message: "Error counting rows affected in the update operation",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		if count == 0 {
			w.WriteHeader(http.StatusNotImplemented)
			res := middleware.Response{
				Error:   true,
				Message: "No row was affected in the update operation",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error:   false,
		Message: "Privilege updated successfully",
	}
	json.NewEncoder(w).Encode(res)
}
