package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

type MysqlHandler struct {
	Connection *sql.DB
}

var instance *MysqlHandler
var once sync.Once

func GetInstance() *MysqlHandler {
	once.Do(func() {
		instance = &MysqlHandler{}
	})

	return instance
}

func (handler *MysqlHandler) OpenConnection(host, port, user, password, db string) (err error) {
	connDb, connErr := sql.Open("mysql", handler.BuildDsn(host, port, user, password, db))
	if connErr != nil {
		log.Fatal(connErr)
	}

	handler.Connection = connDb

	return connErr
}

func (handler *MysqlHandler) CloseConnection() (err error) {
	closeErr := handler.Connection.Close()

	return closeErr
}

func (handler *MysqlHandler) BuildDsn(host, port, user, password, db string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?parseTime=true"
}
