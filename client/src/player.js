import {sendMove} from './redux/actions';
import { TILE_SIZE } from './tilemap';

class Player {
    constructor(data, tilemap) {
        this.uuid = data.uuid;
        this.name = data.name;
        this.tile = data.tile;

        this.x = data.x;
        this.y = data.y;

        this.tilemap = tilemap;

        this.movementTime = 0.25; // in s, TODO: populate from server
        this.movementTimer = 0.0;
    }

    handleMove(x, y) {
        this.x = x;
        this.y = y;
    }

    update(dt, store) {
        const state = store.getState();

        if (this.movementTimer > 0) {
            this.movementTimer -= dt;
        }

        if (state.ui.isTyping) {
            return; // don't take input
        }
    
        const up = state.ui.keyPressed['ArrowUp'];
        const down = state.ui.keyPressed['ArrowDown'];
        const left = state.ui.keyPressed['ArrowLeft'];
        const right = state.ui.keyPressed['ArrowRight'];

        if (this.movementTimer <= 0) {
            let moveX = this.x;
            let moveY = this.y

            if (up) {
                moveY -= 1;
            } else if (down) {
                moveY += 1;
            }
    
            if (left) {
                moveX -= 1;
            } else if (right) {
                moveX += 1;
            }

            if (up || down || left || right) {
                store.dispatch(sendMove(moveX, moveY));
                this.movementTimer = this.movementTime;
            }
        }
    }

    draw(ctx, dt) {
        this.tilemap.drawTile(ctx, this.tile, this.x, this.y);
    }
}

export default Player;