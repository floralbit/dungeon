import {NETWORK_RECV_MESSAGE, KEY_DOWN, KEY_UP, SET_TYPING} from './actions';

const initialState = {
  accountUUID: null,
  zone: null,
  messages: [],
  keyPressed: {},
  isTyping: false,
};

export default function gameReducer(state = initialState, action) {
  switch (action.type) {
    case NETWORK_RECV_MESSAGE:
      const data = action.payload;

      if (data.connect) {
        return {
          ...state,
          accountUUID: data.connect.uuid,
        }
      }

      if (data.zone && data.zone.load) {
        return {
          ...state,
          zone: data.zone.load,
        }
      }

      return {
        ...state,
        messages: [
          ...state.messages,
          data,
        ],
      };

    case KEY_DOWN:
      const downCode = action.payload;
      return {
        ...state,
        keyPressed: {
          ...state.keyPressed,
          [downCode]: true,
        },
      };
    
    case KEY_UP:
      const upCode = action.payload;
      return {
        ...state,
        keyPressed: {
          ...state.keyPressed,
          [upCode]: false,
        },
      };
    
    case SET_TYPING:
      const isTyping = action.payload;
      return {
        ...state,
        isTyping,
      };

    default:
      return state;
  }
};