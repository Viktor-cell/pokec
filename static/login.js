const nameInput = document.getElementById("name")
const baseUrl = window.location.origin;

document.getElementById("btn").addEventListener("click", e => {
	e.preventDefault()

	const nameVal = nameInput.value.trim()
	nameInput.style = ""

	if (nameVal === "") {
		nameInput.style = "border: 2px solid red"
		return
	}

	fetch(`${baseUrl}/login`, {
		method: "POST",
		headers: {
			"Content-type": "application/json"
		},
		body: JSON.stringify({ name: nameVal })
	})
		.then(r => r.json())
		.then(r => handleLogin(r, nameVal))
})

function handleLogin(serverResponse, nameVal) {
	const isOk = serverResponse.ok
	const message = serverResponse.message
	const redirect = serverResponse.redirect

	let redirectUrl = redirect;

	if (redirect.startsWith('/')) {
		redirectUrl = baseUrl + redirect;
	}

	console.log("Server response: isOk:", isOk, "message:", message)

	if (!isOk) {
		alert(message)
		nameInput.style = "border: 2px solid red"
	} else {
		localStorage.setItem("name", nameVal)
		window.location.href = redirectUrl;
	}
}
