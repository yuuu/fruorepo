package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fruorepo"
	app.Usage = "We enjoy repository."
	app.Version = "0.1.0"
	app.UsageText = "fruorepo [options]"
	app.Author = "k-masatany"
	app.Email = "sonntag902@gmail.com"
	licence, _ := ioutil.ReadFile("LICENSE")
	app.Copyright = string(licence)

	app.Action = func(c *cli.Context) error {
		run(c.Bool("dry-run"))
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

func run(dry bool) {
	f := new(Fruorepo)
	f.SelectRepository()
	f.PrintRepositoryOverview()

	prompt := promptui.Prompt{
		Label:     "Do you want to modify this repository?",
		IsConfirm: true,
	}
	accept, err := prompt.Run()
	if err != nil {
		MessageAndDie(fmt.Sprintf("Prompt failed %v\n", err))
	} else if accept != "y" {
		MessageAndDie("Exiting on user command")
	}

	labels, err := f.GetLabels()
	if err != nil {
		MessageAndDie(err.Error())
	}
	for _, label := range labels {
		fmt.Println(label.GetName())
	}
}

// MessageAndDie
func MessageAndDie(s string) {
	fmt.Println(s)
	os.Exit(-1)
}
