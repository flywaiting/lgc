package server

import (
	"database/sql"
	"fmt"
	"lgc/com"
	"lgc/util"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB      *sql.DB
	dbMutex sync.RWMutex
)

func init() {
	db, err := sql.Open("sqlite3", com.DbPath())
	util.ErrCheck(err)
	DB = db

	_, err = db.Exec(com.SqlStr)
	util.ErrCheck(err)
}

func insertTask(t *Task) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt, err := DB.Prepare(`insert into t_task(ip, team, pattern, branch, create) values(?,?,?,?,?)`)
	if err != nil {
		return
	}
	stmt.Close()

	res, err := stmt.Exec(t.Ip, t.Team, t.Pattern, t.Branch, time.Now())
	if err != nil {
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		id = 0
	}
	t.ID = int(id)
	return
}

func getTask() (t *Task) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	err := DB.QueryRow(`select * from t_task where state=0 order by create limit 1`).Scan(t)
	if err != nil {
		fmt.Println("获取任务失败: ", err)
		return nil
	}
	return
}

func todoList() (arr []*Task) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	// arr := make([]*Task, 5)
	stmt, err := DB.Prepare(`select * from t_task where state=? order by create limit 5`)
	if err != nil {
		return nil
	}
	stmt.Close()
	rows, err := stmt.Query(com.Ready)
	if err != nil {
		return nil
	}
	defer rows.Close()

	arr = []*Task{}
	for rows.Next() {
		t := &Task{}
		err = rows.Scan(t)
		if err != nil {
			fmt.Println(err)
		} else {
			arr = append(arr, t)
		}
	}
	return
}

func doneList() (arr []*Task) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	// arr := make([]*Task, 5)
	stmt, err := DB.Prepare(`select * from t_task where state!=? and state!=? order by end desc limit 5`)
	if err != nil {
		return nil
	}
	stmt.Close()
	rows, err := stmt.Query(com.Ready, com.Running)
	if err != nil {
		return nil
	}
	defer rows.Close()

	arr = []*Task{}
	for rows.Next() {
		t := &Task{}
		err = rows.Scan(t)
		if err != nil {
			fmt.Println(err)
		} else {
			arr = append(arr, t)
		}
	}
	return
}

func updateTask(id, state int) (err error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt, err := DB.Prepare(`update t_task set state=? end=? where id=?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(state, time.Now(), id)
	return
}
