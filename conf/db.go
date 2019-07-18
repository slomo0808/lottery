package conf

const DriveName = "mysql"

type DBConfig struct {
	Name      string
	Pwd       string
	Host      string
	Port      int
	Database  string
	IsRunning bool
}

var DBMasterList = []DBConfig{
	{
		Name:      "root",
		Pwd:       "123456",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  "lottery",
		IsRunning: true,
	},
}

var DbMaster DBConfig = DBMasterList[0]
