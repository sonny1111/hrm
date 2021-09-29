package user

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


func AddUserRole(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//Extract user id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["user_id"])
	if err != nil {
		res := middleware.Response{
			Error: true,
			Message: "Unable to extract req params",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Use role name to extract role id
	user := UserModel{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message:  "Unable to extract req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
		//Use role_name to extract roleId
	stmt := `SELECT role_id FROM roles WHERE role_name =$1`
	row := db.QueryRow(stmt, user.RoleName)
	var roleId uint64
	err = row.Scan(&roleId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Unable to scan role id",
		}
		json.NewEncoder(w).Encode(res)
	}

	//Check if user belongs to a group
	stmt = `SELECT * FROM user_group WHERE user_id = $1`
	row = db.QueryRow(stmt, userId)
	var groupId uint64
	err = row.Scan(&groupId)
	//Check for no data found error
	if err == sql.ErrNoRows{
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error: true,
			Message: "User does not belong to any group",
		}
		json.NewEncoder(w).Encode(res)
	   
		}
		//check for other errors
	if err, ok := err.(*pq.Error); ok{
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	   
		}	
	//Check if user's group has that role
	stmt = `SELECT * FROM group_roles WHERE role_id = $1`	
	row = db.QueryRow(stmt, roleId)
	err = row.Scan(&roleId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		res := middleware.Response{
			Error: true,
			Message:  "User's group does not have the role you are trying to assign the user",
		}
		json.NewEncoder(w).Encode(res)
	}
//Now assign role to user
	stmt = `INSERT INTO user_roles(role_id, user_id) VALUES ($1, $2)`
	_, err = db.Exec(stmt, uint64(roleId), uint64( userId))
	//Check for errors
	if err, ok := err.(*pq.Error); ok {
		//Check for duplicate data/ unique violation
		if err.Code == "42701" || err.Code == "23505"{
			w.WriteHeader(http.StatusFound)
			res := middleware.Response{
				Error: true,
				Message: "User already has the role you are trying to assign",
			}
			json.NewEncoder(w).Encode(res)
		} else{
			//Report other errors
			res := middleware.Response{
				Error: true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If no error, return response 
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "Role assigned to user successfully",
	}
	json.NewEncoder(w).Encode(res)
}

//For revoking a single role from a user
func RevokeUserRole(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//Extract user id and role id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["user_id"])	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message:  "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	roleId, err := strconv.Atoi(params["role_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message:  "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	result, err := db.Exec(stmt, userId, roleId)
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
				Message: "Revoke operation was NOT carried out",
			}
			json.NewEncoder(w).Encode(result)
		}
	}
	//If everything went well, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Error: false,
		Message: "User role revoked successfully",
	}

	json.NewEncoder(w).Encode(res)

}

//For revoking all the roles assigned to a user
func RevoleAllUserRoles(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//Extract user id from req params
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["id"])
	if err != nil {
		res := middleware.Response{
			Error: true,
			Message: "Unable to convert req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `DELETE FROM user_roles WHERE user_id = $1`
	result, err := db.Exec(stmt, userId)
	//Check for any possible error
	if err, ok := err.(*pq.Error); ok {
		res := middleware.Response{
			Error: true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}
	//Check if any row was affected by the delete operation
	if count, err := result.RowsAffected(); err != nil{
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error: true,
			Message: "Error counting rows affected" + err.Error(),
		}
		json.NewEncoder(w).Encode(res)
	}else{
		if count == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := middleware.Response{
				Error: true,
				Message: "Nothing to revoke",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went fine, return response
	w.WriteHeader(http.StatusOK)
	res := middleware.Response{
		Message: "All roles assigned to the user has been revoke. User has no more roles.",
	}
	json.NewEncoder(w).Encode(res)
}



package group

import (
	"encoding/json"
	"hrm/db"
	"hrm/middleware"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

//For adding user to a group
func AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group := GroupModel{}
	//Parse req body to json
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error:   true,
			Message: "Unable to parser req body to json",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Use group name to get group id to insert along with userid into user_group
		//Call db connection
	db := db.ConnectDB()
	defer db.Close()
	stmt := `SELECT group_id FROM groups WHERE group_name = $1`
	row := db.QueryRow(stmt, group.GroupName)
	var groupId uint64
	err := row.Scan(&groupId)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		res := middleware.Response{
			Error: true,
			Message: "Unable to scan result",
		}
		json.NewEncoder(w).Encode(res)
	}
	//Extract user id from req parasm
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["user_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := middleware.Response{
			Error: true,
			Message: "Error converting req params to int",
		}
		json.NewEncoder(w).Encode(res)
	}

  //insert group id and user id into user_group
  stmt = `INSERT INTO user_group(user_id, group_id) VALUES ($1, $2)`
	_, err = db.Exec(stmt, uint64(userId), uint64(groupId))
	//Checking for errors
	if err, ok := err.(*pq.Error); ok {
		//Check for duplicate data
		if err.Code == "42701" || err.Code == "23505"{
			w.WriteHeader(http.StatusConflict)
			res := middleware.Response{
				Error: true,
				Message: "User is already part of this group.",
			}
			json.NewEncoder(w).Encode(res)
		}else{
			//For all other errors
			w.WriteHeader(http.StatusInternalServerError)
			res := middleware.Response{
				Error: true,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(res)
		}
	}
	//If everything went fine, return response
	w.WriteHeader(http.StatusCreated)
	res := middleware.Response{
		Error: false,
		Message: "User assigned to group successfully",
	}
	json.NewEncoder(w).Encode(res)

}

//For removing a user from a group
func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {

}

//For changing user from one group to another
func EditUserGroup(w http.ResponseWriter, r *http.Request) {

}

-- CREATE TABLE user_group(
--     group_id INT,
--     user_id INT UNIQUE
-- );

-- ALTER TABLE user_group ADD CONSTRAINT usr_grp_grpid_fk FOREIGN KEY(group_id) REFERENCES groups(id)
-- ON DELETE CASCADE ON UPDATE SET NULL;
-- ALTER TABLE user_group ADD CONSTRAINT usr_grp_usrid_fk FOREIGN KEY(user_id) REFERENCES users(id)
-- ON DELETE CASCADE ON UPDATE SET NULL;


-- CREATE TABLE user_role(
--     role_id INT NOT NULL,
--     user_id INT UNIQUE NOT NULL
-- );
-- ALTER TABLE user_roles ADD CONSTRAINT usr_rl_rid_fk FOREIGN KEY(role_id) REFERENCES roles(id)
-- ON DELETE CASCADE ON UPDATE SET NULL;
-- ALTER TABLE user_roles ADD CONSTRAINT usr_rl_usrid_fk FOREIGN KEY(user_id) REFERENCES users(id)
-- ON DELETE CASCADE ON UPDATE SET NULL;

SELECT group.group_id FROM groups 
	 JOIN user_groups.group_id
	 ON group.group_id = user_group.user_id
	 JOIN group_roles.
	 ON group.group_id = groups.group_id = group_roles.group_id
	 WHERE group_id = $1
	 LIMIT 1
