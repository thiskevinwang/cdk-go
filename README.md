CHANGE FOR RandoApp9000 #3!
<html>
<div align="center">
<h1>Go CDK</h1>
<h3>AWS Infrastructure as Code with <code>go</code></h3>

</div>
</html>

---

## Quick start

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

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `cdk destroy` clean up any provisioned infrastructure
- `go test` run unit tests

---

# Learnings

Looking at [/aws/aws-cdk-go/awscdk@v1.121](https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk), there are two Lambda packages to choose from:

- 👎 [`"github.com/aws/aws-cdk-go/awscdk/awslambda"`][awslambda]
- 👍 [`"github.com/aws/aws-cdk-go/awscdk/awslambdago"`][awslambdago]

👎 If you use `awslambda`, you'll need to build your `go` lambda function binary manually first

- `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./lambda/main ./lambda/main.go`

👍 If you use `awslambdago`, you do not need to manually build any binary

[awslambda]: https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk@v1.121.0-devpreview/awslambda
[awslambdago]: https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk@v1.121.0-devpreview/awslambdago

## Add AWS SDK code to Lambda fn

https://github.com/aws/aws-sdk-go-v2

```bash
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/dynamodb
```
