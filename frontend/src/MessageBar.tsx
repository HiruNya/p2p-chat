import React, {FormEvent, useState} from "react"
import {Button, Paper, TextField} from "@material-ui/core"
import {useSelector} from "react-redux"
import {MessageData, State} from "./redux/store"

function MessageBar() {
	const [messageText, setMessageText] = useState("")
	const ws = useSelector((state: State) => state.ws)
	const state = useSelector((state: State) => state.type)
	const name = useSelector((state: State) => state.nickname)

	function onSubmit(event: FormEvent<HTMLFormElement>) {
		event.preventDefault()
		const message: MessageData = {
			Type: "MESSAGE",
			Text: messageText,
			User: name,
			Date: ""
		}
		ws?.send(JSON.stringify(message))
		setMessageText("")
	}

	return (
		<Paper elevation={5} variant="outlined">
			<form className="search-bar" onSubmit={onSubmit}>
				<TextField
					variant="outlined"
					placeholder="Enter your message here..."
					value={messageText}
					disabled={(state !== "CONNECTED")}
					onChange={(event) => setMessageText(event.target.value)}
				/>
				<Button type="submit" variant="contained" color="primary" disabled={(state !== "CONNECTED")}>Send</Button>
			</form>
		</Paper>
	)
}

export default MessageBar;
