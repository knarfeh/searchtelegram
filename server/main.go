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
	ESHOSTPORT, _ := conf.String("ESHOSTPORT")
	REDISHOST, _ := conf.String("REDISHOST")
	REDISPORT, _ := conf.String("REDISPORT")
	AWSACCESSKEY, _ := conf.String("AWSACCESSKEY")
	AWSSECRETKEY, _ := conf.String("AWSSECRETKEY")
	AWSREGION, _ := conf.String("AWSREGION")

	fmt.Println("Diagnose...")
	hauler, _ := CreateConsumer(
		ESHOSTPORT,
		REDISHOST,
		REDISPORT,
		AWSACCESSKEY,
		AWSSECRETKEY,
		AWSREGION,
	)
	redisClient := hauler.redisClient
	esClient := hauler.esClient
	reporter, _ := diagnose.New()
	reporter.Add(redisClient)
	reporter.Add(esClient)
	reporterResult := reporter.Check()
	fmt.Println(reporterResult)

	hauler.Query2ES()
}
