package bastion

import (
	"github.com/urfave/cli"
)

func (app *ShellApp) AclCommand() cli.Command {
	return cli.Command{
		Name:  "acl",
		Usage: "Manages ACLs",
		Subcommands: []cli.Command{
			app.AclCreateCommand(),
			app.AclInspectCommand(),
			app.AclLsCommand(),
			app.AclRmCommand(),
			app.AclUpdateCommand(),
			app.AclRmCommand(),
		},
	}
}
