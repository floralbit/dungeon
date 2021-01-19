import {NETWORK_RECV_MESSAGE} from './actions';

const initialState = {
  messages: [],
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

    default:
      return state;
  }
};