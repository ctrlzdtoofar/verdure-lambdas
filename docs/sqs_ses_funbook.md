# Troubleshooting Confirmation Messages

1. Is there a new user in the database?
2. Did the message make it to SQS?
3. Was the message read by the lambda?
4. Is the SQS Intake Disabled?
4. Does CloudWatch have a log from the lambda?
5. Are there enough permissions on the lambda's role?