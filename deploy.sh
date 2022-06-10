#!/bin/bash

telegram_bot_token="$1"

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
    fi
fi

rm -rf ./bin
mkdir ./bin
mkdir ./bin/reply
cp *.csv ./bin/reply/

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

terraform apply --var "telegram_bot_token=$telegram_bot_token"

echo -e "\n+++++ Deployment done +++++\n"