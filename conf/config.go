package conf

type Makedb struct {
	Server Server
	// Zap    Zap
}

type Server struct {
	HttpPort  string `mapstructure:"http_port"`
	RedisPort string `mapstructure:"redis_port"`
	DataPath  string `mapstructure:"data_path"`
}

type Zap struct {
	Path     string `mapstructure:"path"`
	FileName string `mapstructure:"file_name"`
}
