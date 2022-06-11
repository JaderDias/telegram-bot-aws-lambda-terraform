#!/bin/sh
telegram_bot_tokens="$1"
cd terraform
terraform apply -destroy --var "telegram_bot_tokens=$telegram_bot_tokens"
cd ..