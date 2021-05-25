package server

/*
import (
	"io"
	"github.com/gliderlabs/ssh"
	"fmt"
)

func StartSSHServer() {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})

	fmt.Println(ssh.ListenAndServe(":2222", nil))
}

// https://github.com/leechristensen/GolangSSHServer
*/