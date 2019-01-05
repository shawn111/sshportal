package bastion

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

func (app *ShellApp) AclCreateCommand() cli.Command {
	return cli.Command{
		Name:        "create",
		Usage:       "Creates a new ACL",
		Description: "$> acl create -",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "hostgroup, hg", Usage: "Assigns `HOSTGROUPS` to the acl"},
			cli.StringSliceFlag{Name: "usergroup, ug", Usage: "Assigns `USERGROUP` to the acl"},
			cli.StringFlag{Name: "pattern", Usage: "Assigns a host pattern to the acl"},
			cli.StringFlag{Name: "comment", Usage: "Adds a comment"},
			cli.StringFlag{Name: "action", Usage: "Assigns the ACL action (allow,deny)", Value: string(dbmodels.ACLActionAllow)},
			cli.UintFlag{Name: "weight, w", Usage: "Assigns the ACL weight (priority)"},
		},
		Action: app.AclCreate(),
	}
}

func (app *ShellApp) AclCreate() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if err := app.myself.CheckRoles([]string{"admin"}); err != nil {
			return err
		}
		acl := dbmodels.ACL{
			Comment:     c.String("comment"),
			HostPattern: c.String("pattern"),
			UserGroups:  []*dbmodels.UserGroup{},
			HostGroups:  []*dbmodels.HostGroup{},
			Weight:      c.Uint("weight"),
			Action:      c.String("action"),
		}
		if acl.Action != string(dbmodels.ACLActionAllow) && acl.Action != string(dbmodels.ACLActionDeny) {
			return fmt.Errorf("invalid action %q, allowed values: allow, deny", acl.Action)
		}
		if _, err := govalidator.ValidateStruct(acl); err != nil {
			return err
		}

		var userGroups []*dbmodels.UserGroup
		if err := dbmodels.UserGroupsPreload(dbmodels.UserGroupsByIdentifiers(db, c.StringSlice("usergroup"))).Find(&userGroups).Error; err != nil {
			return err
		}
		acl.UserGroups = append(acl.UserGroups, userGroups...)
		var hostGroups []*dbmodels.HostGroup
		if err := dbmodels.HostGroupsPreload(dbmodels.HostGroupsByIdentifiers(db, c.StringSlice("hostgroup"))).Find(&hostGroups).Error; err != nil {
			return err
		}
		acl.HostGroups = append(acl.HostGroups, hostGroups...)

		if len(acl.UserGroups) == 0 {
			return fmt.Errorf("an ACL must have at least one user group")
		}
		if len(acl.HostGroups) == 0 && acl.HostPattern == "" {
			return fmt.Errorf("an ACL must have at least one host group or host pattern")
		}

		if err := db.Create(&acl).Error; err != nil {
			return err
		}
		fmt.Fprintf(app.Writer, "%d\n", acl.ID)
		return nil
	}
}
