package main

import (
	"log"
	"os"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	style := styles.Get("github")
	if style == nil {
		log.Fatal("스타일을 찾을 수 없습니다.")
	}

	formatter := html.New(html.WithClasses(true))

	f, err := os.Create("server/web/assets/code-styles.css")
	if err != nil {
		log.Fatalf("파일 생성 오류: %v", err)
	}
	defer f.Close()

	f.WriteString("pre { overflow-x: auto; }\n")

	err = formatter.WriteCSS(f, style)
	if err != nil {
		log.Fatalf("CSS 작성 오류: %v", err)
	}

	log.Println("CSS 파일이 성공적으로 생성되었습니다.")
}
