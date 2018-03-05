package main

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

type Options struct {
	DryRun bool
}

func main() {
	// コマンドツールの基本情報
	app := cli.NewApp()
	app.Name = "fruorepo"
	app.Usage = "We enjoy repository."
	app.Version = "0.1.0"
	app.UsageText = "fruorepo [options]"
	app.Author = "k-masatany"
	app.Email = "sonntag902@gmail.com"
	app.Copyright = `
	MIT License

	Copyright (c) 2018 Kensuke Masatani

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
`
	opt := new(Options)
	app.Action = func(c *cli.Context) error {
		opt.DryRun = c.Bool("dry-run")
		run(opt)
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

func run(opt *Options) {
	f := new(Fruorepo)
	f.SelectRepository()
	f.PrintRepositoryOverview()

	accept := confirm("Do you want to modify this repository?")
	if accept != "y" {
		MessageAndDie("Exiting on user command")
	}

	// FetchLabels
	labels, err := f.FetchLabels()
	if err != nil {
		MessageAndDie(err.Error())
	}

	// DeleteLabels
	for _, label := range labels {
		err := f.DeleteLabel(label, opt)
		if err != nil {
			MessageAndDie(err.Error())
		}
	}

	// CreateLabels
	labelSet := []map[string]string{}
	labelSet = append(labelSet, map[string]string{"name": ":dragon:ドラゴン", "color": "FF6969"})
	labelSet = append(labelSet, map[string]string{"name": ":crossed_swords:討伐中", "color": "6EC4FF"})

	for _, label := range labelSet {
		err := f.CreateLabel(label["name"], label["color"], opt)
		if err != nil {
			MessageAndDie(err.Error())
		}
	}

	fmt.Println("Repository has been changed.\nHave a nice day!!")
}

// MessageAndDie
func MessageAndDie(s string) {
	fmt.Println(s)
	os.Exit(-1)
}

// bool prompt wrapper
func confirm(message string) string {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	ret, err := prompt.Run()
	if err != nil {
		MessageAndDie(fmt.Sprintf("Prompt failed %v\n", err))
	}

	return ret
}
