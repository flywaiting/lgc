package model

import (
	"errors"
	"lgc/server"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB
var usrMux sync.Mutex
var taskMux sync.Mutex

func init() {
	fn := "data.db"
	// db 文件创建
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		_, err = os.Create(fn)
		if err != nil {
			panic(err)
		}
	}
	// 打开文件
	_db, err := sqlx.Open("sqlite3", fn)
	if err != nil {
		panic(err)
	}
	db = _db
	// 表格创建
	createTable()
}

func createTable() {
	sql := `
create table if not exists lgc_usr (
    uuid integer primary key AUTOINCREMENT,
    ctime integer,
    uname text not null,
    pwd text not null
);
create table if not exists lgc_task (
    tid integer primary key,
    tstatus integer default 0,
    -- 创建时间
    ctime integer,
    -- 开始时间
    stime integer,
    -- 结束时间
    etime integer,
    -- 执行人
    rname text not null,
    -- 任务执行命令
    cmds text not null
);
	`
	db.MustExec(sql)
}

func Register(usr, pwd string) (int64, error) {
	usrMux.Lock()
	defer usrMux.Unlock()

	if Login(usr, pwd) > 0 {
		return 0, errors.New("usr has exist")
	}
	res, err := db.Exec("insert into lgc_usr (uname, pwd, ctime) values (?, ?, ?)", usr, pwd, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Login(usr, pwd string) int64 {
	var uid int64
	db.Get(&uid, "select usr_id from lgc_usr where uname=? and pwd=?", usr, pwd)
	return uid
}

func AddTask(t *server.Task) (int64, error) {
	taskMux.Lock()
	defer taskMux.Unlock()

	res, err := db.NamedExec("insert into lgc_task (tid, tstatus, ctime, stime, etime, rname, cmds) values (:tid, :tstatus, :ctime, :stime, :etime, :rname, :cmds)", t)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Tasks() []*server.Task {
	res := []*server.Task{}
	err := db.Select(&res, "select * from lgc_task order by etime desc limit ?", 2)
	if err != nil {
		// fmt.Println(err)
		return nil
	}
	return res
}
