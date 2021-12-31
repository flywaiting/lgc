function f(s) {
    // console.log("load success", s);
    console.log(document.location.host);
}
f("");

let conn
function ws(){
    conn = new WebSocket(`ws://${document.location.host}/ws`);
    conn.onopen = function (evt){
        console.log("websocket open...");
        conn.send("ws for golang");
    }

    conn.onmessage = function (evt){
        console.log("received info: ", evt);
    }

    conn.onclose = function (evt){
        console.log("ws close...", evt);
    }

    conn.onerror = function (evt){
        console.log("error: ", evt);
    }
}