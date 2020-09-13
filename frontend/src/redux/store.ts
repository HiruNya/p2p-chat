import {applyMiddleware, createStore} from "redux"
import thunk from "redux-thunk"

const store = createStore<State, Action, any, any>(reducer, applyMiddleware(thunk))

function reducer(
	state: State = {type: "UNCONNECTED", messages: [], ws: null, room: null, nickname: "", peer: (process.env.REACT_APP_DEFAULT_SERVER)? process.env.REACT_APP_DEFAULT_SERVER: "" },
	action: Action): State {

	switch (action.type) {
		case "CONNECTION_START":
			return { ...state, type: "CONNECTING", peer: action.peer, ws: null, messages: [] }
		case "CONNECTION_ESTABLISHED":
			return { ...state, type: "CONNECTED", ws: action.ws }
		case "CONNECTION_ENDED":
			return { ...state, type: "UNCONNECTED" }
		case "MESSAGE":
			const messages = [...state.messages, action.msg]
			return { ...state, messages }
		case "ROOM_JOINED":
			return { ...state, messages: [], room: action.room }
		case "SET_NICKNAME":
			return { ...state, nickname: action.name }
	}
	return state
}

export type State = {
	type: "UNCONNECTED" | "CONNECTING" | "CONNECTED",
	messages: MessageData[],
	ws: WebSocket | null,
	room: string | null,
	peer: string | null,
	nickname: string,
}

export type MessageData = {
	Type: "MESSAGE",
	Text: string,
	User: string,
	Date: string,
}

export type Action = ConnectionStartAction | ConnectionEstablishedAction | ConnectionEndedAction | MessageAction
	| RoomJoinAction | SetNicknameAction
type ConnectionStartAction = { type: "CONNECTION_START", peer: string }
type ConnectionEstablishedAction = { type: "CONNECTION_ESTABLISHED", ws: WebSocket }
type ConnectionEndedAction = { type: "CONNECTION_ENDED" }
type MessageAction = { type: "MESSAGE", msg: MessageData }
type RoomJoinAction = { type: "ROOM_JOINED", room: string }
type SetNicknameAction = { type: "SET_NICKNAME", name: string }

type RoomJoinData = {
	Type: "JOIN",
	Room: string,
}

function connect(peer: string | null, wsOld: WebSocket | null) {
	if (wsOld !== null) {
		wsOld.close()
	}
	return (dispatch: any) => {
		dispatch({ type: "CONNECTION_START", peer })
		if (peer == null) {
			dispatch({type: "CONNECTION_ENDED"})
			return
		}
		let ws: WebSocket;
		try {
			ws = new WebSocket(peer)
		} catch (e) {
			dispatch({type: "CONNECTION_ENDED"})
			return
		}
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
