package server

// import (
// 	"database/sql"
// 	"fmt"
// 	"lgc/com"
// 	"lgc/util"
// 	"sync"
// 	"time"

// 	_ "github.com/mattn/go-sqlite3"
// )

// var (
// 	DB      *sql.DB
// 	dbMutex sync.RWMutex
// )

// func InitDB() {
// 	db, err := sql.Open("sqlite3", com.DbPath())
// 	util.ErrCheck(err)
// 	DB = db

// 	_, err = db.Exec(com.SqlStr)
// 	util.ErrCheck(err)
// }

// func insertTask(t *Task) {
// 	dbMutex.Lock()
// 	defer dbMutex.Unlock()

// 	stmt, err := DB.Prepare(`insert into t_task(ip, team, pattern, branch, ct) values(?,?,?,?,?)`)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer stmt.Close()

// 	res, err := stmt.Exec(t.Ip, t.Team, t.Pattern, t.Branch, time.Now())
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		id = 0
// 	}
// 	t.ID = int(id)
// 	// return
// }

// func getTask() (t *Task) {
// 	dbMutex.RLock()
// 	defer dbMutex.RUnlock()

// 	t = &Task{}
// 	err := DB.QueryRow(`select id,state,ip,pattern,team,branch from t_task where state=0 order by ct limit 1`).Scan(&t.ID, &t.State, &t.Ip, &t.Pattern, &t.Team, &t.Branch)
// 	if err != nil {
// 		return nil
// 	}
// 	return
// }

// func todoList() (arr []*Task) {
// 	dbMutex.RLock()
// 	defer dbMutex.RUnlock()

// 	// arr := make([]*Task, 5)
// 	stmt, err := DB.Prepare(`select id,state,ip,pattern,team,branch,ct from t_task where state=? order by ct limit 5`)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query(com.Ready)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	defer rows.Close()

// 	arr = []*Task{}
// 	var ct time.Time
// 	for rows.Next() {
// 		t := &Task{}
// 		err = rows.Scan(&t.ID, &t.State, &t.Ip, &t.Pattern, &t.Team, &t.Branch, &ct)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		} else {
// 			t.CT = ct.Format("2006/01/02 15:04:05")
// 			arr = append(arr, t)
// 		}
// 	}
// 	return
// }

// func doneList() (arr []*Task) {
// 	dbMutex.RLock()
// 	defer dbMutex.RUnlock()

// 	// arr := make([]*Task, 5)
// 	stmt, err := DB.Prepare(`select id,state,ip,pattern,team,branch,et from t_task where state!=? and state!=? order by et desc limit 5`)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query(com.Ready, com.Running)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	defer rows.Close()

// 	var et time.Time
// 	arr = []*Task{}
// 	for rows.Next() {
// 		t := &Task{}
// 		err = rows.Scan(&t.ID, &t.State, &t.Ip, &t.Pattern, &t.Team, &t.Branch, &et)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		} else {
// 			t.ET = et.Format("2006/01/02 15:04:05")
// 			arr = append(arr, t)
// 		}
// 	}
// 	return
// }

// func updateTask(id, state int) (err error) {
// 	dbMutex.Lock()
// 	defer dbMutex.Unlock()

// 	stmt, err := DB.Prepare(`update t_task set state=?, et=? where id=?`)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(state, time.Now(), id)
// 	return
// }
