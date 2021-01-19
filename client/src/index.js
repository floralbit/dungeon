import React from 'react';
import ReactDOM from 'react-dom';

import UI from './components/ui';
import Game from './game';

const ui = ReactDOM.render(
  <UI />,
  document.getElementById('ui')
);

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

// kick off render loop
const game = new Game(canvas, ctx, ui);

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
