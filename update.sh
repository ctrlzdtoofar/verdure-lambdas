#!/bin/zsh

[ -f bootstrap ] && rm bootstrap

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/ses/sqs_ses.go
chmod 755 bootstrap

[ -f function.zip ] && rm function.zip
zip function.zip bootstrap

aws lambda update-function-code --function-name SesConfirm --zip-file fileb://function.zip
