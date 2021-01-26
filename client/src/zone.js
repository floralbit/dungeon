import { TILE_SIZE } from "./tilemap";

class Zone {
    constructor(data, tilemap) {
        this.width = data.width;
        this.width = data.width;
        this.height = data.height;

        // populate tile data
        this.map = [];
        for (let x = 0; x < this.width; x++) {
            this.map.push([]);
            for (let y = 0; y < this.height; y++) {
                const t = data.tiles[(y * this.width) + x];
                this.map[x][y] = t.id;
            }
        }

        // draw once to hidden canvas for speed
        this.canvas = document.createElement('canvas');
        this.canvas.width = this.width * TILE_SIZE;
        this.canvas.height = this.height * TILE_SIZE;

        const ctx = this.canvas.getContext('2d');
        for (let x = 0; x < this.width; x++) {
            for (let y = 0; y < this.height; y++) {
                const tileID = this.map[x][y];
                tilemap.drawTile(ctx, tileID, x, y);
            }
        }
    }

    draw(ctx, dt) {
        ctx.drawImage(this.canvas, 0, 0);
    }
}

export default Zone;