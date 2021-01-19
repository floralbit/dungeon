import {SEND_CHAT} from './actions';

export function networkMiddleware(network) {
  return function networkMiddlewareInner(storeAPI) {
    return function wrapDispatch(next) {
      return function handleAction(action) {
        if (action.type === SEND_CHAT) {
          network.sendChat(action.payload);
        }
  
        return next(action);
      }
    }
  }
}