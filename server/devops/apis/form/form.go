package form

type MachineForm struct {
	Id   uint64 `json:"id"`
	Name string `json:"name" binding:"required"`
	// IP地址
	Ip string `json:"ip" binding:"required"`
	// 用户名
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	// 端口号
	Port int `json:"port" binding:"required"`
}

type MachineRunForm struct {
	MachineId int64  `binding:"required"`
	Cmd       string `binding:"required"`
}

type MachineFileForm struct {
	Id        uint64
	Name      string `binding:"required"`
	MachineId uint64 `binding:"required"`
	Type      int    `binding:"required"`
	Path      string `binding:"required"`
}

type MachineScriptForm struct {
	Id          uint64
	Name        string `binding:"required"`
	MachineId   uint64 `binding:"required"`
	Type        int    `binding:"required"`
	Description string `binding:"required"`
	Params      string
	Script      string `binding:"required"`
}

type DbSqlSaveForm struct {
	Sql  string `binding:"required"`
	Type int    `binding:"required"`
}

type MachineFileUpdateForm struct {
	Content string `binding:"required"`
	Id      uint64 `binding:"required"`
	Path    string `binding:"required"`
}
