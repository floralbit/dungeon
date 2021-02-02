import {combineReducers} from 'redux';
import gameReducer from './game.js';
import logReducer from './log.js';
import uiReducer from './ui.js';

export default combineReducers({
    game: gameReducer,
    ui: uiReducer,
    log: logReducer,
});