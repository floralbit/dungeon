import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux';

import UI from './components/ui';
import Game from './game';
import Network from './network';
import { networkMiddleware } from './redux/middleware';
import gameReducer from './redux/reducer';

// setup canvas
const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;
ctx.imageSmoothingEnabled = false;

window.onresize = () => {
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
  ctx.imageSmoothingEnabled = false;
};

// connect to the server before anything...
const network = new Network();
network.connect().then(() => {
  const store = createStore(gameReducer, applyMiddleware(networkMiddleware(network)));
  const game = new Game(canvas, ctx, store, network);
  network.setStore(store);

  // setup UI in react
  const ui = ReactDOM.render(
    <Provider store={store}>
      <UI />
    </Provider>,
    document.getElementById('ui')
  );
  
  // kick off game loop
  game.load().then(() => {
    window.requestAnimationFrame(loop);
  });
  
  let lastRender = 0;
  function loop(timestamp) {
    const dt = (timestamp - lastRender) / 1000;
  
    game.update(dt);
    game.draw(dt);
  
    lastRender = timestamp;
    window.requestAnimationFrame(loop);
  }
});

