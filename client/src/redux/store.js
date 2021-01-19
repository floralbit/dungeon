import { createStore, applyMiddleware } from 'redux';

import gameReducer from './reducer';
import { networkMiddleware } from './network';

const store = createStore(gameReducer, applyMiddleware(networkMiddleware()));

export default store;