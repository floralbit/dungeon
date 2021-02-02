import React, {Component} from 'react';

function InfoBox(props) {
    const {game} = props;

    if (game.zone?.entities && game.accountUUID in game.zone.entities) {
        const player = game.zone.entities[game.accountUUID];
        return (
            <div id="info-box-wrapper">
                <div id="info-box">
                    <strong>{player.name}</strong> <br />
                    <strong>HP</strong>: {player.stats.hp} | <strong>LV</strong>: {player.stats.level} <br />
                    <hr />
                    CHA: {player.stats.charisma} | CON: {player.stats.constitution} | DEX: {player.stats.dexterity} <br />
                    INT: {player.stats.intelligence} | STR: {player.stats.strength} | WIS: {player.stats.wisdom}
                    <hr />
                    <strong>Zone</strong>: {game.zone.name} <br />
                    <strong>Looking at</strong>: Nothing
                </div>
            </div>
        );
    }

    return (
        <></>
    );
}

export default InfoBox;