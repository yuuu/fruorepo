package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"golang.org/x/oauth2"
	"github.com/briandowns/spinner"
)

type Fruorepo struct {
	Token      string
	Client     *github.Client
	Repository *github.Repository
}

func (f *Fruorepo) SetTokenFromEnv() {
	f.Token = os.Getenv("GITHUB_AUTH_TOKEN")
}

func (f *Fruorepo) SetClient() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: f.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	f.Client = github.NewClient(tc)
}

// SelectRepository get taraget *github.Repository
func (f *Fruorepo) SelectRepository() {
	fmt.Println("access github.com...")
	s := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	s.Start()

	f.SetTokenFromEnv()
	f.SetClient()
	// get all pages of results
	repositories, _, err := f.Client.Repositories.List(context.Background(), "", nil)
	if err != nil {
		MessageAndDie(err.Error())
	}

	var items = []string{}
	for _, r := range repositories {
		items = append(items, r.GetFullName())
	}
	s.Stop()

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
		MessageAndDie(fmt.Sprintf("Prompt failed %v\n", err))
	}

	f.Repository = repositories[n]
}

// PrintRepositoryOverview
func (f *Fruorepo) PrintRepositoryOverview() {
	fmt.Println("------------------------------------------------")
	fmt.Println("       Name:", f.Repository.GetName())
	fmt.Println("        URL:", f.Repository.GetURL())
	fmt.Println("Description:", f.Repository.GetDescription())
	fmt.Println("------------------------------------------------")
}

func ChangeLabel(repo *github.Repository) {}

// GetLabels fetches all labels in a repository, iterating over pages for 50 at a time.
func (f *Fruorepo) GetLabels() ([]*github.Label, error) {
	var labelsRemote []*github.Label

	pagination := &github.ListOptions{
		PerPage: 50,
		Page:    1,
	}

	for {
		fmt.Printf("Fetching labels from Github, page %d\n", pagination.Page)

		labels, resp, err := f.Client.Issues.ListLabels(
			context.Background(),
			f.Repository.GetOwner().GetLogin(),
			f.Repository.GetName(),
			pagination,
		)

		if err != nil {
			fmt.Println("Failed to fetch labels from Github")
			return nil, err
		}

		labelsRemote = append(labelsRemote, labels...)

		if resp.NextPage == 0 {
			fmt.Println("Fetched all labels from Github")
			break
		}
		pagination.Page = resp.NextPage
	}

	return labelsRemote, nil
}
