package ssh

import (
	"fmt"
	"github.com/sfreiberg/simplessh"
)

var client *simplessh.Client

func Conn(host *string, user *string, pass *string) {
	/*
		Leave privKeyPath empty to use $HOME/.ssh/id_rsa.
		If username is blank simplessh will attempt to use the current user.
	*/
	//client, err := simplessh.ConnectWithKeyFile("localhost:22", "root", "/home/user/.ssh/id_rsa")
	var err error
	client, err = simplessh.ConnectWithPassword(*host, *user, *pass)
	if err != nil {
		defer client.Close()
		panic(err)
	}
}

func Close() {
	client.Close()
}

func Run(cmd *string) {
	output, err := client.Exec(*cmd)
	if err != nil {
		//panic(err)
		fmt.Printf("Error from command: %s\n", err)
	}
	//fmt.Printf("Run command: %s\n", *cmd)
	fmt.Printf("Output:\n %s\n", output)
}

func Download(remote string, local string) {
	err := client.Download(remote, local)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download: %s -> %s\n", remote, local)
}
