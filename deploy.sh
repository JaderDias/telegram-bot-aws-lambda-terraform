#!/bin/bash

echo -e "\n+++++ Starting deployment +++++\n"

SH_DICT="sh.csv"

if [ ! -f "$SH_DICT" ]; then
    DUMP_XML_BZ2="enwiktionary-latest-pages-articles-multistream.xml.bz2"
    if [ ! -f "../$DUMP_XML_BZ2" ]; then
        wget "https://dumps.wikimedia.org/enwiktionary/latest/$DUMP_XML_BZ2"
        mv $DUMP_XML_BZ2 ../
    fi
    if [ ! -f "$SH_DICT" ]; then
        python3 python/parser/filter_wiktionary.py Serbo-Croatian A-ZÁČĆĐÍĽŇÔŠŤÚÝŽ ../$DUMP_XML_BZ2 | tee $SH_DICT
        source upload.sh "sh"
    fi
fi

rm -rf ./bin
mkdir ./bin
mkdir ./bin/reply

echo "+++++ build go packages +++++"

cd golang/reply
go get
go test ./...
env GOOS=linux GOARCH=amd64 go build -o ../../bin/reply/reply

echo "+++++ apply terraform +++++"
cd ../../terraform
if [ ! -f 'terraform.tfstate' ]; then
  terraform init
fi

telegram_bot_token=`aws ssm get-parameter --name telegram_bot_token --output text --with-decryption | cut -f7`
if [ -z "$telegram_bot_token" ]
then
    printf "paste the telegram bot token for the SH language: "
    read telegram_bot_token
fi

terraform apply --var "telegram_bot_token=$telegram_bot_token"

echo -e "\n+++++ Deployment done +++++\n"