import React, {useEffect} from 'react'
import {useDispatch, useSelector} from "react-redux"
import {Container, Paper} from '@material-ui/core'
import MessageWindow from "./MessageWindow"
import MessageBar from "./MessageBar"
import {connect, State} from "./redux/store"
import "./App.css"

function App() {
	const state = useSelector((state: State) => state.type)
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
					<MessageWindow />
					<MessageBar />
				</div>
			</Paper>
		</Container>
	)
}

export default App
