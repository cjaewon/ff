/*
message protocal {
  command: string
  data: string
}
*/

const relativePath = window.location.pathname.replace("/tree/", "");
const ws = new WebSocket(`http://${document.location.host}/livereload?relative-path=${encodeURI(relativePath)}`);

ws.onmessage = (raw) => {
  const message = JSON.parse(raw.data);
  const { command, data } = message;

  switch (command) {
    case "refresh":
      sessionStorage.setItem(`${document.location.href}-scrollY`, window.scrollY);  
    
      document.location.reload(true);


      break;
    default:
      console.error("undefined command, response = ", message)
      break;
  }
}

window.onload = function() {
  const scrollY = sessionStorage.getItem(`${document.location.href}-scrollY`);
  
  if (scrollY !== null) {
      window.scrollTo(0, parseInt(scrollY, 10));
      sessionStorage.removeItem(`${document.location.href}-scrollY`);
  }
};
