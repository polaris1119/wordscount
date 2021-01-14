package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/polaris1119/wordscount"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "wordscount",
		Usage: "统计文章字数",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "The article file `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			file := c.String("file")
			if file == "" {
				return cli.ShowAppHelp(c)
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			counter := &wordscount.Counter{}
			counter.Stat(string(content))
			fmt.Printf("%#v\n", counter)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
