module example.com/main

go 1.18

replace example.com/telegram => ../telegram

require (
	example.com/telegram v0.0.0-00010101000000-000000000000
	gopkg.in/telegram-bot-api.v4 v4.6.4
)

require (
	github.com/aws/aws-sdk-go-v2 v1.16.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.15.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.27.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.7 // indirect
	github.com/aws/smithy-go v1.11.3 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
