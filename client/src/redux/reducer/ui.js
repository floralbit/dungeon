import { KEY_DOWN, KEY_UP, SET_HOVERING, SET_TYPING } from "../actions";

const initialState = {
    keyPressed: {},
    isTyping: false,
    hovering: {
        isHovering: false,
        x: 0,
        y: 0,
    }
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
        
        case SET_HOVERING:
            const {isHovering, hoveringX, hoveringY} = action.payload;
            return {
                ...state,
                hovering: {
                    isHovering,
                    x: hoveringX,
                    y: hoveringY,
                }
            }
        
        default:
            return state;
    }
}