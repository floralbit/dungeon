import Tilemap, {TILE_SIZE} from './tilemap';
import Player from './player';
import Zone from './zone';
import {lerp} from './util';

class Game {
  constructor(canvas, ctx, store) {
    this.canvas = canvas;
    this.ctx = ctx;
    
    this.camera = {
      x: 0,
      y: 0,
      zoom: 4
    };
    this.cameraSpeed = 10;

    // game data
    this.player = null;
    this.zone = null;
  }

  addStore(store) {
    this.store = store; // this feels like a crappy hack, TODO: figure out a better way
  }

  load() {
    // asset loads
    this.tilemap = new Tilemap();

    // handler registering
    this.canvas.onclick = this.mouseClickHandler.bind(this);

    return Promise.all([this.tilemap.load()]);
  }

  initPlayer(data) {
    this.player = new Player(data, this.tilemap);
  }

  changeZone(data) {
    this.zone = new Zone(data, this.tilemap);
  }

  handleMove(x, y) {
    this.player.handleMove(x, y);
  }

  update(dt) {
    if (this.player) {
      this.player.update(dt, this.store);
    }

    // center camera on player
    if (this.player) {
      const targetCamX = (this.player.x * TILE_SIZE) - (this.canvas.width / this.camera.zoom)/2;
      const targetCamY = (this.player.y * TILE_SIZE) - (this.canvas.height / this.camera.zoom)/2;
      this.camera.x = lerp(this.camera.x, targetCamX, this.cameraSpeed * dt);
      this.camera.y = lerp(this.camera.y, targetCamY, this.cameraSpeed * dt);
    }
  }

  draw(dt) {
    this.ctx.save();
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    // this.ctx.fillStyle = '#292929'; // grey
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)

    // move camera
    this.ctx.scale(this.camera.zoom, this.camera.zoom);
    this.ctx.translate(-this.camera.x, -this.camera.y);

    // draw world objects
    if (this.zone) {
      this.zone.draw(this.ctx, dt);
    }
    if (this.player) {
      this.player.draw(this.ctx, dt);
    }

    this.ctx.restore();
  }

  mouseClickHandler(event) {
    const { tileX, tileY } = this.canvasToWorldCoordinates(event.x, event.y);
    console.log(tileX, tileY);
  }

  canvasToWorldCoordinates(x, y) {
    const tileX = Math.floor(
      (x / this.camera.zoom + this.camera.x) / TILE_SIZE
    );
    const tileY = Math.floor(
      (y / this.camera.zoom + this.camera.y) / TILE_SIZE
    );
    return { tileX, tileY };
  }
}

export default Game;