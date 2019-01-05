package bastion

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

func (app *ShellApp) AclLsCommand() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "Removes one or more ACLs",
		ArgsUsage: "ACL...",
		Action:    app.AclRm(),
	}
}
func (app *ShellApp) AclLs() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if err := app.myself.CheckRoles([]string{"admin"}); err != nil {
			return err
		}

		var acls []*dbmodels.ACL
		query := db.Order("created_at desc").Preload("UserGroups").Preload("HostGroups")
		if c.Bool("latest") {
			var acl dbmodels.ACL
			if err := query.First(&acl).Error; err != nil {
				return err
			}
			acls = append(acls, &acl)
		} else {
			if err := query.Find(&acls).Error; err != nil {
				return err
			}
		}
		if c.Bool("quiet") {
			for _, acl := range acls {
				fmt.Fprintln(app.Writer, acl.ID)
			}
			return nil
		}

		table := tablewriter.NewWriter(app.Writer)
		table.SetHeader([]string{"ID", "Weight", "User groups", "Host groups", "Host pattern", "Action", "Updated", "Created", "Comment"})
		table.SetBorder(false)
		table.SetCaption(true, fmt.Sprintf("Total: %d ACLs.", len(acls)))
		for _, acl := range acls {
			userGroups := []string{}
			hostGroups := []string{}
			for _, entity := range acl.UserGroups {
				userGroups = append(userGroups, entity.Name)
			}
			for _, entity := range acl.HostGroups {
				hostGroups = append(hostGroups, entity.Name)
			}

			table.Append([]string{
				fmt.Sprintf("%d", acl.ID),
				fmt.Sprintf("%d", acl.Weight),
				strings.Join(userGroups, ", "),
				strings.Join(hostGroups, ", "),
				acl.HostPattern,
				acl.Action,
				humanize.Time(acl.UpdatedAt),
				humanize.Time(acl.CreatedAt),
				acl.Comment,
			})
		}
		table.Render()
		return nil
	}
}
