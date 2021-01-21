import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux'

import UI from './components/ui';
import Game from './game';
import store from './redux/store';
import {networkConnect} from './redux/actions';

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

// connect to ws via middleware
store.dispatch(networkConnect());

// setup UI in react
const ui = ReactDOM.render(
  <Provider store={store}>
    <UI />
  </Provider>,
  document.getElementById('ui')
);

// kick off game loop
const game = new Game(canvas, ctx, store);

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

