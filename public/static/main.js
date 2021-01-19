window.onload = function() {
    const chatLog = document.getElementById("chat-log");
    const chatInput = document.getElementById("chat");

    const ws = new WebSocket('ws://' + window.location.host + '/ws');
    ws.addEventListener('message', function(e) {
        const msg = JSON.parse(e.data);

        if (msg.join) {
            addChat(chatLog, `${msg.join.From} joined the game.`);
        }

        if (msg.leave) {
            addChat(chatLog, `${msg.leave.From} left the game.`);
        }

        if (msg.chat) {
            addChat(chatLog, `${msg.chat.From}: ${msg.chat.Message}`);
        }

        console.log(msg);
    });

    ws.addEventListener('open', function() {
        console.log('ws connected');
    });

    chatInput.onkeydown = function(e) {
        if (e.key == 'Enter') {
            ws.send(
                JSON.stringify({
                    chat: {
                        message: chatInput.value,
                    }
                })
            );
            chatInput.value = '';
        }
    }
}

function addChat(log, text) {
    const liNode = document.createElement('li');
    const textNode = document.createTextNode(text);
    liNode.appendChild(textNode);
    log.appendChild(liNode);
    log.scrollTop = log.scrollHeight;
}