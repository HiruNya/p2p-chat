import React from "react"
import ReactMarkdown from "react-markdown"

type MessageProps = {
	user: string,
	text: string,
	date: string,
}

function Message(props: MessageProps) {
	return (
		<div className="message">
			<div className="message-user">{props.user}:</div>
			<ReactMarkdown
				className="message-contents"
				source={props.text}
				escapeHtml={true}
				allowedTypes={[
					"root", "text", "break", "paragraph", "emphasis", "strong", "thematicBreak", "blockquote", "delete",
					"link", "image", "table", "tableHead", "tableBody", "tableRow", "tableCell", "list", "listItem",
					"inlineCode", "code"
				]}
			/>
			<div className="message-date">{props.date}</div>
		</div>
	)
}

export default Message;
