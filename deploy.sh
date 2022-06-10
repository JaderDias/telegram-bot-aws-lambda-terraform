#!/bin/bash

telegram_bot_token="$1"

echo -e "\n+++++ Starting deployment +++++\n"

rm -rf ./bin

echo "+++++ build go packages +++++"

cd golang/reply
go get
go test ./...
env GOOS=linux GOARCH=amd64 go build -o ../../bin/reply

echo "+++++ apply terraform +++++"
cd ../../terraform
if [ ! -f 'terraform.tfstate' ]; then
  terraform init
fi

terraform apply --var "telegram_bot_token=$telegram_bot_token"

echo -e "\n+++++ Deployment done +++++\n"