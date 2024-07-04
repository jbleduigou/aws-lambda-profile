# profile: AWS Lambda performance profiling

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/jbleduigou/aws-lambda-profile?status.svg)](https://pkg.go.dev/github.com/jbleduigou/aws-lambda-profile)
![Build Status](https://github.com/jbleduigou/aws-lambda-profile/actions/workflows/go.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/jbleduigou/aws-lambda-profile)](https://goreportcard.com/report/github.com/jbleduigou/aws-lambda-profile)
[![Contributors](https://img.shields.io/github/contributors/jbleduigou/aws-lambda-profile)](https://github.com/jbleduigou/aws-lambda-profile/graphs/contributors)
[![License](https://img.shields.io/github/license/jbleduigou/aws-lambda-profile)](./LICENSE)

An [AWS Lambda Function](https://aws.amazon.com/lambda/) performance profiling tool based on [profile package](https://github.com/pkg/profile) by [Dave Cheney](https://github.com/davecheney).  
The idea is to provide an adapter on top of the existing package to make it easier to use in the context of AWS Lambda.  
The profiling output file is uploaded to an S3 bucket.

## 🚀 Install

```sh
go get github.com/jbleduigou/aws-lambda-profile
```

**Compatibility**: go >= 1.21


## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/jbleduigou/aws-lambda-profile](https://pkg.go.dev/github.com/jbleduigou/aws-lambda-profile)

### Basic example

To enable profiling in your application, you need to add one line in your main function:

```go
import profile "github.com/jbleduigou/aws-lambda-profile"

func main(ctx context.Context) {
    defer profile.Start(profile.S3Bucket("pprof-bucket"), profile.AWSRegion("eu-west-1")).Stop(ctx)
}
```

### CPU profiling

By default, the CPU profiling is enabled.  
You can still explicitly enable it by using the `CPUProfile` option:

```go
import profile "github.com/jbleduigou/aws-lambda-profile"

defer profile.Start(profile.CPUProfile, profile.S3Bucket("pprof-bucket"), profile.AWSRegion("eu-west-1")).Stop(ctx)
```

### Memory profiling

To enable memory profiling, you can use the `MemProfile` option:

```go
import profile "github.com/jbleduigou/aws-lambda-profile"

defer profile.Start(profile.MemProfile, profile.S3Bucket("pprof-bucket"), profile.AWSRegion("eu-west-1")).Stop(ctx)
```



## 🤝 Contributing

Feel free to contribute to this project, either my opening issues or submitting pull requests.  
Don't hesitate to contact me, by sending me a PM on [LinkedIn](www.linkedin.com/in/jbleduigou).

## 👤 Contributors

The only contributor so far is me, Jean-Baptiste Le Duigou.
Feel free to check my blog to about my other projects: http://www.jbleduigou.com

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/jbleduigou?style=for-the-badge)](https://github.com/sponsors/jbleduigou)

## 📝 License

This project is [BSD](./LICENSE) licensed.
