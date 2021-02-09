import Tilemap, {TILE_SIZE} from './tilemap';
import {lerp} from './util';
import {sendMove, sendAttack, setHovering} from '../redux/actions';

const ENTITY_LERP_SPEED = 18;
const LERP_MIN_DIST = 0.3;
const LERP_CAM_MAX_DIST = 30;
const LERP_MAX_DIST = 5;

const TOGGLE_DEBUG_BOXES = false;

class Game {
    constructor(canvas, ctx, store) {
        this.canvas = canvas;
        this.ctx = ctx;
        this.store = store;
        
        this.camera = {
            x: 0, y:0,
            zoom: 4,
        };
        this.cameraSpeed = 8;

        this.renderedZoneUUID = undefined;
        this.zoneCanvas = document.createElement('canvas');

        this.actionTimer = 0.0;
        this.actionTime = 0.25; // in s, TODO: populate from server

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
                const lerpDistX = Math.abs(this.camera.x - targetCamX);
                const lerpDistY = Math.abs(this.camera.y - targetCamY);

                if (lerpDistX < LERP_MIN_DIST || lerpDistX > LERP_CAM_MAX_DIST) {
                    this.camera.x = targetCamX; 
                }
                if (lerpDistY < LERP_MIN_DIST || lerpDistY > LERP_CAM_MAX_DIST) {
                    this.camera.y = targetCamY;
                }

                // handle actions
                if (this.actionTimer > 0) {
                    this.actionTimer -= dt;
                }
        
                if (ui.isTyping) {
                    return; // don't handle input
                }
        
                const up = ui.keyPressed['ArrowUp'];
                const down = ui.keyPressed['ArrowDown'];
                const left = ui.keyPressed['ArrowLeft'];
                const right = ui.keyPressed['ArrowRight'];

                const lightAttack = ui.keyPressed['KeyX'];
        
                if (this.actionTimer <= 0) {
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
                        if (lightAttack) {
                            this.store.dispatch(sendAttack(moveX, moveY));
                            this.actionTimer = this.actionTime;
                        } else {
                            // collision detection
                            if (moveX >= 0 && moveX < game.zone.width && moveY >= 0 && moveY < game.zone.height) {
                                const t = game.zone.tiles[(moveY * game.zone.width) + moveX];
                                if (!t.solid) {
                                    this.store.dispatch(sendMove(moveX, moveY));
                                    this.actionTimer = this.actionTime;
                                }
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
                const lerpDistX = Math.abs(entity.x - this.entityLerpMap[entityUUID].x);
                const lerpDistY = Math.abs(entity.y - this.entityLerpMap[entityUUID].y);

                if (lerpDistX < LERP_MIN_DIST || lerpDistX > LERP_MAX_DIST) {
                    this.entityLerpMap[entityUUID].x = entity.x
                }
                if (lerpDistY < LERP_MIN_DIST || lerpDistY > LERP_MAX_DIST) {
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
        this.ctx.fillStyle = '#000000';
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

        if (game.queuedAction && TOGGLE_DEBUG_BOXES ) {
            if (game.queuedAction.type === "move") {
                this.ctx.fillStyle = 'rgba(0, 0, 255, .2)';
            } else if (game.queuedAction.type === "attack") {
                this.ctx.fillStyle = 'rgba(255, 0, 0, .2)';
            }
            this.ctx.fillRect(game.queuedAction.x * TILE_SIZE, game.queuedAction.y * TILE_SIZE, TILE_SIZE, TILE_SIZE);
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