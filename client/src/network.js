import {receiveMessage} from './redux/actions';

class Network {
  connect() {
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', this.handleMessage.bind(this))

    this.messages = [];

    return new Promise(resolve => {
      this.ws.addEventListener('open', resolve);
    });
  }

  // this is a hack
  setStore(store) {
    this.store = store;
  }

  handleMessage(event) {
    const data = JSON.parse(event.data);

    console.log(data);

    if (data.join || data.leave || data.chat) {
      this.messages.push(data);
      this.store.dispatch(receiveMessage(data));
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