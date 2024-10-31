package cfg

import (
	"fmt"
	"testing"
)

/*
	MysqlDataBase      = "chatZoo" //"happytest" //数据库名字
	MysqlUser          = "root"
	MysqlPwd           = "111111"
	MysqlAddr          = "127.0.0.1:3306"
	MysqlCmdTimeoutSec = 3 * time.Second
	mysqlTableUser     = "User"

	redisAddr          = "127.0.0.1:6379"
	redisPassword      = ""
	redisDB            = 8
	redisCmdTimeoutSec = 3 * time.Second
*/

func TestWriteAnything(t *testing.T) {
	cfg := _ChatZooServerConfig{
		App: make(map[string]_AppConfig),
	}
	cfg.App["gate"] = _AppConfig{
		ListenAddr: "127.0.0.1:7788",
		PprofAddr:  "127.0.0.1:6640",
	}
	cfg.App["login"] = _AppConfig{
		ListenAddr: "127.0.0.1:7789",
		PprofAddr:  "127.0.0.1:6641",
	}
	cfg.Common = _CommonConfig{
		Mysql: MysqlConfig{
			MysqlAddr:          "127.0.0.1:3306",
			MysqlDataBase:      "chatZoo",
			MysqlUser:          "root",
			MysqlPwd:           "111111",
			MysqlCmdTimeoutSec: 3,
		},
		Redis: RedisConfig{
			RedisAddr:          "127.0.0.1:6379",
			RedisPwd:           "",
			RedisDB:            8,
			RedisCmdTimeoutSec: 3,
		},
	}
	WriteAnything("./server.json", cfg)
}

func TestReadAnything(t *testing.T) {
	cfg := _ChatZooServerConfig{
		App: make(map[string]_AppConfig),
	}
	ReadAnything("./server.json", &cfg)
	fmt.Printf("%+v\n", cfg)
}
