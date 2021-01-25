import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux'

import UI from './components/ui';
import Game from './game';
import buildStore from './redux/store';
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

const game = new Game(canvas, ctx);

// connect to ws via middleware
const store = buildStore(game);
game.addStore(store);
store.dispatch(networkConnect());

// setup UI in react
// kick off game loop
const ui = ReactDOM.render(
  <Provider store={store}>
    <UI />
  </Provider>,
  document.getElementById('ui')
);


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

