
class Network {
  connect() {
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', this.handleMessage.bind(this))

    return new Promise(resolve => {
      this.ws.addEventListener('open', resolve);
    });
  }

  handleMessage(event) {
    const data = JSON.parse(event.data);
    console.log(data);
  }
}

export default Network;