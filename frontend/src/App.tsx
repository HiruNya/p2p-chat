import React, {useEffect} from 'react'
import {useDispatch, useSelector} from "react-redux"
import {AppBar, Container, Paper, Toolbar, Typography} from '@material-ui/core'
import MessageWindow from "./MessageWindow"
import MessageBar from "./MessageBar"
import {connect, State} from "./redux/store"
import "./App.css"

function App() {
	const state = useSelector((state: State) => state.type)
	const room = useSelector((state: State) => state.room)
	const dispatch = useDispatch()

	useEffect(() => {
		if (state === "UNCONNECTED") {
			dispatch(connect())
		}
	})

	return (
		<Container maxWidth="xl" className="container">
			<Paper>
				<div className="app">
					<AppBar position="static">
						<Toolbar>
							<Typography variant="h6">
								{ (state==="CONNECTED")? `Connected to ${room}!`: "Not Connected" }
							</Typography>
						</Toolbar>
					</AppBar>
					<MessageWindow />
					<MessageBar />
				</div>
			</Paper>
		</Container>
	)
}

export default App
