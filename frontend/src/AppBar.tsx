import React from "react"
import {AppBar as MaterialAppBar, IconButton, Toolbar, Typography} from "@material-ui/core"
import {Settings} from "@material-ui/icons"
import {useSelector} from "react-redux"
import {State} from "./redux/store"

type AppBarProps = {
	onSettingsButtonPressed: () => void,
}


function AppBar(props: AppBarProps) {
	const room = useSelector((state: State) => state.room)
	const state = useSelector((state: State) => state.type)
	const name = useSelector((state: State) => state.nickname)

	return (
		<MaterialAppBar position="static">
			<Toolbar className="app-bar">
				<Typography variant="h6">
					{ (state==="CONNECTED")? `Connected to ${room}!`: "Not Connected" }
				</Typography>
				<div className="filler" />
				<Typography variant="h6">
					{(name !== "")? ("@"+name) : ""}
				</Typography>
				<IconButton color="inherit" onClick={props.onSettingsButtonPressed} >
					<Settings />
				</IconButton>
			</Toolbar>
		</MaterialAppBar>
	)
}

export default AppBar;
