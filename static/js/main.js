console.log("Hello!");

const socket = new WebSocket("ws://localhost:4000/ws");
document.addEventListener('DOMContentLoaded', () => {
    let ws = null
    const form = document.getElementById("formJoin");
    const userInput = document.getElementById("userInput");
    const status = document.getElementById("status");
    const submitBtn = document.getElementById("submitBtn");

    const updateStatus = (message, type) => {
        status.textContent = message;
        status.classList = `status ${type}`
        status.style.display = "block";
    }

    const initializeWebsocket = (userName) => {
        updateStatus('Connecting websocket...', 'connecting');
        submitBtn.disabled = true;

        ws = new WebSocket("http://localhost:4000/ws")

        ws.onopen = (event) => {
            console.log("Websocket connected!");
            updateStatus('websocket connected!', 'connected')

            const joinMsg = {
                type: "join",
                payload: {
                    userName: userName
                }
            }
            ws.send(JSON.stringify(joinMsg));
        }
        ws.onclose = (event) => {
            console.log("Websocket connection closed");
            updateStatus("Connection Closed!", "disonnected")
            submitBtn.disabled = false;
            ws = null;
        }
    }

    form.addEventListener("submit", function(ev) {
        ev.preventDefault();
        const userName = userInput.value.trim();
        if (!userName) {
            alert("please enter your name!");
            return;

        }
        console.log(`Name: ${userName}`);
        initializeWebsocket(userName);

    })

    window.addEventListener('beforeunload', () => {
        if (ws) {
            ws.close();
        }
    })
})


