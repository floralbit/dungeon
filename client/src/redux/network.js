import {NETWORK_CONNECT, SEND_CHAT, networkConnected, networkRecvMessage} from './actions';

// we can pass game here if we want side effects
export const networkMiddleware = () => {
  let ws = null;

  const handleOpen = store => event => {
    console.log('websocket connected', event.target.url);
    store.dispatch(networkConnected());
  }

  const handleMessage = store => event => {
    const data = JSON.parse(event.data);
    console.log(data);
    store.dispatch(networkRecvMessage(data));
  }

  return store => next => action => {
    switch (action.type) {
      case NETWORK_CONNECT:
        ws = new WebSocket('ws://' + window.location.host + '/ws');
        ws.addEventListener('open', handleOpen(store));
        ws.addEventListener('message', handleMessage(store));
        break;
      
      case SEND_CHAT:
        ws.send(JSON.stringify({
          chat: {message: action.payload},
        }));
        break;
      
      default:
        return next(action);
    }
  }
};