window.onload = function() {
    const chatForm = document.getElementById("chat");

    const ws = new WebSocket('ws://' + window.location.host + '/ws');
    ws.addEventListener('message', function(e) {
        const msg = JSON.parse(e.data);
        console.log(msg);
    });

    ws.addEventListener('open', function() {
        console.log('ws connected');
    });

    chatForm.onsubmit = function(e) {
        const message = document.getElementById("chat-message").value;

        ws.send(
            JSON.stringify({
                chat: {
                    message,
                }
            })
        )

        return false;
    }
} 