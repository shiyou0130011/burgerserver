package main

import (
	"io"
	"log"
	"net/http"

	"github.com/shiyou0130011/burgerserver"
)

func main() {
	s := burgerserver.DefaultServer()
	s.Logger.SetFlags(log.Llongfile | log.LUTC)

	s.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		io.WriteString(w, `{"status": "successfully"}`)
	})
	s.HandleFunc("/sample.txt", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Print("sample text")
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ornare lacus at erat mattis venenatis. Curabitur eleifend dolor varius est consectetur, eu gravida leo feugiat. `)
	})
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8"/>
				<title>Sample Page</title>
			</head>
			<body>
				<h1>Sample Page</h1>
				<article>Hello World!</article>
			</body>
		</html>`)
	})

	s.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// sample image from https://play.golang.org/
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00,
			0x10, 0x00, 0x00, 0x00, 0x0f, 0x04, 0x03, 0x00, 0x00, 0x00, 0x1f, 0x5d, 0x52, 0x1c, 0x00, 0x00, 0x00, 0x0f, 0x50,
			0x4c, 0x54, 0x45, 0x7a, 0xdf, 0xfd, 0xfd, 0xff, 0xfc, 0x39, 0x4d, 0x52, 0x19, 0x16, 0x15, 0xc3, 0x8d, 0x76, 0xc7,
			0x36, 0x2c, 0xf5, 0x00, 0x00, 0x00, 0x40, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x95, 0xc9, 0xd1, 0x0d, 0xc0, 0x20,
			0x0c, 0x03, 0xd1, 0x23, 0x5d, 0xa0, 0x49, 0x17, 0x20, 0x4c, 0xc0, 0x10, 0xec, 0x3f, 0x53, 0x8d, 0xc2, 0x02, 0x9c,
			0xfc, 0xf1, 0x24, 0xe3, 0x31, 0x54, 0x3a, 0xd1, 0x51, 0x96, 0x74, 0x1c, 0xcd, 0x18, 0xed, 0x9b, 0x9a, 0x11, 0x85,
			0x24, 0xea, 0xda, 0xe0, 0x99, 0x14, 0xd6, 0x3a, 0x68, 0x6f, 0x41, 0xdd, 0xe2, 0x07, 0xdb, 0xb5, 0x05, 0xca, 0xdb,
			0xb2, 0x9a, 0xdd, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	http.ListenAndServe(":8081", s.ToHttpHandler())
}
