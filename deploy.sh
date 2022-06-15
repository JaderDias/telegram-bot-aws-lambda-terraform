#!/bin/bash

echo -e "\n+++++ Starting deployment +++++\n"

NL_DICT="golang/upload/nl.csv"
SH_DICT="golang/upload/sh.csv"

if [ ! -f "$NL_DICT" ] ||  [ ! -f "$SH_DICT" ]; then
    DUMP_XML_BZ2="enwiktionary-latest-pages-articles-multistream.xml.bz2"
    if [ ! -f "../$DUMP_XML_BZ2" ]; then
        wget "https://dumps.wikimedia.org/enwiktionary/latest/$DUMP_XML_BZ2"
        mv $DUMP_XML_BZ2 ../
    fi
    if [ ! -f "$NL_DICT" ]; then
        python3 python/parser/filter_wiktionary.py Dutch A-ZÁÉÍÓÚÀÈËÏÖÜĲ ../$DUMP_XML_BZ2 | tee $NL_DICT
    fi
    if [ ! -f "$SH_DICT" ]; then
        python3 python/parser/filter_wiktionary.py Serbo-Croatian A-ZÁČĆĐÍĽŇÔŠŤÚÝŽ ../$DUMP_XML_BZ2 | tee $SH_DICT
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

telegram_bot_tokens=`aws ssm get-parameter --name telegram_bot_tokens --output text --with-decryption | cut -f7`
if [ -z "$telegram_bot_tokens" ]
then
    printf "paste the telegram_bot_tokens JSON: "
    read telegram_bot_tokens
fi

terraform apply --auto-approve --var "telegram_bot_tokens=$telegram_bot_tokens" 

echo -e "\n+++++ Deployment done +++++\n"