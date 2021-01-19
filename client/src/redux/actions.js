// from server
export const MESSAGE_RECEIVED = 'MESSAGE_RECEIVED';

// client side, just for middlewares
export const SEND_CHAT = 'SEND_CHAT';

export const receiveMessage = (data) => ({
  type: MESSAGE_RECEIVED,
  payload: data,
});

export const sendChat = (message) => ({
  type: SEND_CHAT,
  payload: message,
});