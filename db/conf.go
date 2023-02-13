package db

type DBConfig struct {
	Host     string            `yaml:"Host"`
	Port     int               `yaml:"Port"`
	User     string            `yaml:"User"`
	Password string            `yaml:"Password"`
	Database string            `yaml:"Database"`
	Options  map[string]string `yaml:"Options"`
	ReadDB   *DBServerConfig   `yaml:"Read"`
	WriteDB  *DBServerConfig   `yaml:"Write"`
}

type DBServerConfig struct {
	Host     string            `yaml:"Host"`
	Port     int               `yaml:"Port"`
	User     string            `yaml:"User"`
	Password string            `yaml:"Password"`
	Database string            `yaml:"Database"`
	Options  map[string]string `yaml:"Options"`
}

// TODO 重试
type DBServerOptionsConfig struct {
	Charset      string `yaml:"Charset"`
	ConnTimeout  string `yaml:"ConnTimeout"`
	ReadTimeout  string `yaml:"ReadTimeout"`
	WriteTimeout string `yaml:"WriteTimeout"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
}
