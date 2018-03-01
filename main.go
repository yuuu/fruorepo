package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	app := cli.NewApp()
	app.Name = "sampleApp"
	app.Usage = "This app echo input arguments"
	app.Version = "0.0.1"

	app.Action = func(context *cli.Context) error {
		if context.Bool("token") {
			repo, err := command(context.Args().Get(0))
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, r := range repo {
				fmt.Println(r.GetFullName())
			}
		} else {
			fmt.Println(context.Args().Get(0))
		}
		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "token, t",
			Usage: "GitHub Access Token",
		},
	}

	app.Run(os.Args)
}

func command(token string) ([]*github.Repository, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)

	return repos, err
}
