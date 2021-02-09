// network
export const NETWORK_CONNECT = 'NETWORK_CONNECT';
export const NETWORK_CONNECTED = 'NETWORK_CONNECTED';
export const NETWORK_RECV_MESSAGE = 'NETWORK_RECV_MESSAGE';

export const KEY_DOWN = 'KEY_DOWN';
export const KEY_UP = 'KEY_UP';

export const SET_TYPING = 'SET_TYPING';
export const SET_HOVERING = 'SET_HOVERING';

// client side, just for middlewares
export const SEND_CHAT = 'SEND_CHAT';
export const SEND_MOVE = 'SEND_MOVE';
export const SEND_ATTACK = 'SEND_ATTACK';

export const networkConnect = () => ({
  type: NETWORK_CONNECT,
});

export const networkConnected = () => ({
  type: NETWORK_CONNECTED,
});

export const networkRecvMessage = data => ({
  type: NETWORK_RECV_MESSAGE,
  payload: data,
});

export const sendChat = (message) => ({
  type: SEND_CHAT,
  payload: message,
});

export const sendMove = (x, y) => ({
  type: SEND_MOVE,
  payload: {x, y},
});

export const sendAttack = (x, y) => ({
  type: SEND_ATTACK,
  payload: {x, y},
});

export const keyDown = (code) => ({
  type: KEY_DOWN,
  payload: code,
});

export const keyUp = (code) => ({
  type: KEY_UP,
  payload: code,
});

export const setTyping = (status) => ({
  type: SET_TYPING,
  payload: status,
});

export const setHovering = (isHovering, hoveringX, hoveringY) => ({
  type: SET_HOVERING,
  payload: {
    isHovering,
    hoveringX,
    hoveringY,
  },
});