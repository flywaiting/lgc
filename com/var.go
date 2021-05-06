package com

var CfgStr = `
{
  "sufTasks": ["git add -A", "git commit -m 'jt'", "git push"],
  "root": "./",
  "workDir": "./",
  "db": "data",
  "logDir": "log",
  "pattern": ["product", "client-product", "server-product"]
}
`

var SqlStr = `
drop table if exists "t_task";
create table "t_task" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "state" integer default 0,
    "ip" text not null,
    "team" text not null,
    "pattern" text not null,
    "branch" text not null,
    "ct" datetime,
    "et" datetime
);
`
