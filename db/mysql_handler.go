package db

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
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



// GetCheckpoint returns the last scanned commit hash
func GetCheckpoint(repoId string, db *sql.DB) (string, error) {
	var commitHash sql.NullString
	query := sq.Select("commit_hash").
		From("scan_history").
		Where("repo_id = ?", repoId).
		RunWith(db).
		QueryRow()

	err := query.Scan(&commitHash)
	if err != nil {
		return "", err
	}

	return commitHash.String, nil
}

// UpdateCheckpoint insert (or updates if exists) the repo id
// and its latest commit hash in the DB
func UpdateCheckpoint(dir, repoId, latestCommitHash string, db *sql.DB) error {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	_, err := sq.Insert("scan_history").
		Columns("repo_id", "commit_hash").
		Values(repoId, latestCommitHash).
		Suffix("ON DUPLICATE KEY UPDATE commit_hash = ?", latestCommitHash).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
