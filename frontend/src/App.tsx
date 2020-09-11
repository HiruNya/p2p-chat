import React from 'react'
import {Container, Paper} from '@material-ui/core'
import MessageWindow from "./MessageWindow"
import SearchBar from "./SearchBar"
import "./App.css"

function App() {
	return (
		<Container maxWidth="xl" className="container">
			<Paper>
				<div className="app">
					<MessageWindow />
					<SearchBar />
				</div>
			</Paper>
		</Container>
	)
}

export default App
