package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/knarfeh/searchtelegram/server/diagnose"
	"github.com/olebedev/config"
)

func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	Run(os.Args)
}

// Run creates, configures and runs
// main cli.App
func Run(args []string) {
	app := cli.NewApp()
	app.Name = "app"
	app.Usage = "Search telegram"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Runs server",
			Action: RunServer,
		},
		{
			Name:   "worker",
			Usage:  "Runs worker",
			Action: RunWorker,
		},
		{
			Name:   "telebot",
			Usage:  "Runs telegram bot",
			Action: RunTelebot,
		},
	}
	app.Run(args)
}

// RunServer creates, configures and runs
// main server.App
func RunServer(c *cli.Context) {
	app := NewApp(AppOptions{
	// see server/app.go:150
	})
	app.Run()
}

// RunWorker runs worker
func RunWorker(c *cli.Context) {
	conf, err := config.ParseYaml(confString)
	conf.Env()
	Must(err)

	fmt.Println("Diagnose...")
	hauler, _ := CreateConsumer(
		conf,
	)
	redisClient := hauler.redisClient
	esClient := hauler.esClient
	reporter, _ := diagnose.New()
	reporter.Add(redisClient)
	reporter.Add(esClient)
	reporterResult := reporter.Check()
	fmt.Println(reporterResult)

	hauler.Query2ES()
	hauler.tb.Start()
}

// RunTelebot ...
func RunTelebot(c *cli.Context) {
	conf, err := config.ParseYaml(confString)
	conf.Env()
	Must(err)

	TGBOTTOKEN, _ := conf.String("TGBOTTOKEN")
	REDISHOST, _ := conf.String("REDISHOST")
	REDISPORT, _ := conf.String("REDISPORT")

	tgBot, _ := CreateTeleBot(TGBOTTOKEN, REDISHOST, REDISPORT)

	redisClient := tgBot.redisClient
	reporter, _ := diagnose.New()
	reporter.Add(redisClient)

	reporterResult := reporter.Check()
	fmt.Println(reporterResult)

	tgBot.tb.Start()

}
