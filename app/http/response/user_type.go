package response

type UserResponse struct {
	UserId   string         `json:"userId"`
	Username string         `json:"username"`
	RealName string         `json:"realName"`
	Avatar   string         `json:"avatar"`
	Desc     string         `json:"desc"`
	Password string         `json:"-"`
	Token    string         `json:"token"`
	HomePath string         `json:"homePath"`
	Roles    []RoleResponse `json:"roles"`
}
