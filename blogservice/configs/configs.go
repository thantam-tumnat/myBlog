package configs

type Configs struct {
	PostgreSQL PostgreSQL
	App        Fiber
	Cors       Cors
	Kafka      Kafka
	// Minio      Minio
	Redis      Redis
}

type Fiber struct {
	Host string
	Port string
}

// type Minio struct {
// 	Host      string
// 	Port      string
// 	AccessKey string
// 	SecretKey string
// 	Name      string
// 	Secure    string
// 	TLSDISABLE string
// }

// Database
type PostgreSQL struct {
	Host     string
	Port     string
	Protocol string
	Username string
	Password string
	Database string
	SSLMode  string
	Schema   string
}

type Kafka struct {
	Host     string
	Port     string
	Group    string
	ClientID string
}

type Redis struct {
	Host     string
	Port     string
	// Password string
}

type Cors struct {
	AllowOrigins string
}
