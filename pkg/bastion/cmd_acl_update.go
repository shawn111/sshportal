package bastion

import (
	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

func (app *ShellApp) AclUpdateCommand() cli.Command {
	return cli.Command{
		Name:      "update",
		Usage:     "Updates an existing acl",
		ArgsUsage: "ACL...",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "action, a", Usage: "Update action"},
			cli.StringFlag{Name: "pattern, p", Usage: "Update host-pattern"},
			cli.UintFlag{Name: "weight, w", Usage: "Update weight"},
			cli.StringFlag{Name: "comment, c", Usage: "Update comment"},
			cli.StringSliceFlag{Name: "assign-usergroup, ug", Usage: "Assign the ACL to new `USERGROUPS`"},
			cli.StringSliceFlag{Name: "unassign-usergroup", Usage: "Unassign the ACL from `USERGROUPS`"},
			cli.StringSliceFlag{Name: "assign-hostgroup, hg", Usage: "Assign the ACL to new `HOSTGROUPS`"},
			cli.StringSliceFlag{Name: "unassign-hostgroup", Usage: "Unassign the ACL from `HOSTGROUPS`"},
		},
		Action: app.AclUpdate(),
	}
}

func (app *ShellApp) AclUpdate() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}
		if err := app.myself.CheckRoles([]string{"admin"}); err != nil {
			return err
		}

		var acls []dbmodels.ACL
		if err := dbmodels.ACLsByIdentifiers(db, c.Args()).Find(&acls).Error; err != nil {
			return err
		}

		tx := db.Begin()
		for _, acl := range acls {
			model := tx.Model(&acl)
			update := dbmodels.ACL{
				Action:      c.String("action"),
				HostPattern: c.String("pattern"),
				Weight:      c.Uint("weight"),
				Comment:     c.String("comment"),
			}
			if err := model.Updates(update).Error; err != nil {
				tx.Rollback()
				return err
			}

			// associations
			var appendUserGroups []dbmodels.UserGroup
			var deleteUserGroups []dbmodels.UserGroup
			if err := dbmodels.UserGroupsByIdentifiers(db, c.StringSlice("assign-usergroup")).Find(&appendUserGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := dbmodels.UserGroupsByIdentifiers(db, c.StringSlice("unassign-usergroup")).Find(&deleteUserGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := model.Association("UserGroups").Append(&appendUserGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if len(deleteUserGroups) > 0 {
				if err := model.Association("UserGroups").Delete(deleteUserGroups).Error; err != nil {
					tx.Rollback()
					return err
				}
			}

			var appendHostGroups []dbmodels.HostGroup
			var deleteHostGroups []dbmodels.HostGroup
			if err := dbmodels.HostGroupsByIdentifiers(db, c.StringSlice("assign-hostgroup")).Find(&appendHostGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := dbmodels.HostGroupsByIdentifiers(db, c.StringSlice("unassign-hostgroup")).Find(&deleteHostGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := model.Association("HostGroups").Append(&appendHostGroups).Error; err != nil {
				tx.Rollback()
				return err
			}
			if len(deleteHostGroups) > 0 {
				if err := model.Association("HostGroups").Delete(deleteHostGroups).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit().Error
	}
}
