import { NETWORK_RECV_MESSAGE, SEND_MOVE } from "../actions";

const initialState = {
    accountUUID: null,
    zone: null,
    messages: [],
};

export default function gameReducer(state = initialState, action) {
    // sends should be guaranteed to happen before recvs, so push player to expected pos
    if (action.type === SEND_MOVE) {
        const {x, y} = action.payload;

        return {
            ...state,
            zone: {
                ...state.zone,
                entities: {
                    ...state.zone.entities,
                    [state.accountUUID]: {
                        ...state.zone.entities[state.accountUUID],
                        x, y,
                    }
                }
            }
        };
    }

    if (action.type === NETWORK_RECV_MESSAGE) {
        const data = action.payload;

        if (data.connect) {
            return {
                ...state,
                accountUUID: data.connect.uuid,
            };
        }

        if (data.zone?.load) {
            return {
                ...state,
                zone: data.zone.load,
            }
        }

        // things that can only happen if a zone is loaded
        if (state.zone) {
            if (data.entity?.spawn) {
                return {
                    ...state,
                    zone: {
                        ...state.zone,
                        entities: {
                            ...state.zone.entities,
                            [data.entity.uuid]: data.entity.spawn,
                        }
                    }
                };
            }

            if (data.entity?.despawn) {
                const {[data.entity.uuid]: value, ...withoutEntity} = state.zone.entities;
                return {
                    ...state,
                    zone: {
                        ...state.zone,
                        entities: {
                        ...withoutEntity
                        }
                    }
                };
            }

            if (data.entity?.move) {
                return {
                    ...state,
                    zone: {
                        ...state.zone,
                        entities: {
                            ...state.zone.entities,
                            [data.entity.uuid]: {
                              ...state.zone.entities[data.entity.uuid],
                              x: data.entity.move.x,
                              y: data.entity.move.y,
                            }
                        }
                    }
                };
            }
        }

        // none of the special handling, just append to log
        return {
            ...state,
            messages: [
                ...state.messages,
                data,
            ],
        };
    }

    return state;
}