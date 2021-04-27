package server

import (
	"database/sql"
	"lgc/com"
	"lgc/util"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *sql.DB
	dm sync.Mutex
)

//go:embed cfg/sql.sql
var sqlStr string

func init() {
	db, err := sql.Open("sqlite3", com.DbPath())
	util.ErrCheck(err)
	DB = db

	_, err = db.Exec(sqlStr)
	util.ErrCheck(err)
}

func InsertTask(t *Task) int {
	dm.Lock()
	defer dm.Unlock()

	stmt, err := DB.Prepare(`insert into t_task(ip, team, pattern, branch, create) values(?,?,?,?,?)`)
	if err != nil {
		return 0
	}

	res, err := stmt.Exec(t.Ip, t.Team, t.Pattern, t.Branch, time.Now())
	if err != nil {
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		id = 0
	}
	return int(id)
}

func updateTask(t *Task) (err error) {
	dm.Lock()
	defer dm.Unlock()

	stmt, err := DB.Prepare(`update t_task set state=? end=? where id=?`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(t.State, time.Now(), t.ID)
	return
}
