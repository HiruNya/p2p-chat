import React, {useState} from "react"
import {Button, Divider, Drawer, List, ListItem, TextField, Typography} from "@material-ui/core"
import {useDispatch, useSelector} from "react-redux"
import {Action, connect, State} from "./redux/store"

type SettingsPaneProps = {
	open: boolean,
	setOpen: (_: boolean) => void,
}

function SettingsPane(props: SettingsPaneProps) {
	const dispatch = useDispatch()
	const ws = useSelector((state: State) => state.ws)
	const peerAddress = useSelector((state: State) => state.peer)
	const [peerValue, setPeerValue] = useState("")
	const nickname = useSelector((state: State) => state.nickname)
	const [nameValue, setNameValue] = useState("")

	function onConnectClicked(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault()
		if (peerValue && peerValue !== "") {
			dispatch(connect(peerValue, ws))
		}
	}

	function onSetNickname(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault()
		const setNickname: Action = {
			type: "SET_NICKNAME",
			name: nameValue,
		}
		dispatch(setNickname)
	}

	return (
		<Drawer
			anchor="right"
			open={props.open}
			onClose={() => props.setOpen(false)}
		>
			<List>
				<ListItem>
					<form className="settings-pane-form" onSubmit={onSetNickname} >
						<Typography variant="h6">Set Nickname</Typography>
						<TextField
							variant="outlined"
							value={nameValue}
							onChange={(event) => setNameValue(event.target.value)}
							placeholder={(nickname)? nickname : ""}
						/>
						<Button type="submit" color="primary" variant="contained">Set</Button>
					</form>
				</ListItem>
				<ListItem>
					<form className="settings-pane-form" onSubmit={onConnectClicked} >
						<Typography variant="h6">Connect to a different peer</Typography>
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
