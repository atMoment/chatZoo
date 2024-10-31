package cfg

// 只许读不许改
// todo 不够通用, 重新设计, 那不可能每加一个配置量, 就改一遍这个代码吧？
// 先暂时这么地

type IChatZooServerConfig interface {
	GetAppListenAddr(typ string) string
	GetAppPProfAddr(typ string) string
	GetMysqlCfg() MysqlConfig
	GetRedisCfg() RedisConfig
}

func NewServerConfig() IChatZooServerConfig {
	cfg := &_ChatZooServerConfig{}
	ReadAnything("./server.json", cfg)
	return cfg
}

func (c *_ChatZooServerConfig) GetAppListenAddr(typ string) string {
	info, ok := c.App[typ]
	if !ok {
		panic("this typ not config")
	}
	return info.ListenAddr
}

func (c *_ChatZooServerConfig) GetAppPProfAddr(typ string) string {
	info, ok := c.App[typ]
	if !ok {
		panic("this typ not config")
	}
	return info.PprofAddr
}

func (c *_ChatZooServerConfig) GetMysqlCfg() MysqlConfig {
	return c.Common.Mysql
}

func (c *_ChatZooServerConfig) GetRedisCfg() RedisConfig {
	return c.Common.Redis
}
