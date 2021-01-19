import {MESSAGE_RECEIVED} from './actions';

const initialState = {
  messages: [],
};

export default function gameReducer(state = initialState, action) {
  switch (action.type) {
    case MESSAGE_RECEIVED:
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