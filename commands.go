package main

import (
	"github.com/reizist/ec2sshtb/utils"
	"github.com/urfave/cli"
)

var commandSync = cli.Command{
	Name:        "sync",
	Usage:       "Sync instances to local file ~/.ec2sshtb",
	Description: "",
	Action:      doSync,
}

var commandSSH = cli.Command{
	Name:        "ssh",
	Usage:       "ssh to instance selected by peco",
	Description: "",
	Action:      doSSH,
}

var Commands = []cli.Command{
	commandSync,
	commandSSH,
}

func doSync(c *cli.Context) error {
	utils.Sync()
	return nil
}

func doSSH(c *cli.Context) error {
	utils.SSH()
	return nil
}
