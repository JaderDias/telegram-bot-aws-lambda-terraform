#!/bin/sh
environment="$1"
if [ -z "$environment" ]
then
    echo "Usage: destroy.sh <environment> <aws_region>"
    exit 1
fi

aws_region="$2"
if [ -z "$aws_region" ]
then
    echo "Usage: destroy.sh <environment> <aws_region>"
    exit 1
fi

telegram_bot_tokens=`aws ssm get-parameter --region $aws_region --name "${environment}_telegram_bot_tokens" --output text --with-decryption | cut -f7`
if [ -z "$telegram_bot_tokens" ]
then
    printf "paste the telegram_bot_tokens JSON: "
    read telegram_bot_tokens
fi

cd terraform

terraform workspace new $environment
terraform workspace select $environment

terraform apply -destroy \
    --var "aws_region=$aws_region" \
    --var "telegram_bot_tokens=$telegram_bot_tokens"
cd ..