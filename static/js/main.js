console.log("Hello!");

const socket = new WebSocket("ws://localhost:4000/ws");

socket.onmessage = (ev) => {
  console.log("Message from server: ", ev.data);
};

