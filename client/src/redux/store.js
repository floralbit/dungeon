import { createStore, applyMiddleware } from 'redux';

import gameReducer from './reducer';
import { networkMiddleware } from './network';

// todo: pass game for action side-chain if needed
const store = createStore(gameReducer, applyMiddleware(networkMiddleware()));

export default store;