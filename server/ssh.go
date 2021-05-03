package server

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
)

func StartSSHServer() {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
// https://github.com/leechristensen/GolangSSHServer