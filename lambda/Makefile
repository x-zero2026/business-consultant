.PHONY: build clean deploy test deps

build:
	sam build

clean:
	rm -rf .aws-sam

deploy:
	sam build && sam deploy

local:
	sam local start-api

test:
	go test ./...

deps:
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/jackc/pgx/v5/pgxpool
	go get github.com/golang-jwt/jwt/v5
	go get github.com/google/uuid
	go mod tidy

# SAM build targets (called by sam build from cmd/*/ directories)
LAMBDA_ROOT := /Users/chilly/go/src/github.com/x-zero2026/business-consultant/lambda

build-ChatFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/chat/main.go

build-SaveReportFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/save-report/main.go

build-GetReportsFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/get-reports/main.go

build-GetReportFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/get-report/main.go

build-DeleteReportFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/delete-report/main.go

build-UpdateReportItemFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/update-report-item/main.go

build-IdentifyProfessionTagsFunction:
	cd $(LAMBDA_ROOT) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/identify-profession-tags/main.go
