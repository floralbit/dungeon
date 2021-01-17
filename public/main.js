window.onload = function() {
    console.log('running...');

    const ws = new WebSocket('ws://' + window.location.host + '/ws');
    ws.addEventListener('message', function(e) {
        const msg = JSON.parse(e.data);
        console.log(msg);
    });

    ws.addEventListener('open', function() {
        console.log('ws connected');
        for (let i = 0; i < 10; i++) {
            ws.send(
                JSON.stringify({
                    chatMessage: {
                        data: `${i}`
                    }
                })
            );
        }
    });
}