window.onload = function() {
    const registerForm = document.getElementById("register");
    const loginForm = document.getElementById("login");
    const chatForm = document.getElementById("chat");

    const ws = new WebSocket('ws://' + window.location.host + '/ws');
    ws.addEventListener('message', function(e) {
        const msg = JSON.parse(e.data);
        console.log(msg);
    });

    ws.addEventListener('open', function() {
        console.log('ws connected');
    });

    registerForm.onsubmit = function(e) {
        const username = document.getElementById("register-username").value;
        const password = document.getElementById("register-password").value;

        ws.send(
            JSON.stringify({
                register: {
                    username,
                    password,
                }
            })
        ); 

        return false; // stop action 
    }

    loginForm.onsubmit = function(e) {
        const username = document.getElementById("login-username").value;
        const password = document.getElementById("login-password").value;

        ws.send(
            JSON.stringify({
                login: {
                    username,
                    password,
                }
            })
        ); 

        return false; // stop action 
    }

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