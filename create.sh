#!/bin/zsh

source .env
[ -f bootstrap ] && rm bootstrap

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/ses/sqs_ses.go
chmod 755 bootstrap

[ -f function.zip ] && rm function.zip
zip function.zip bootstrap

aws lambda create-function --function-name SesConfirm \
    --zip-file fileb://function.zip --handler bootstrap \
    --architectures arm64 \
    --runtime provided.al2023 \
    --role "${CONFIRM_ROLE_ARN}" \

aws lambda update-function-configuration --function-name SesConfirm \
    --dead-letter-config TargetArn="${CONFIRM_DLQ_SQS_ARN}"

aws lambda create-event-source-mapping \
    --function-name SesConfirm \
    --batch-size 10 \
    --event-source-arn "${CONFIRM_SQS_ARN}"

aws logs create-log-group --log-group-name /aws/lambda/SesConfirm
aws logs put-retention-policy --log-group-name /aws/lambda/SesConfirm --retention-in-days 3

# Create email templates
aws ses create-template --cli-input-json file://templates/email_confirmation_en.json
aws ses create-template --cli-input-json file://templates/email_confirmation_es.json
aws ses create-template --cli-input-json file://templates/email_confirmation_de.json

aws ses create-template --cli-input-json file://templates/reset_confirmation_en.json
aws ses create-template --cli-input-json file://templates/reset_confirmation_es.json
aws ses create-template --cli-input-json file://templates/reset_confirmation_de.json


