# verdure-lambdas

###To build everything:
> go build -v ./...

### To run tests:
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