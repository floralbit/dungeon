import {combineReducers} from 'redux';
import gameReducer from './game.js';
import uiReducer from './ui.js';

export default combineReducers({
    game: gameReducer,
    ui: uiReducer,
});