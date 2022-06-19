#!/bin/bash
environment="$1"
if [ -z "$environment" ]
then
    echo "Usage: deploy.sh <environment> <aws_region>"
    exit 1
fi

aws_region="$2"
if [ -z "$aws_region" ]
then
    echo "Usage: deploy.sh <environment> <aws_region>"
    exit 1
fi

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
mkdir ./bin/upload
cp *.csv ./bin/upload/

echo "+++++ build go packages +++++"

cd golang/reply
go get
env GOOS=linux GOARCH=amd64 go build -o ../../bin/reply/reply
if [ $? -ne 0 ]
then
    echo "build reply packages failed"
    exit 1
fi

cd ../upload
go get
env GOOS=linux GOARCH=amd64 go build -o ../../bin/upload/upload
if [ $? -ne 0 ]
then
    echo "build upload packages failed"
    exit 1
fi

echo "+++++ apply terraform +++++"
cd ../../terraform
terraform init
if [ $? -ne 0 ]
then
    echo "terraform init failed"
    exit 1
fi

terraform workspace new $environment
terraform workspace select $environment

telegram_bot_tokens=`aws ssm get-parameter --region $aws_region --name "${environment}_telegram_bot_tokens" --output text --with-decryption | cut -f7`
if [ -z "$telegram_bot_tokens" ]
then
    printf "paste the telegram_bot_tokens JSON: "
    read telegram_bot_tokens
fi

terraform apply --auto-approve \
    --var "aws_region=$aws_region" \
    --var "telegram_bot_tokens=$telegram_bot_tokens"

echo -e "\n+++++ Deployment done +++++\n"