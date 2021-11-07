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

-- insert into 