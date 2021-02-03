import React, {Component} from 'react';

function InfoBox(props) {
    const {game, ui} = props;

    if (game.zone?.entities && game.accountUUID in game.zone.entities) {
        const player = game.zone.entities[game.accountUUID];
        let lookingAt = "Nothing";

        if (ui.hovering.isHovering) {
            if (ui.hovering.x >= 0 && ui.hovering.x < game.zone.width && ui.hovering.y >= 0 && ui.hovering.y < game.zone.height) {
                const tile = game.zone.tiles[(ui.hovering.y * game.zone.width) + ui.hovering.x];
                if (tile.name !== "") {
                    lookingAt = tile.name;
                }
                for (let worldObjectUUID in game.zone.world_objects) {
                    const worldObject = game.zone.world_objects[worldObjectUUID];
                    if (worldObject.x == ui.hovering.x && worldObject.y == ui.hovering.y) {
                        lookingAt = worldObject.name;
                    }
                }
                for (let entityUUID in game.zone.entities) {
                    const entity = game.zone.entities[entityUUID];
                    if (entity.x == ui.hovering.x && entity.y == ui.hovering.y) {
                        lookingAt = entity.name;
                    }
                }
            }
        }

        return (
            <div id="info-box-wrapper">
                <div id="info-box">
                    <strong>{player.name}</strong> <br />
                    <strong>HP</strong>: {player.stats.hp}/{player.stats.max_hp} | <strong>AC</strong>: {player.stats.ac} | <strong>LV</strong>: {player.stats.level} | <strong>XP</strong>: {player.stats.xp}/{player.stats.xp_to_next_level} <br />
                    <hr />
                    STR: {player.stats.strength} | DEX: {player.stats.dexterity} | CON: {player.stats.constitution} <br />
                    INT: {player.stats.intelligence} | WIS: {player.stats.wisdom} | CHA: {player.stats.charisma}
                    <hr />
                    <strong>Zone</strong>: {game.zone.name} <br />
                    <strong>Looking at</strong>: {lookingAt}
                </div>
            </div>
        );
    }

    return (
        <></>
    );
}

export default InfoBox;