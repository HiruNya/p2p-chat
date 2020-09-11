import React from "react"
import {Button, Paper, TextField} from "@material-ui/core"

function SearchBar() {
	return (
		<Paper elevation={5} variant="outlined">
			<div className="search-bar">
				<TextField variant="outlined" placeholder="Enter your message here..."/>
				<Button variant="contained" color="primary">Send</Button>
			</div>
		</Paper>
	)
}

export default SearchBar;
