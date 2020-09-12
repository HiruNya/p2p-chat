import React, {useState} from 'react'
import {useDispatch, useSelector} from "react-redux"
import {Container, Paper} from '@material-ui/core'
import AppBar from "./AppBar"
import MessageWindow from "./MessageWindow"
import MessageBar from "./MessageBar"
import SettingsPane from "./SettingsPane"
import {connect, State} from "./redux/store"
import "./App.css"

function App() {
	const dispatch = useDispatch()
	const [drawerOpen, setDrawerOpen] = useState(false)
	const peer = useSelector((state: State) => state.peer)
	const [initialLoad, setInitialLoad] = useState(true)

	if (initialLoad) {
		dispatch(connect(peer))
		setInitialLoad(false)
	}

	return (
		<Container maxWidth="xl" className="container">
			<Paper>
				<SettingsPane open={drawerOpen} setOpen={setDrawerOpen} />
				<div className="app">
					<AppBar onSettingsButtonPressed={() => setDrawerOpen(true)} />
					<MessageWindow />
					<MessageBar />
				</div>
			</Paper>
		</Container>
	)
}

export default App
