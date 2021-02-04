import Tilemap, {TILE_SIZE} from './tilemap';
import {lerp} from './util';
import {sendMove, setHovering} from '../redux/actions';

const ENTITY_LERP_SPEED = 18;

class Game {
    constructor(canvas, ctx, store) {
        this.canvas = canvas;
        this.ctx = ctx;
        this.store = store;
        
        this.camera = {
            x: 0, y:0,
            zoom: 4,
        };
        this.cameraSpeed = 10;

        this.renderedZoneUUID = undefined;
        this.zoneCanvas = document.createElement('canvas');

        this.movementTimer = 0.0;
        this.movementTime = 0.25; // in s, TODO: populate from server

        this.entityLerpMap = {};
    }

    load() {
        this.tilemap = new Tilemap();
        
        this.canvas.onclick = this.mouseClickHandler.bind(this);
        this.canvas.onmouseover = this.mouseOverHandler.bind(this);
        this.canvas.onmouseout = this.mouseOutHandler.bind(this);
        this.canvas.onmousemove = this.mouseMoveHandler.bind(this);

        return Promise.all([this.tilemap.load()]);
    }

    update(dt) {
        const {game, ui} = this.store.getState();

        // do player actions
        if (game.zone?.entities) {
            if (game.accountUUID in game.zone.entities) {
                const player = game.zone.entities[game.accountUUID];

                // center camera on player
                const targetCamX = Math.floor((player.x * TILE_SIZE) - (this.canvas.width / this.camera.zoom)/2);
                const targetCamY = Math.floor((player.y * TILE_SIZE) - (this.canvas.height / this.camera.zoom)/2);
                this.camera.x = lerp(this.camera.x, targetCamX, this.cameraSpeed * dt);
                this.camera.y = lerp(this.camera.y, targetCamY, this.cameraSpeed * dt);
                // camera correction if we get close (avoid artifacts)
                if (Math.abs(this.camera.x - targetCamX) < .3) {
                    this.camera.x = targetCamX; 
                }
                if (Math.abs(this.camera.y - targetCamY) < .3) {
                    this.camera.y = targetCamY;
                }

                // handle movement
                if (this.movementTimer > 0) {
                    this.movementTimer -= dt;
                }
        
                if (ui.isTyping) {
                    return; // don't handle input
                }
        
                const up = ui.keyPressed['ArrowUp'];
                const down = ui.keyPressed['ArrowDown'];
                const left = ui.keyPressed['ArrowLeft'];
                const right = ui.keyPressed['ArrowRight'];
        
                if (this.movementTimer <= 0) {
                    let moveX = player.x;
                    let moveY = player.y;

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
                        // collision detection
                        if (moveX >= 0 && moveX < game.zone.width && moveY >= 0 && moveY < game.zone.height) {
                            const t = game.zone.tiles[(moveY * game.zone.width) + moveX];
                            if (!t.solid) {
                                this.store.dispatch(sendMove(moveX, moveY));
                                this.movementTimer = this.movementTime;
                            }
                        }
                    }
                }
            }

            for (let entityUUID in game.zone.entities) {
                const entity = game.zone.entities[entityUUID];
                if (!(entityUUID in this.entityLerpMap)) {
                    this.entityLerpMap[entityUUID] = {x: entity.x, y: entity.y};
                }
                this.entityLerpMap[entityUUID].x = lerp(this.entityLerpMap[entityUUID].x, entity.x, ENTITY_LERP_SPEED * dt);
                this.entityLerpMap[entityUUID].y = lerp(this.entityLerpMap[entityUUID].y, entity.y, ENTITY_LERP_SPEED * dt);
                if (Math.abs(entity.x - this.entityLerpMap[entityUUID].x) < .3) {
                    this.entityLerpMap[entityUUID].x = entity.x
                }
                if (Math.abs(entity.y - this.entityLerpMap[entityUUID].y) < .3) {
                    this.entityLerpMap[entityUUID].y = entity.y
                }
            }
        }
    }

    draw(dt) {
        const {game, ui} = this.store.getState();

        if (game.zone?.uuid !== this.renderedZoneUUID) {
            this.renderedZoneUUID = game.zone.uuid;
            this.drawZone(game.zone);
        }

        this.ctx.save();
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

        this.ctx.scale(this.camera.zoom, this.camera.zoom);
        this.ctx.translate(-this.camera.x, -this.camera.y);

        if (game.zone) {
            this.ctx.drawImage(this.zoneCanvas, 0, 0);
            for (let worldObjectUUID in game.zone.world_objects) {
                const worldObject = game.zone.world_objects[worldObjectUUID];
                this.tilemap.drawTile(this.ctx, worldObject.tile, worldObject.x, worldObject.y);
            }

            for (let entityUUID in game.zone.entities) {
                const entity = game.zone.entities[entityUUID];
                const lerpData = this.entityLerpMap[entityUUID];
                this.tilemap.drawTile(this.ctx, entity.tile, lerpData.x, lerpData.y);
            }
        }

        this.ctx.restore();
    }

    drawZone(zone) {
        this.zoneCanvas.width = zone.width * TILE_SIZE;
        this.zoneCanvas.height = zone.height * TILE_SIZE;
        const ctx = this.zoneCanvas.getContext('2d');
        for (let x = 0; x < zone.width; x++) {
            for (let y = 0; y < zone.height; y++) {
                const tile = zone.tiles[(y * zone.width) + x];
                this.tilemap.drawTile(ctx, tile.id, x, y);
            }
        }
    }

    mouseClickHandler(event) {
    }

    mouseOverHandler(event) {
        const { tileX, tileY } = this.canvasToWorldCoordinates(event.x, event.y);
        this.store.dispatch(setHovering(true, tileX, tileY));
    }

    mouseOutHandler(event) {
        const { tileX, tileY } = this.canvasToWorldCoordinates(event.x, event.y);
        this.store.dispatch(setHovering(false, tileX, tileY));
    }

    mouseMoveHandler(event) {
        const { tileX, tileY } = this.canvasToWorldCoordinates(event.x, event.y);
        this.store.dispatch(setHovering(true, tileX, tileY));
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