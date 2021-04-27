
drop table if exists "t_task"
create table "t_task" (
    "id" integer primary key autoincrement,
    "state" integer default 0,
    "ip" varchar not null,
    "team" varchar not null,
    "pattern" varchar not null,
    "branch" varchar not null,
    "create" datetime,
    "end" datetime,
);