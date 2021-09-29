package group

type GroupModel struct{
	UserId uint64 `json:"id"`
	GroupName string `json:"group_name"`
	Description string `json:"description"`
	RoleId uint64 `json:"role_id"`
	GroupId uint64 `json:"group_id"`
}