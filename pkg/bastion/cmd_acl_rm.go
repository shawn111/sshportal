package bastion

import (
	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

func (app *ShellApp) AclRmCommand() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "Removes one or more ACLs",
		ArgsUsage: "ACL...",
		Action:    app.AclRm(),
	}
}

func (app *ShellApp) AclRm() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		if err := app.myself.CheckRoles([]string{"admin"}); err != nil {
			return err
		}

		return dbmodels.ACLsByIdentifiers(db, c.Args()).Delete(&dbmodels.ACL{}).Error
	}
}
