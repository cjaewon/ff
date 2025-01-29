// todo: update url
// todo: timeout
// scrollY must be localized
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
      sessionStorage.setItem('scrollY', window.scrollY);  
    
      document.location.reload(true);


      break;
    default:
      console.error("undefined command, response = ", message)
      break;
  }
}

window.onload = function() {
  const scrollY = sessionStorage.getItem('scrollY');
  
  if (scrollY !== null) {
      window.scrollTo(0, parseInt(scrollY, 10));
      sessionStorage.removeItem('scrollY');
  }
};
