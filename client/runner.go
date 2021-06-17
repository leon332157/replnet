package client

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecCommand(cmdString string) {
	cmds := strings.Split(cmdString, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
  


	
	err := cmd.Start()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
  //err = cmd.Wait()
	fmt.Printf("Program exited: %v", cmd)
}

func printStdout(){

}

/*func readData() {
  
}*/