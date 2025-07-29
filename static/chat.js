const yourName = localStorage.getItem("name");
const msgNode = document.getElementById("msg");
const all_msgsNode = document.getElementById("all-msgs")
let lastMsg = {}
const baseUrl = window.location.origin;

document.getElementById("name").innerText = yourName || "Anonymous";

document.getElementById("send-btn").addEventListener("click", () => {
	const msg = msgNode.value.trim();
	if (msg === "") return;


	fetch(`${baseUrl}/message`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			user: {
				name: yourName
			},
			msg: msg
		})
	})

	msgNode.value = "";
});

function printMessages(msgs) {
	msgs.sort((a, b) => a.scn - b.scn);
	all_msgsNode.innerHTML = ""
	
	for (let msg of msgs) {
		console.log(msg.fst.User.name + " said " + msg.fst.msg)
		displayMessage({from: msg.fst.User.name, text: msg.fst.msg})
	}
}

function displayMessage(message) {
	all_msgsNode.innerHTML += `<li>${message.from}: ${message.text}</li>`
}

setInterval(() => {
	fetch(`${baseUrl}/getMessages`).then(r => r.json()).then(m => {
		if (JSON.stringify(m) != JSON.stringify(lastMsg)) {
			printMessages(m)
			lastMsg = msg
		}
	})
}, 1000)
