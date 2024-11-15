package cfg

const (
	AppTypeLogin = "login"
	AppTypeGate  = "gate"
)

type _ChatZooServerConfig struct {
	App    map[string]_AppConfig
	Common _CommonConfig
}

type _AppConfig struct {
	ListenAddr string `json:"listen_addr"`
	OuterAddr  string `json:"outer_addr"`
	PprofAddr  string `json:"pprof_addr"`
}

type _CommonConfig struct {
	Mysql MysqlConfig `json:"mysql"`
	Redis RedisConfig `json:"redis"`
}

type MysqlConfig struct {
	MysqlAddr          string `json:"mysql_addr"`
	MysqlUser          string `json:"mysql_user"`
	MysqlPwd           string `json:"mysql_pwd"`
	MysqlDataBase      string `json:"mysql_data_base"`
	MysqlCmdTimeoutSec uint8  `json:"mysql_cmd_timeout_sec"`
}

type RedisConfig struct {
	RedisAddr          string `json:"redis_addr"`
	RedisPwd           string `json:"redis_pwd"`
	RedisDB            int    `json:"redis_db"`
	RedisCmdTimeoutSec uint8  `json:"redis_cmd_timeout_sec"`
}
