package conf

type Config struct {
	Server Server `mapstructure:"server"`
	Zap    Zap    `mapstructure:"zap"`
}

type Server struct {
	HttpPort  string `mapstructure:"http_port"`
	RedisPort string `mapstructure:"redis_port"`
	DataPath  string `mapstructure:"data_path"`
}

type Zap struct {
	Director  string `mapstructure:"director"`
	LogFile   string `mapstructure:"log_file"`
	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`
}
