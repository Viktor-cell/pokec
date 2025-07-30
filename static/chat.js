const yourName = localStorage.getItem("name");
const msgNode = document.getElementById("msg");
const all_msgsNode = document.getElementById("all-msgs")
const baseUrl = window.location.origin;

const socket = new WebSocket(`${baseUrl.replace("https", "wss")}/message`)

socket.addEventListener("message", e => {
	console.log(e)
	let msgs = JSON.parse(e.data).sort((a, b) => a.scn - b.scn);
	all_msgsNode.innerHTML = ""
	
	for (let msg of msgs ) {
		all_msgsNode.innerHTML += `<li>${msg.fst.User.name}: ${msg.fst.msg}</li>`
	}
})

document.getElementById("name").innerText = yourName

document.getElementById("send-btn").addEventListener("click", () => {
	const msg = msgNode.value.trim();
	if (msg === "") return;

	socket.send(JSON.stringify({
		user: {
			name: yourName
		},
		msg: msg
	})
	)

	msgNode.value = "";
});
