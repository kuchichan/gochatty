console.log("Hello!");

const socket = new WebSocket("ws://localhost:4000/ws");
document.addEventListener('DOMContentLoaded', () => {
  let ws = null
  const form = document.getElementById("formJoin");
  const userInput = document.getElementById("userInput");
  const status = document.getElementById("status");

  form.addEventListener("submit", function(ev) {
    ev.preventDefault();
    const userName = userInput.value.trim();
    if (!userName) {
      alert("please enter your name!");
      return;

    }
    console.log(`Name: ${userName}`);

  })

  window.addEventListener('beforeunload', () => {
    if (ws) {
      ws.close();
    }
  })
})


