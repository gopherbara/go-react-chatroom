var socket = new WebSocket('ws://localhost:4000/ws');

let connect = (cb) => {
    console.log("connecting")

    socket.onopen = () => {
        console.log("connected")
    }

    socket.onmessage = (msg) =>{
        console.log("Message from socket: ", msg);
        cb(msg);
    }

    socket.onclose = (event) => {
        console.log("socket closed: ", event)
    }

    socket.onerror = (error) => {
        console.log("socket on error: ", error)
    }
};

let sendMessage = (msg) => {
    console.log("sending message", msg);
    socket.send(msg);
};

export {connect, sendMessage};