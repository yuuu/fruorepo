package main

import (
	"context"
	"fmt"
	"log"
	"os"
	_ "strconv"

	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	app := cli.NewApp()
	app.Name = "fruorepo"
	app.Usage = "We enjoy repository."
	app.Version = "0.1.0"

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	app.Action = func(c *cli.Context) error {
		command(token, c.Bool("dry-run"))
		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Run, But Not Change Repository.",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func selectRepository(client *github.Client) *github.Repository {
	// get all pages of results
	repositories, _, err := client.Repositories.List(context.Background(), "", nil)
	if err != nil {
		MessageAndDie(err.Error())
	}

	var items = []string{}
	for _, r := range repositories {
		items = append(items, r.GetFullName())
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ . | cyan }}",
		Inactive: "  {{ .FullName | cyan }}",
		Selected: "\U00002728 {{ . | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select Repository",
		Items:     items,
		Templates: templates,
	}

	n, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	fmt.Printf("You choose %s\n", repositories[n].GetFullName())
	return repositories[n]
}

func command(token string, dry bool) {
	client := getClient(token)

	selectRepository(client)

	// PrintRepositoryOverview(repo)
	// accept := prompter.YesNo("Do you want to modify this repository?", false)

	// if accept {
	// 	labels, err := GetLabels(client, target)
	// 	if err != nil {
	// 		MessageAndDie(err.Error())
	// 	}
	// 	for _, label := range labels {
	// 		fmt.Println(label.GetName())
	// 	}
	// } else {
	// 	fmt.Print("Exiting on user command")
	// }
}

func getClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

// MessageAndDie
func MessageAndDie(s string) {
	fmt.Println(s)
	os.Exit(-1)
}

// PrintRepositoryOverview
func PrintRepositoryOverview(repo *github.Repository) {
	fmt.Println("------------------------------------------------")
	fmt.Println("       Name:", repo.GetName())
	fmt.Println("        URL:", repo.GetURL())
	fmt.Println("Description:", repo.GetDescription())
	fmt.Println("------------------------------------------------")
}

func ChangeLabel(repo *github.Repository) {}

// GetLabels fetches all labels in a repository, iterating over pages for 50 at a time.
func GetLabels(client *github.Client, repo *github.Repository) ([]*github.Label, error) {
	var labelsRemote []*github.Label

	pagination := &github.ListOptions{
		PerPage: 50,
		Page:    1,
	}

	for {
		fmt.Printf("Fetching labels from Github, page %d\n", pagination.Page)

		labels, resp, err := client.Issues.ListLabels(
			context.Background(),
			repo.GetOwner().GetName(),
			repo.GetName(),
			pagination,
		)

		if err != nil {
			fmt.Println("Failed to fetch labels from Github")
			return nil, err
		}
		fmt.Printf("Response: %s\n", resp)

		labelsRemote = append(labelsRemote, labels...)

		if resp.NextPage == 0 {
			fmt.Println("Fetched all labels from Github")
			break
		}
		pagination.Page = resp.NextPage
	}

	return labelsRemote, nil
}
