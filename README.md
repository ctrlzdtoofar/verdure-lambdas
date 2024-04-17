# verdure-lambdas

###To build everything:
> go build -v ./...

### To run tests:
> go test -v ./...

### Build ses lambda locally
> go build -o ses ./cmd/ses/sqs_ses.go

### Build and create/deploy the ses lambda, one time for each function
```zsh
source .env
[ -f bootstrap ] && rm bootstrap

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/ses/sqs_ses.go
chmod 755 bootstrap

[ -f function.zip ] && rm function.zip
zip function.zip bootstrap

aws lambda delete-function --function-name SesConfirm
aws lambda create-function --function-name SesConfirm \
    --zip-file fileb://function.zip --handler bootstrap \
    --architectures arm64 \
    --runtime provided.al2023 \
    --role ${CONFIRM_ROLE_ARN}

aws lambda create-event-source-mapping \
    --function-name SesConfirm \
    --batch-size 10 \
    --event-source-arn ${CONFIRM_SQS_ARN}
    
aws logs create-log-group --log-group-name /aws/lambda/SesConfirm
aws logs put-retention-policy --log-group-name /aws/lambda/SesConfirm --retention-in-days 3

# verification
aws lambda get-function --function-name SesConfirm 
aws lambda list-event-source-mappings
aws logs describe-log-groups --log-group-name-prefix /aws/lambda/SesConfirm
```

### Build and update/deploy the ses lambda
```zsh
./update.sh
```

### Submit SES Email Templates:
```zsh
aws ses create-template --cli-input-json file://templates/email_confirmation.json
aws ses create-template --cli-input-json file://templates/reset_confirmation.json

aws ses update-template --cli-input-json file://templates/email_confirmation.json
aws ses update-template --cli-input-json file://templates/reset_confirmation.json
# or
./update-templates.sh

aws ses list-templates 

aws ses delete-template --template-name EmailConfirmation
aws ses delete-template --template-name PasswordResetConfirmation
```