package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
	"time"

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
		{
			Name:   "download_cert",
			Usage:  "Download certificate file for https",
			Action: Download,
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
	hauler, _ := CreateConsumer(conf)

	fmt.Println("Diagnose...")
	redisClient := hauler.redisClient
	redisearchClient := hauler.redisearchClient
	esClient := hauler.esClient
	reporter, _ := diagnose.New()
	reporter.Add(redisClient)
	reporter.Add(redisearchClient)
	reporter.Add(esClient)
	reporterResult := reporter.Check()
	fmt.Println(reporterResult)

	go hauler.Submit2ES()
	go hauler.Search2Redisearch()
	go hauler.Scheduler(30 * time.Second)
	hauler.tb.Start()
}

// RunTelebot ...
func RunTelebot(c *cli.Context) {
	conf, err := config.ParseYaml(confString)
	conf.Env()
	Must(err)

	tgBot, _ := CreateTeleBot(conf)

	fmt.Println("Diagnose...")
	redisClient := tgBot.redisClient
	redisearchClient := tgBot.redisearchClient
	esClient := tgBot.esClient
	reporter, _ := diagnose.New()
	reporter.Add(redisClient) // TODO: add telebot Client???
	reporter.Add(redisearchClient)
	reporter.Add(esClient)
	reporterResult := reporter.Check()
	fmt.Println(reporterResult)

	tgBot.tb.Start()
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// Donwload ...
func Download(c *cli.Context) {
	conf, err := config.ParseYaml(confString)
	conf.Env()
	Must(err)

	AWSACCESSKEY, _ := conf.String("AWSACCESSKEY")
	AWSSECRETKEY, _ := conf.String("AWSSECRETKEY")
	AWSREGION_WITHSPACE, _ := conf.String("AWSREGION")
	AWSREGION := strings.TrimSpace(AWSREGION_WITHSPACE)
	creds := credentials.NewStaticCredentials(strings.TrimSpace(AWSACCESSKEY), strings.TrimSpace(AWSSECRETKEY), "")
	_, err = creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s\n", err)
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      &AWSREGION,
	})

	file, err := os.Create("searchtelegramdotcom.key")
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("knarfeh-website-certs"),
			Key:    aws.String("searchtelegramdotcom.key"),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", "searchtelegramdotcom.key", err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	file, err = os.Create("searchtelegramdotcom_bundle.crt")
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()
	numBytes, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("knarfeh-website-certs"),
			Key:    aws.String("searchtelegramdotcom_bundle.crt"),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", "searchtelegramdotcom_bundle.crt", err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
