

const app = new Vue({
    el: '#app',
    data: {
        cmdInfo: {},
        createInfo: {},
        taskInfo: {},

        curTask: "",
        loopFlag: false,
    },

    methods: {
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
        addTask() {
            let self = this;
            let info = self.createInfo;
            if (Object.keys(info).length != 3) {
                alert("信息不全");
                return;
            }

            let ajax = getAjax("/addTask");
            ajax.onload = () => {
                if (ajax.status != 200) {
                    alert("添加失败");
                    return;
                }
                self.upTaskInfo()
            }

            ajax.send(JSON.stringify(info));
        },
        stopTask(id, url="/removeTask") {
            let self = this;
            if (!id && self.taskInfo.cur) {
                id = self.taskInfo.cur.id;
                url = "/stopTask";
            }
            if (!id) {
                alert("无任务ID");
                return;
            }

            let ajax = getAjax(url)
            ajax.onload = () => {
                self.upTaskInfo();
            };
            ajax.send(id);
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

function getAjax(url, method = "POST") {
    let ajax = new XMLHttpRequest();
    ajax.open(method, url);
    ajax.responseType = "json";
    return ajax;
}


window.onload = () => {
    app.getCmdInfo(false)
    app.upTaskInfo()
}


function initWebsocket() {
    const socket = new WebSocket('ws://127.0.0.1:6464/ws')
    socket.onopen = function (event) {
        socket.send('init')
        // 完成连接
    }
    socket.onmessage = function (event) {
        // 接受消息
    }
    socket.onclose = function (event) {
        // 关闭
    }
    socket.onerror = function (error) {
        console.log('websocket error:', error);
    }
    return socket
}