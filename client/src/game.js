import Tilemap, {TILE_SIZE} from './tilemap';
import Terrain from './terrain';

class Game {
  constructor(canvas, ctx) {
    this.canvas = canvas;
    this.ctx = ctx;
    
    this.camera = {
      x: 0,
      y: 0,
      zoom: 2
    };
    this.cameraSpeed = 250;

    this.keyPressedState = {};
  }

  load() {
    // asset loads
    this.tilemap = new Tilemap();

    // todo: load map by network
    this.currentMap = new Terrain(20, 20);

    // handler registering
    window.onkeydown = this.keyDownHandler.bind(this);
    window.onkeyup = this.keyUpHandler.bind(this);
    this.canvas.onclick = this.mouseClickHandler.bind(this);

    return Promise.all([this.tilemap.load()]);
  }

  update(dt) {
    if (this.isKeyPressed("ArrowLeft")) {
      this.camera.x -= this.cameraSpeed * dt;
    } else if (this.isKeyPressed("ArrowRight")) {
      this.camera.x += this.cameraSpeed * dt;
    }

    if (this.isKeyPressed("ArrowUp")) {
      this.camera.y -= this.cameraSpeed * dt;
    } else if (this.isKeyPressed("ArrowDown")) {
      this.camera.y += this.cameraSpeed * dt;
    }
  }

  draw(dt) {
    this.ctx.save();
    // this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    this.ctx.fillStyle = '#292929'; // grey
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)

    // move camera
    this.ctx.scale(this.camera.zoom, this.camera.zoom);
    this.ctx.translate(-this.camera.x, -this.camera.y);

    // draw world objects
    this.currentMap.draw(this.ctx, this.tilemap, dt);
    this.tilemap.drawTile(this.ctx, 89, 2, 1); // dude

    this.ctx.restore();
  }

  keyDownHandler(event) {
    this.keyPressedState[event.code] = true;
  }

  keyUpHandler(event) {
    this.keyPressedState[event.code] = false;
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

  isKeyPressed(keyCode) {
    return this.keyPressedState[keyCode] || false;
  }
}

export default Game;