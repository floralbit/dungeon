import { GRASS_TILE } from "./tilemap";

class Terrain {
  constructor(mapWidth, mapHeight) {
    this.mapWidth = mapWidth;
    this.mapHeight = mapHeight;

    this.map = [];
    for (let x = 0; x < mapWidth; x++) {
      this.map.push([]);
      for (let y = 0; y < mapHeight; y++) {
        this.map[x][y] = GRASS_TILE;
      }
    }
  }

  draw(ctx, tilemap, dt) {
    for (let x = 0; x < this.mapWidth; x++) {
      for (let y = 0; y < this.mapHeight; y++) {
        const tileID = this.map[x][y];
        tilemap.drawTile(ctx, tileID, x, y);
      }
    }
  }
}

export default Terrain;