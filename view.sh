#!/bin/zsh

source .env

echo "CONFIRM_ROLE_ARN: ${CONFIRM_ROLE_ARN}"
echo "CONFIRM_DLQ_SQS_ARN: ${CONFIRM_DLQ_SQS_ARN}"

# verification
aws lambda get-function --function-name SesConfirm
aws lambda list-event-source-mappings
aws logs describe-log-groups --log-group-name-prefix /aws/lambda/SesConfirm

# List email templates
aws ses list-templates