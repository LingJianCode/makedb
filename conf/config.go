package conf

type Makedb struct {
	HttpPort  string `mapstructure:"http_port"`
	RedisPort string `mapstructure:"redis_port"`
	DataPath  string `mapstructure:"data_path"`
}
