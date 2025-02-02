---
title: "ff architecture"
date: 2025-01-30T01:46:03+09:00
---
추후 복기를 위한 ff 내부 구현에 대한 설명 입니다.

## ./main.go, ./commands
cli 구현을 위해 spf13/cobra를 사용합니다.  
commands/root.go의 rootCmd로 모든 커맨드(new, server)를 통합한 후 main.go에서 `commands.Execute()`로 실행합니다. 

## ./server, ./server/server.go
서버 가동을 위한 패키지 입니다.
./server/server.go에 다음과 같은 Server 구조체가 정의되어 있습니다.
이 부분을 commands/server.go가 초기화 후 Run 메서드를 이용하여 실행합니다.

```go
// server/server.go
type Server struct {
	Port        int
	Watch       bool
	Bind        string
	RootDirName string
	RootDirPath string
}

// commands/server.go
s := server.Server{ /* ... */} // 서버 구현체
err := s.Run() // 서버 실행
```

코드 내에 있는 `relativePath`과 `absPath`의 의미는 매우 중요합니다.

`relativePath`은 http://localhost:1234/tree/(여기서부터 relativePath) 같이 /tree/ 뒷 부분에 오는 것을 말합니다.  
`absPath`은 말 그대로 절대 경로를 의미합니다. 코드에서는 모두 다음과 같은 방식으로 정의되어 사용됩니다.
```go
absPath := filepath.Join(s.RootDirPath, relativePath)
// s.RootDirPath 또한 절대 경로 입니다.
```

### ./server/livereload.go
라이브 리로딩 구현을 위한 코드가 위치해있습니다.

```go
// ./server/server.go
// --watch=true(기본) 일 때 livereload 활성화

if s.Watch {
  go watch()
  http.HandleFunc("/livereload", s.liveReloadHandler)
  // + s.Watch 유무에 따라 html에 livereload.js 링킹
}
```

라이브 리로딩을 위해서는 감시할 파일 상대경로를 합친 /livereload?relative-path=(감시할 파일 상대경로) 주소로 웹소켓을 연결합니다. 감시하는 파일이 변경될 경우 해당 웹소켓으로 메세지를 보냅니다. 이를 받은 브라우저는 자동으로 새로고침하여 업데이트 합니다.

파일 감시의 경우에는 매 1초마다 감시하는 파일이 수정됐는지 여부를 비교하여 확인하는 것으로 구현되었습니다.

### ./server/markdown.go
로컬 파일, 외부 이미지(http://..) 모두 링크 할 수 있도록 구현하였습니다.

특히 로컬 파일의 경우에는 주로 마크다운 파일 위치를 기준으로 상대 경로이므로  `imgRepathize()` 함수를 통해  /files?src=(절대 경로) 형식으로 변경합니다. 다만 이를 통해 모든 파일에 접근 할 수 있으므로 localhost 사용을 권장합니다.


