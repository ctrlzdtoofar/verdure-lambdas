#!/bin/zsh

aws ses update-template --cli-input-json file://templates/email_confirmation.json
aws ses update-template --cli-input-json file://templates/reset_confirmation.json
