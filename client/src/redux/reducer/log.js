import { NETWORK_RECV_MESSAGE } from "../actions";

const initialState = {
    accountUUID: null,
    uuidToName: {},
    log: [],
};

export default function logReducer(state = initialState, action) {
    if (action.type === NETWORK_RECV_MESSAGE) {
        const data = action.payload;

        if (data.connect) {
            return {
                ...state,
                accountUUID: data.connect.uuid,
            };
        }

        if (data.zone?.load) {
            const uuidToName = state.uuidToName;
            for (let entityUUID in data.zone.load.entities) {
                uuidToName[entityUUID] = data.zone.load.entities[entityUUID].name;
            }

            return {
                ...state,
                uuidToName,
                log: [
                    ...state.log,
                    {zone: {
                        name: data.zone.load.name,
                    }},
                ]
            };
        }

        if (data.entity?.spawn) {
            if (data.entity.spawn.uuid == state.accountUUID) {
                return state;
            }

            return {
                ...state,
                uuidToName: {
                    ...state.uuidToName,
                    [data.entity.spawn.uuid]: data.entity.spawn.name,

                },
                log: [
                    ...state.log,
                    {spawn: {
                        name:  data.entity.spawn.name,
                    }},
                ]
            };
        }

        if (data.entity?.despawn && !data.entity.die) {
            if (data.entity.uuid === state.accountUUID) {
                return state;
            }

            return {
                ...state,
                log: [
                    ...state.log,
                    {despawn: {
                        name: state.uuidToName[data.entity.uuid],
                    }},
                ],
            };
        }

        if (data.entity?.chat) {
            return {
                ...state,
                log: [
                    ...state.log,
                    {chat: {
                        name: state.uuidToName[data.entity.uuid],
                        message: data.entity.chat.message,
                    }}
                ]
            }
        }

        if (data.entity?.attack) {
            return {
                ...state,
                log: [
                    ...state.log,
                    {attack: {
                        attacker: state.uuidToName[data.entity.uuid],
                        target: state.uuidToName[data.entity.attack.target],
                        hit: data.entity.attack.hit,
                        damage: data.entity.attack.damage,
                    }}
                ]
            }
        }

        if (data.entity?.die) {
            return {
                ...state,
                log: [
                    ...state.log,
                    {die: {
                        name: state.uuidToName[data.entity.uuid],
                    }}
                ]
            }
        }

        if (data.message?.message) {
            return {
                ...state,
                log: [
                    ...state.log,
                    {serverMessage: {
                        message: data.message.message
                    }}
                ]
            }
        }
    }

    return state;
}