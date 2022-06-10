#!/bin/sh
telegram_bot_token="$1"
cd terraform
terraform apply -destroy --var "telegram_bot_token=$telegram_bot_token"
cd ..