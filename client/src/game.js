import Tilemap, {TILE_SIZE} from './tilemap';
import Zone from './zone';

class Game {
  constructor(canvas, ctx, store) {
    this.canvas = canvas;
    this.ctx = ctx;
    
    this.camera = {
      x: 0,
      y: 0,
      zoom: 4
    };
    this.cameraSpeed = 250;
  }

  addStore(store) {
    this.store = store; // this feels like a crappy hack, TODO: figure out a better way
  }

  load() {
    // asset loads
    this.tilemap = new Tilemap();
    this.zone = null;

    // handler registering
    this.canvas.onclick = this.mouseClickHandler.bind(this);

    return Promise.all([this.tilemap.load()]);
  }

  changeZone(data) {
    this.zone = new Zone(data, this.tilemap);
  }

  update(dt) {
    const state = this.store.getState();

    if (state.isTyping) {
      return; // don't take input
    }
     
    if (state.keyPressed['ArrowLeft']) {
      this.camera.x -= this.cameraSpeed * dt;
    } else if (state.keyPressed['ArrowRight']) {
      this.camera.x += this.cameraSpeed * dt;
    }

    if (state.keyPressed['ArrowUp']) {
      this.camera.y -= this.cameraSpeed * dt;
    } else if (state.keyPressed['ArrowDown']) {
      this.camera.y += this.cameraSpeed * dt;
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

    this.tilemap.drawTile(this.ctx, 21, 2, 1); // dude, temp

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