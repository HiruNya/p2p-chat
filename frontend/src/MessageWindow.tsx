import React from "react"
import Message from "./Message"
import {useSelector} from "react-redux"
import {State} from "./redux/store"

function MessageWindow() {
	const messages = useSelector((state: State) => state.messages)

	return (
		<div className="message-window">
			<div className="message-box">
				{ messages.map((msg) => <Message user={msg.User} text={msg.Text} date={msg.Date} />) }
			</div>
		</div>
	)
}

export default MessageWindow
