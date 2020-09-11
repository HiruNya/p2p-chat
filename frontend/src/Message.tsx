import React from "react"

type MessageProps = {
	user: string,
	text: string,
	date: string,
}

function Message(props: MessageProps) {
	return (
		<div className="message">
			<div className="message-user">{props.user}:</div>
			<div className="message-text">{props.text}</div>
			<div className="message-date">{props.date}</div>
		</div>
	)
}

export default Message;
