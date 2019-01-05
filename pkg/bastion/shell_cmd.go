package bastion

import (
	"github.com/urfave/cli"
	"moul.io/sshportal/pkg/dbmodels"
)

type ShellApp struct {
	*cli.App

	actx   *authContext
	myself *dbmodels.User
}

func NewApp(actx *authContext) *ShellApp {
	app := cli.NewApp()
	return &ShellApp{app, actx, &actx.user}
}

//func (app *ShellCmd) NewCommand(name string) cli.Command {
//}
