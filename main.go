// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"

// 	"github.com/gliderlabs/ssh"
// )

// func main() {
// 	ssh.Handle(func(s ssh.Session) {
// 		io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
// 	})

// 	log.Println("starting ssh server on port 3000...")
// 	log.Fatal(ssh.ListenAndServe(":3000", nil))
// }
