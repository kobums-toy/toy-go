package config

var (
	Database			string
	ConnectionString	string
)

func init() {
	Database = "mysql"
}