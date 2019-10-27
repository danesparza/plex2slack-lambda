package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/danesparza/plex2slack-lambda/data"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
		mr := multipart.NewReader(strings.NewReader(string(s)), params["boundary"])
		log.Printf("DEBUG:: boundary: %v\n", params["boundary"])

		for {
			part, err := mr.NextPart()

			//	We're done parsing
			if err == io.EOF {
				log.Printf("PARSE EOF\n")
				break
			}

			//	We have an error
			if err != nil {
				log.Printf("There was an error: %v", err)
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Headers: map[string]string{
						"Content-Type":                "application/json",
						"Access-Control-Allow-Origin": "*",
					},
					Body: fmt.Sprintf("Error while parsing: %v\n", err),
				}, nil
			}

			//	If we have the right part (payload) parse the JSON request and log what we found
			if part.FormName() == "payload" {
				partBytes, _ := ioutil.ReadAll(part)
				msg := data.PlexMessage{}
				if err := json.Unmarshal(partBytes, &msg); err != nil {
					log.Printf("There was an error parsing JSON payload: %v\n", err)
				}
				log.Printf("Event: %v / Title: %v \n", msg.Event, msg.Metadata.Title)
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
