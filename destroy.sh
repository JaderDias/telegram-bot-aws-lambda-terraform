#!/bin/sh
environment="$1"
if [ -z "$environment" ]
then
    echo "Usage: deploy.sh <environment>"
    exit 1
fi

telegram_bot_tokens=`aws ssm get-parameter --name "${environment}_telegram_bot_tokens" --output text --with-decryption | cut -f7`
if [ -z "$telegram_bot_tokens" ]
then
    printf "paste the telegram_bot_tokens JSON: "
    read telegram_bot_tokens
fi

cd terraform

terraform workspace new $environment
terraform workspace select $environment

terraform apply -destroy --var "telegram_bot_tokens=$telegram_bot_tokens"
cd ..