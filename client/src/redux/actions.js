// network
export const NETWORK_CONNECT = 'NETWORK_CONNECT';
export const NETWORK_CONNECTED = 'NETWORK_CONNECTED';
export const NETWORK_RECV_MESSAGE = 'NETWORK_RECV_MESSAGE';

// client side, just for middlewares
export const SEND_CHAT = 'SEND_CHAT';

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
