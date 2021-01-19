export const GRASS_TILE = 831;
export const TILE_SIZE = 32;

export class Tilemap {
  load() {
    this.tilesetImage = new Image();
    this.tilesetImage.src = "/static/tileset.png";
    return new Promise(resolve => {
      this.tilesetImage.onload = resolve;
    });
  }

  drawTile(ctx, tileID, tileX, tileY) {
    const worldX = tileX * TILE_SIZE;
    const worldY = tileY * TILE_SIZE;

    const setX = (tileID * TILE_SIZE) % this.tilesetImage.width;
    const setY =
      Math.floor((tileID * TILE_SIZE) / this.tilesetImage.width) *
      TILE_SIZE;

    ctx.drawImage(
      this.tilesetImage,
      setX,
      setY,
      TILE_SIZE,
      TILE_SIZE,
      worldX,
      worldY,
      TILE_SIZE,
      TILE_SIZE
    );
  }
}

export default Tilemap;