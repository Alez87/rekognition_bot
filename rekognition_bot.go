package main

import (
	"fmt"
        "log"
        "os"
        "encoding/json"
	"github.com/yanzay/tbot"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type Configuration struct {
	Telegram_token	string
	Region		string
	Bucket_name	string
	File_name	string
	MaxLabels	int64
        MinConfidence	float64
}

var configuration Configuration

func load_configuration() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
	  fmt.Println("error:", err)
	}
	fmt.Println(configuration.Telegram_token)
	fmt.Println(configuration.Region)
	fmt.Println(configuration.Bucket_name)
	fmt.Println(configuration.File_name)
	fmt.Println(configuration.MaxLabels)
	fmt.Println(configuration.MinConfidence)
}

//func exitErrorf(msg string, args ...interface{}) {
//    fmt.Fprintf(os.Stderr, msg+"\n", args...)
//    os.Exit(1)
//}

func videoHandler(m *tbot.Message){
	m.Reply("Functionality not implemented yet")
	// TODO
}

func imageHandler(m *tbot.Message) {
	m.Reply("Analysis started")

	session, _ := session.NewSession(&aws.Config{Region: aws.String(configuration.Region)},)
        svc := rekognition.New(session)
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(configuration.Bucket_name),
				Name: aws.String(configuration.File_name),
			},
		},
		MaxLabels:     aws.Int64(configuration.MaxLabels),
		MinConfidence: aws.Float64(configuration.MinConfidence),
	}

	result, err := svc.DetectLabels(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidS3ObjectException:
				fmt.Println(rekognition.ErrCodeInvalidS3ObjectException, aerr.Error())
			case rekognition.ErrCodeInvalidParameterException:
				fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				fmt.Println(rekognition.ErrCodeImageTooLargeException, aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				fmt.Println(rekognition.ErrCodeInvalidImageFormatException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Print("result:")
	fmt.Print(result)
        m.Reply("Finished")
}

func main() {
	load_configuration()
//        bot, err := tbot.NewServer(os.Getenv(configuration.Telegram_token))
	bot, err := tbot.NewServer(os.Getenv("TOKEN"))
        if err != nil {
		fmt.Print(err)
                log.Fatal(err)
        }
	bot.HandleFunc("/image", imageHandler)
	bot.HandleFunc("/video", videoHandler)
        bot.ListenAndServe()
}
