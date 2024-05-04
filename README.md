# verdure-lambdas

###To build everything:
> go build -v ./...

### To run tests:
> source .env
> go test -v ./...

### Build ses lambda locally
> go build -o ses ./cmd/ses/sqs_ses.go

### Build and create/deploy the ses lambda, one time for each function
```zsh
./create.sh
```

### Build and update/deploy the ses lambda
```zsh
./update.sh
```

### SES Email Templates:
```zsh
./update-templates.sh
```

### AWS CloudWatch
https://us-east-2.console.aws.amazon.com/cloudwatch/home?region=us-east-2#logsV2:log-groups/log-group/$252Faws$252Flambda$252FSesConfirm
