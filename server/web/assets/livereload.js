
// todo: update url
const ws = new WebSocket("http://localhost:1234/livereload?=server");

ws.onmessage = (message) => {
  console.log(message)
}