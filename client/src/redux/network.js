import {NETWORK_CONNECT, SEND_CHAT, SEND_MOVE, networkConnected, networkRecvMessage} from './actions';

// we can pass game here if we want side effects
export const networkMiddleware = (game) => {
  let ws = null;

  const handleOpen = store => event => {
    console.log('websocket connected', event.target.url);
    store.dispatch(networkConnected());
  }

  const handleMessage = store => event => {
    const data = JSON.parse(event.data);
    console.log(data);
    store.dispatch(networkRecvMessage(data));

    // side effects to game
    if (data.join) {
      if (data.join.You) {
        game.initPlayer(data.join.Data);
      }
    }

    if (data.zone) {
      game.changeZone(data.zone.Data);
    }

    if (data.move) {
      if (data.move.You) {
        game.handleMove(data.move.X, data.move.Y);
      }
    }
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
      
      case SEND_MOVE:
        ws.send(JSON.stringify({
          move: {x: action.payload.x, y: action.payload.y},
        }));
        break;

      default:
        return next(action);
    }
  }
};