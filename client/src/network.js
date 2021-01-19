
class Network {
  constructor(ui) {
    this.ui = ui; // this sucks lol, we need to store state somewhere else, let's use redux eventually or something.
  }

  connect() {
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', this.handleMessage.bind(this))

    this.messages = [];

    return new Promise(resolve => {
      this.ws.addEventListener('open', resolve);
    });
  }

  handleMessage(event) {
    const data = JSON.parse(event.data);

    console.log(data);

    if (data.join || data.leave || data.chat) {
      this.messages.push(data);
      this.ui.setState({
        messages: this.messages,
        network: this,
      });
    }
  }

  sendChat(message) {
    this.ws.send(
      JSON.stringify({
          chat: {
              message,
          }
      })
  );
  }
}

export default Network;