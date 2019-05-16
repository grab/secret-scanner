package db

type Store interface {
	OpenConnection(host, port, user, password, db string) error
	CloseConnection() error
}
