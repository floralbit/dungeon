import { KEY_DOWN, KEY_UP, SET_TYPING } from "../actions";

const initialState = {
    keyPressed: {},
    isTyping: false,
};

export default function uiReducer(state = initialState, action) {
    switch (action.type) {
        case KEY_DOWN:
            const downCode = action.payload;
            return {
                ...state,
                keyPressed: {
                    ...state.keyPressed,
                    [downCode]: true,
                }
            };
        
        case KEY_UP:
            const upCode = action.payload;
            return {
                ...state,
                keyPressed: {
                    ...state.keyPressed,
                    [upCode]: false,
                }
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
}