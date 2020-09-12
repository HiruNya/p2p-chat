import {applyMiddleware, createStore} from "redux"
import thunk from "redux-thunk"

const store = createStore<State, Action, any, any>(reducer, applyMiddleware(thunk))

function reducer(state: State = { type: "UNCONNECTED", messages: [], ws: null }, action: Action): State {
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
	}
	return state
}

export type State = {
	type: "UNCONNECTED" | "CONNECTING" | "CONNECTED",
	messages: MessageData[],
	ws: WebSocket | null,
}

export type MessageData = {
	Text: string,
	User: string,
	Date: string,
}

export type Action = ConnectionStartAction | ConnectionEstablishedAction | ConnectionEndedAction | MessageAction
type ConnectionStartAction = { type: "CONNECTION_START" }
type ConnectionEstablishedAction = { type: "CONNECTION_ESTABLISHED", ws: WebSocket }
type ConnectionEndedAction = { type: "CONNECTION_ENDED" }
type MessageAction = { type: "MESSAGE", msg: MessageData }

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
		ws.addEventListener("message", (msg: MessageEvent) => dispatch({ type: "MESSAGE", msg: JSON.parse(msg.data) }))
	}
}

export {connect, store}
