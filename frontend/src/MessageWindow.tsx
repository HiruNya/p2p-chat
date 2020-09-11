import React, {useState} from "react"
import Message from "./Message"

function MessageWindow() {
	const [arr, setArr] = useState<number[]>([]);

	setTimeout(() => {
		if (arr.length < 100) {
			setArr((arr) => [...arr, arr.length])
		}
	}, 1000)

	return (
		<div className="message-window">
			<div className="message-box">
				{ arr.map((i: number) => <Message user={"User"+i} text={"The message is "+i} date={i+":"+i} />) }
			</div>
		</div>
	)
}

export default MessageWindow
