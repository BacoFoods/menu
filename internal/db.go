package internal

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWD"`
	Name     string `env:"DB_NAME"`
	Port     int    `env:"DB_PORT"`
}
