# go-minimum-lambda

This repository is one of minimum templates for AWS Lambda in Go with AWS CloudFormation/SAM.

## Prerequisite

- Go >= 1.9
- [dep](https://github.com/golang/dep) >= 0.4.1
- [mage](https://github.com/magefile/mage) >= 2.1
- aws-cli >= 1.14.40

## Setup

### Configuration file

Prepare `param.cfg` file like below.

```conf
StackName=mizutani-test
CodeS3Bucket=mizutani-test
CodeS3Prefix=functions

LambdaRoleArn=arn:aws:iam::123456789:role/LambdaTestRole
```

- `StackName` : CloudFormation base stack name
- `CodeS3Bucket` : AWS S3 bucket name to save your lambda code (Go binary)
- `CodeS3Prefix` : Prefix of AWS S3 key to save your lambda code
  - *NOTE: `/` is automatically completed at suffix by aws-cli command.
- `LambdaRoleArn` : IAM role ARN of AWS Lambda

### AWS API credential

Please prepare your AWS user and setup credentials to use aws-cli command in your terminal. See [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html) for details.


## Deploy

```
$ dep ensure
$ mage -v deploy
```
