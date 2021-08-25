<html>
<div align="center">
<h1>Go Lambda</h1>
<h3>Create a Golang AWS lambda function, using the AWS CDK</h3>

</div>
</html>

---

## Quick start

Build the lambda function binary
- `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./lambda/main ./lambda/main.go`

Deploy infrastructure
- `cdk deploy`

Clean up infrastructure
- `cdk destroy`

---

This is a blank project for Go development with CDK.

**NOTICE**: Go support is still in Developer Preview. This implies that APIs may
change while we address early feedback from the community. We would love to hear
about your experience through GitHub issues.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `cdk destroy`     clean up any provisioned infrastructure
 * `go test`         run unit tests
