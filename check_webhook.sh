#!/bin/bash

telegram_bot_token="$1"

curl "https://api.telegram.org/bot${telegram_bot_token}/getWebhookInfo"