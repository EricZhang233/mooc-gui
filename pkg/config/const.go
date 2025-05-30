package config

type Config struct {
	Global Global `json:"global"`
	Users  []User `json:"users"`
}

type Global struct {
	Server string `json:"server"`
	Limit  int    `json:"limit"`
}

type User struct {
	BaseURL  string `json:"base_url"`
	SchoolID int    `json:"school_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	// 新增字段
	Name     string `json:"name"`      // 姓名
	Platform string `json:"platform"`  // 平台
	Remark   string `json:"remark"`    // 备注
}

var Conf Config

const VERSION = "v1.3.2-GUI"
