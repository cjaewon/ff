// todo: update url
const relativePath = window.location.pathname.replace("/tree/", "");
const ws = new WebSocket(`http://localhost:1234/livereload?relative-path=${encodeURI(relativePath)}`);

/*
message protocal {
  command: string
  data: string
}
*/

ws.onmessage = (raw) => {
  const message = JSON.parse(raw.data);
  const { command, data } = message;

  switch (command) {
    case "refresh":
      location.reload(true);

      break;
    default:
      console.error("undefined command, response = ", message)
      break;
  }
}
