package db

// Store defines instance for a persistent storage type
type Store interface {
	OpenConnection(host, port, user, password, db string) error
	CloseConnection() error
}
