

const app = new Vue({
    el: '#app',
    data: {},

    methods: {
        addTask() {
            let ajax = new XMLHttpRequest()
            ajax.open("POST", "/addTask");
            ajax.onload = () => {
                if (ajax.status != 200) {
                    alert("添加失败");
                    return;
                }
                console.log(ajax);
            }

            ajax.send(null)
        },
    }
})