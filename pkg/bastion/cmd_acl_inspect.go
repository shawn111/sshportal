package bastion

import (
	"encoding/json"

	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

func (app *ShellApp) AclInspectCommand() cli.Command {
	return cli.Command{
		Name:      "inspect",
		Usage:     "Shows detailed information on one or more ACLs",
		ArgsUsage: "ACL...",
		Action:    app.AclInspect(),
	}

}
func (app *ShellApp) AclInspect() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		if err := app.myself.CheckRoles([]string{"admin"}); err != nil {
			return err
		}

		var acls []dbmodels.ACL
		if err := dbmodels.ACLsPreload(dbmodels.ACLsByIdentifiers(db, c.Args())).Find(&acls).Error; err != nil {
			return err
		}

		enc := json.NewEncoder(app.Writer)
		enc.SetIndent("", "  ")
		return enc.Encode(acls)
	}
}
