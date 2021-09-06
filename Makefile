GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

# .PHONY default:
# default: build

.PHONY build:
build: # Build lambda function go binaries
	@go build -o ./lambda/main ./lambda/main.go

.PHONY deploy:
deploy: build
deploy: # Deploy to CloudFormation
	cdk deploy

.PHONY destroy:
destroy: build
destroy: # Tear down resources
	cdk destroy

.PHONY synth:
synth: # Synth
	cdk synth