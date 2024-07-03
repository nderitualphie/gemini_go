import React, { useState } from "react";
import "./App.css";

function App() {
	const [input, setInput] = useState("");
	const [response, setResponse] = useState("");
	const [loading, setLoading] = useState(false);

	const handleSubmit = async () => {
		setLoading(true); // Set loading to true when the request starts
		setResponse(""); // Clear previous response
		try {
			const res = await fetch("http://localhost:1323/chat", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ body: input }),
			});
			const data = await res.json();
			setResponse(data.response);
		} catch (error) {
			console.error("Error:", error);
			setResponse("An error occurred. Please try again.");
		} finally {
			setLoading(false); // Set loading to false when the request completes
		}
	};

	return (
		<div className="App">
			<header className="App-header">
				<p>Welcome to chat support, how may I help you?</p>
				<div className="input">
					<input
						type="text"
						placeholder="Type your request here"
						value={input}
						onChange={(e) => setInput(e.target.value)}
						disabled={loading} // Disable input while loading
					/>
					<button onClick={handleSubmit} disabled={loading}>
						{loading ? "Loading..." : "Submit"}
						{/* Show "Loading..." text while	loading */}
					</button>
				</div>
				{loading && <div className="loading">Loading...</div>}
				{response && (
					<div className="response">
						<p>{response}</p>
					</div>
				)}
			</header>
		</div>
	);
}

export default App;
