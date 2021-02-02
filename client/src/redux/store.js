import { createStore, applyMiddleware, compose } from 'redux';

// import gameReducer from './reducer';
import reducer from './reducer/reducer';
import { networkMiddleware } from './network';

// todo: pass game for action side-chain if needed
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
function buildStore(game) {
    return createStore(reducer, composeEnhancers(applyMiddleware(networkMiddleware(game))));
}

export default buildStore;