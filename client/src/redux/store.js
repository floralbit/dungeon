import { createStore, applyMiddleware } from 'redux';

// import gameReducer from './reducer';
import reducer from './reducer/reducer';
import { networkMiddleware } from './network';

// todo: pass game for action side-chain if needed
function buildStore(game) {
    return createStore(reducer, applyMiddleware(networkMiddleware(game)));
}

export default buildStore;