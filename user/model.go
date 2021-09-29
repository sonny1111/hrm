package user

type UserModel struct{
	UserId uint64 `json:"id"`
	Firstname string `json:"first_name"`
	Lastname string `json:"last_name"`
	Middlename string `json:"middle_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	// Expires string `json:"expires"`
	// Attempts uint8 `json:"attempts"`
	// DaysB4Expn uint64 `json:"days_b4_expn"`
	RoleId uint64 `json:"role_id"`
	RoleName string `json:"role_name"`
}