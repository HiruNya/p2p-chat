import {applyMiddleware, createStore} from "redux"
import thunk from "redux-thunk"

const store = createStore<State, Action, any, any>(reducer, applyMiddleware(thunk))

function reducer(state: State = {type: "UNCONNECTED", messages: [], ws: null, room: null}, action: Action): State {
	switch (action.type) {
		case "CONNECTION_START":
			return { ...state, type: "CONNECTING" }
		case "CONNECTION_ESTABLISHED":
			return { ...state, type: "CONNECTED", ws: action.ws }
		case "CONNECTION_ENDED":
			return { ...state, type: "UNCONNECTED" }
		case "MESSAGE":
			const messages = [...state.messages, action.msg]
			return { ...state, messages }
		case "ROOM_JOINED":
			return { ...state, messages: [], room: action.room }
	}
	return state
}

export type State = {
	type: "UNCONNECTED" | "CONNECTING" | "CONNECTED",
	messages: MessageData[],
	ws: WebSocket | null,
	room: string | null,
}

export type MessageData = {
	Type: "MESSAGE",
	Text: string,
	User: string,
	Date: string,
}

export type Action = ConnectionStartAction | ConnectionEstablishedAction | ConnectionEndedAction | MessageAction
	| RoomJoinAction
type ConnectionStartAction = { type: "CONNECTION_START" }
type ConnectionEstablishedAction = { type: "CONNECTION_ESTABLISHED", ws: WebSocket }
type ConnectionEndedAction = { type: "CONNECTION_ENDED" }
type MessageAction = { type: "MESSAGE", msg: MessageData }
type RoomJoinAction = { type: "ROOM_JOINED", room: string }

type RoomJoinData = {
	Type: "JOIN",
	Room: string,
}


const CONNECTION_ADDRESS = process.env.REACT_APP_SERVER

function connect() {
	return (dispatch: any) => {
		dispatch({ type: "CONNECTION_START" })
		if (CONNECTION_ADDRESS == null) {
			dispatch({type: "CONNECTION_ENDED"})
			return
		}
		const ws = new WebSocket(CONNECTION_ADDRESS)
		if (ws == null) {
			dispatch({type: "CONNECTION_ENDED"})
			return
		}
		ws.addEventListener("open", (_: Event) => dispatch({ type: "CONNECTION_ESTABLISHED", ws }))
		ws.addEventListener("close", (_: CloseEvent) => dispatch({ type: "CONNECTION_ENDED" }))
		ws.addEventListener("message", (msg: MessageEvent) => {
			const event: MessageData | RoomJoinData = JSON.parse(msg.data)
			switch (event.Type) {
				case "MESSAGE":
					dispatch({ type: "MESSAGE", msg: event })
					break
				case "JOIN":
					dispatch({ type: "ROOM_JOINED", room: event.Room })
					break
			}
		})
	}
}

export {connect, store}
