# plex2slack-lambda [![CircleCI](https://circleci.com/gh/danesparza/plex2slack-lambda.svg?style=shield)](https://circleci.com/gh/danesparza/plex2slack-lambda)
Golang based AWS Lambda handler for sending plex webhook messages to a slack channel.  For a standalone service, see [plex2slack](https://github.com/danesparza/plex2slack)

## Quick start
### AWS Lambda setup
- Download the latest release and setup a new Go based lambda function in AWS.  The handler name is `plex2slack-lambda`
- Add the environment variable `PLEX2SLACK_WEBHOOK_URL` to your lambda handler.  Set this to your Slack incoming webhook url.  It should look something like *https://hooks.slack.com/services/SOMETHING/SOMEID/SOMETHINGELSEHERE*

### AWS API Gateway setup
- Create a new AWS API Gateway API.  
- In the 'Settings' section for the AWS API Gateway for the new API, make sure to add `multipart/form-data` under *Binary Media Types*
- Create a resource (call it something like 'plex') and create a `POST` method on the resource.  Wire up the lambda handler on the method.  Make sure to select *Use Lambda Proxy integration*
- Deploy the API
