import {NETWORK_RECV_MESSAGE, KEY_DOWN, KEY_UP} from './actions';

const initialState = {
  messages: [],
  keyPressed: {},
};

export default function gameReducer(state = initialState, action) {
  switch (action.type) {
    case NETWORK_RECV_MESSAGE:
      const data = action.payload;
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

    default:
      return state;
  }
};