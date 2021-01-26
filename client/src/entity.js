
export function drawEntities(ctx, entities, tilemap, playerUUID, dt) {
  for (let entityUUID in entities) {
    if (entityUUID != playerUUID) {
      const entity = entities[entityUUID];
      tilemap.drawTile(ctx, entity.tile, entity.x, entity.y);
    }
  }
}