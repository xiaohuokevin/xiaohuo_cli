package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
)

func main() {
	app := cli.NewApp()

	app.Name = "xiaohuo"
	app.Usage = "a cli for me!"
	app.Version = "0.0.1"
	app.Author = "xiaohuo"
	app.Email = "xiaohuokevin@gmail.com"

	//app.Flags = []cli.Flag {
	//	cli.StringFlag{
	//		Name: "upload",
	//		Usage: "Upload file to aLi server",
	//	},
	//}

	app.Commands = []cli.Command {
		{
			Name: "upload",
			Aliases: []string {"up"},
			Usage: "Upload file to ali server",
			Action: func(c *cli.Context)  error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func upload()  {
	
}
