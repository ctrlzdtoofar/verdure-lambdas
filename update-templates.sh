#!/bin/zsh
aws ses update-template --cli-input-json file://templates/email_confirmation_en.json
aws ses update-template --cli-input-json file://templates/email_confirmation_es.json
aws ses update-template --cli-input-json file://templates/email_confirmation_de.json

aws ses update-template --cli-input-json file://templates/reset_confirmation_en.json
aws ses update-template --cli-input-json file://templates/reset_confirmation_es.json
aws ses update-template --cli-input-json file://templates/reset_confirmation_de.json