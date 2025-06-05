package pkg

// DBConfig holds the database configuration parameters.
type DBConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Name         string
	SSL          string
	PoolMaxConns int
}
