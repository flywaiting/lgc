

const app = new Vue({
    el: '#app',
    data: {
        cmdInfo: {},
        createInfo: {},
        taskInfo: {},

        loopFlag: false,
        
        socket: null,
        teams: {},
        list: [],    // 仓库分支列表
        pattern: [],  // 打包模式
        curTask: "",    // 当前任务
        logInfo: "",

        upInfo: {}, // 更新映射关系
        env: {}, // 配置分支
        item: {}, // 任务
    },

    mounted() {
        this.socket = this.getWebsocket();
    },

    methods: {
        getWebsocket() {
            let self = this;
            let location = window.location;
            const socket = new WebSocket(`ws://${location.hostname}:${location.port}/ws`);
            socket.onopen = function (event) {
                socket.send('init');
            }
            socket.onmessage = function (event) {
                // let sync = JSON.parse(event.data);
                // if (sync.type) {
                //     alert(sync.info);
                //     return;
                // }
                event.data && self.handler(JSON.parse(event.data));
            }
            socket.onclose = function (event) {
                self.socket = null;
            }
            socket.onerror = function (error) {
                console.log('websocket error:', error);
            }
            return socket;
        },

        handler(sync) {
            if (!sync) return;
            if (sync.type) {
                sync.type > 1 && alert(sync.info);
                if (sync.type == 1) this.logInfo = sync.info;
                return;
            }

            console.log(sync);
            if (sync.team) {
                let teams = {};
                for (let key in (this.teams || {})) {
                    teams[key] = this.teams[key];
                }
                for (let key in sync.team) {
                    teams[key] = sync.team[key];
                }
                this.teams = teams;
            }
            if (sync.list) this.list = sync.list;
            if (sync.pattern) this.pattern = sync.pattern.reverse()
        },
        // 测试服设置
        upTeam() {
            let info = this.upInfo;
            if (!info.team || !info.branch) {
                alert("信息填写不全")
                return;
            }

            let opt = { [info.team]: info.branch };
            this.socket.send(JSON.stringify({ team: opt }));
        },
        // 重新抓取分支列表
        upBranch() {
            this.socket.send("branch");
        },
        // 添加分支
        upEnv() {
            let info = this.env;
            if (!info.branch) {
                alert("选个分支");
                return;
            }

            this.socket.send(JSON.stringify({ branch: info }));
        },
        getLog(id) {
            if (!id || id < 0) {
                alert("日志ID有问题");
                return;
            }
            this.socket.send(JSON.stringify({ log: id }));
        },
        closeLog() {
            this.logInfo = "";
        },
        addTask() {
            let item = this.item;
            if (!item.pattern || !item.team) {
                alert("信息勾选不全");
                return;
            }
            item.branch = this.teams[item.team];
            if (!item.branch) {
                alert("测试服分支未设置");
                return;
            }
            
            this.socket.send(JSON.stringify({ item }));
        },
        stopTask(id) {
            if (!id || id < 0) return;
            this.socket.send(JSON.stringify({ item: { id } }));
        },

        getCmdInfo(refresh) {
            let self = this;
            let method = refresh ? "POST" : "GET";
            
            let ajax = getAjax("/cmdInfo", method);
            ajax.onload = () => {
                // console.log(ajax.response);
                if (ajax.status != 200) {
                    console.error(ajax.response);
                    return
                }
                self.cmdInfo = ajax.response;
                (self.cmdInfo.branchs || []).reverse();
            }

            ajax.send("1");
        },
        upTaskInfo() {
            let self = this;
            let ajax = getAjax("/taskInfo");
            ajax.onload = () => {
                self.curTask = "";
                if (ajax.status != 200) {
                    console.error(ajax.response);
                    return;
                }
                self.taskInfo = ajax.response;
                let cur = self.taskInfo.cur;
                self.curTask = cur && `${cur.ip}: ${cur.pattern}:: ${cur.branch}=>${cur.team}}`;
                cur && self.loop();
            }
            ajax.send()
        },
        loop() {
            let self = this;
            let id = self.taskInfo.cur.id;
            if (self.loopFlag) return;
            
            self.loopFlag = true;
            console.log("loop", id);
            let ajax = getAjax("/loop");
            ajax.onload = () => {
                self.loopFlag = false;

                let done = ajax.status == 200;
                done ? self.upTaskInfo() : self.loop();
                console.log(done);
            };
            setTimeout(() => {
                ajax.send(id);
            }, 3e3);
        },

        getStateClass(state) {
            return {
                succ: state == 2,
                remove: state == 3,
                inter: state == 4,
                kill: state == 5
            }
        },
        getState(state) {
            return ["完成", "移除", "异常", "中断"][state - 2];
        }
    },

    
});

// let socket;

function getAjax(url, method = "POST") {
    let ajax = new XMLHttpRequest();
    ajax.open(method, url);
    ajax.responseType = "json";
    return ajax;
}


window.onload = () => {
    // app.getCmdInfo(false)
    // app.upTaskInfo()
    initWebsocket();
}


function initWebsocket() {
}