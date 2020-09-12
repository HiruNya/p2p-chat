import React, {useState} from "react"
import {Button, Divider, Drawer, List, ListItem, TextField, Typography} from "@material-ui/core"
import {useDispatch, useSelector} from "react-redux"
import {connect, State} from "./redux/store"

type SettingsPaneProps = {
	open: boolean,
	setOpen: (_: boolean) => void,
}

function SettingsPane(props: SettingsPaneProps) {
	const dispatch = useDispatch()
	const peerAddress = useSelector((state: State) => state.peer)
	const [peerValue, setPeerValue] = useState("")

	function onConnectClicked(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault()
		if (peerValue && peerValue !== "") {
			dispatch(connect(peerValue))
		}
	}

	return (
		<Drawer
			anchor="right"
			open={props.open}
			onClose={() => props.setOpen(false)}
		>
			<List>
				<ListItem>
					<form className="settings-pane-form" onSubmit={onConnectClicked} >
						<Typography variant="h6">Connect to a different peer:</Typography>
						<TextField
							variant="outlined"
							value={peerValue}
							onChange={(event) => setPeerValue(event.target.value)}
							placeholder={(peerAddress)? peerAddress : ""}
						/>
						<Button type="submit" color="primary" variant="contained">Connect</Button>
					</form>
				</ListItem>
				<Divider />
			</List>
		</Drawer>
	)
}

export default SettingsPane
