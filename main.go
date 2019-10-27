package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/danesparza/plex2slack-lambda/data"
)

//	Set up our flags
var (
	webhookURL = ""
)

func parseEnvironment() {

	//	Check for allowed origins
	if envWebhook := os.Getenv("PLEX2SLACK_WEBHOOK_URL"); envWebhook != "" {
		webhookURL = envWebhook
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//	Parse our options from the environment
	parseEnvironment()

	//	Verify the webhook url has been set

	//	Step zero: Decode the body:
	s, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Printf("There was an error base64 decoding the body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: fmt.Sprintf("Error getting content type: %v", err),
		}, nil
	}

	//	Print out what we have
	log.Printf("%+v\n", string(s))

	log.Printf("Headers: %+v\n", request.Headers)

	//	First, get the content type.  Throw an error if we can't
	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		log.Printf("There was an error getting the content type: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnsupportedMediaType,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: fmt.Sprintf("Error getting content type: %v", err),
		}, nil
	}

	//	If the content type is multipart...
	if strings.HasPrefix(mediaType, "multipart/") {

		//	Break apart the content into its parts
		readForm := multipart.NewReader(strings.NewReader(string(s)), params["boundary"])
		log.Printf("DEBUG:: boundary: %v\n", params["boundary"])

		for {
			part, errPart := readForm.NextPart()
			if errPart == io.EOF {
				break
			}
			if part.FormName() == "thumb" {
				partBytes, _ := ioutil.ReadAll(part)
				err := ioutil.WriteFile("thumb.jpg", partBytes, 0644)
				if err != nil {
					fmt.Printf("Error saving thumbnail: %v\n", err)
				}
			} else if part.FormName() == "payload" {
				partBytes, _ := ioutil.ReadAll(part)
				msg := data.PlexMessage{}
				if err := json.Unmarshal(partBytes, &msg); err != nil {
					panic(err)
				}

				//	First, see what kind of message it is:
				if msg.Event == "library.new" {

					//	Format our notification:
					slackMsg := data.SlackRequestBody{}

					//	Movie
					if msg.Metadata.Type == "movie" {
						slackMsg = data.SlackRequestBody{
							Text: fmt.Sprintf("%v added to Movies", msg.Metadata.Title),
							Blocks: []data.SlackBlock{
								data.SlackBlock{
									Type: "section",
									Text: &data.SlackText{
										Text: fmt.Sprintf("*%v*", msg.Metadata.Title),
										Type: "mrkdwn",
									},
								},
								data.SlackBlock{
									Type: "context",
									Elements: []data.SlackElement{
										data.SlackElement{
											Type: "mrkdwn",
											Text: "added to Movies",
										},
									},
								},
								data.SlackBlock{
									Type: "divider",
								},
							},
						}
					}

					//	TV show
					if msg.Metadata.Type == "episode" {
						slackMsg = data.SlackRequestBody{
							Text: fmt.Sprintf("%v %v: %v added to TV shows", msg.Metadata.GrandparentTitle, msg.Metadata.ParentTitle, msg.Metadata.Title),
							Blocks: []data.SlackBlock{
								data.SlackBlock{
									Type: "section",
									Text: &data.SlackText{
										Text: fmt.Sprintf("New episode of *%v %v*: _%v_", msg.Metadata.GrandparentTitle, msg.Metadata.ParentTitle, msg.Metadata.Title),
										Type: "mrkdwn",
									},
								},
								data.SlackBlock{
									Type: "context",
									Elements: []data.SlackElement{
										data.SlackElement{
											Type: "mrkdwn",
											Text: "added to TV shows",
										},
									},
								},
								data.SlackBlock{
									Type: "divider",
								},
							},
						}
					}

					//	Send the message to slack
					SendSlackNotification(webhookURL, slackMsg)
				}

				//	Used for debugging:
				log.Printf("Event: %v / Type: %v / Grandparent title: %v / Parent title: %v / Title: %v \n", msg.Event, msg.Metadata.Type, msg.Metadata.GrandparentTitle, msg.Metadata.ParentTitle, msg.Metadata.Title)
			}
		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: "Success",
	}, nil

}

func main() {
	lambda.Start(handler)
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhookURL string, msg data.SlackRequestBody) error {

	slackBody, _ := json.Marshal(msg)

	log.Printf("Sending message to Slack:\n %v\n", string(slackBody))

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
