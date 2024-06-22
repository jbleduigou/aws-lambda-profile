package profile

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	pkgprofile "github.com/pkg/profile"
)

// Profile is an adapter for github.com/pkg/profile
type Profile struct {
	// bucket is the name of the S3 bucket to upload the profile to.
	bucket string

	// region is the AWS region of the S3 bucket.
	region string

	// localPath is the path to the local file where the profile is stored.
	localPath string

	// p is the actual profile object from github.com/pkg/profile
	p interface{ Stop() }
}

// Quiet suppresses informational messages during profiling.
func Quiet(*Profile) func(pkg *pkgprofile.Profile) {
	return pkgprofile.Quiet
}

// CPUProfile enables CPU profiling.
// It disables any previous profiling settings.
func CPUProfile(p *Profile) func(pkg *pkgprofile.Profile) {
	p.localPath = "/tmp/cpu.pprof"
	return pkgprofile.CPUProfile
}

// MemProfile enables memory profiling.
// It disables any previous profiling settings.
func MemProfile(p *Profile) func(pkg *pkgprofile.Profile) {
	p.localPath = "/tmp/mem.pprof"
	return pkgprofile.MemProfile
}

// S3Bucket sets the name of the S3 bucket to upload the profile to.
func S3Bucket(bucket string) func(*Profile) func(pkg *pkgprofile.Profile) {
	return func(p *Profile) func(pkg *pkgprofile.Profile) {
		p.bucket = bucket
		return nil
	}
}

// AWSRegion sets the AWS region of the S3 bucket.
func AWSRegion(region string) func(*Profile) func(pkg *pkgprofile.Profile) {
	return func(p *Profile) func(pkg *pkgprofile.Profile) {
		p.region = region
		return nil
	}
}

// Stop stops the profile and flushes any unwritten data.
func (p *Profile) Stop() {
	p.p.Stop()
	p.uploadToS3()
}

// uploadToS3 uploads the profile to S3 if the bucket name is provided.
func (p *Profile) uploadToS3() {
	if p.bucket == "" {
		log.Fatalf("bucket name is not provided, skipping upload to S3")
		return
	}

	file, err := os.Open(p.localPath)
	if err != nil {
		log.Fatalf("unable to open file %s, %v", p.localPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("unable to close file, %v", err)
		}
	}(file)

	s3Path := p.generateS3Path()

	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(p.region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(s3Path),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("unable to upload file, %v", err)
	}
}

// generateS3Path generates the S3 path based on the lambda context and the current local path.
func (p *Profile) generateS3Path() string {
	// retrieve lambda context
	lc, found := lambdacontext.FromContext(context.Background())
	if !found {
		return strings.Replace(p.localPath, "/tmp/", "unknown/", 1)
	}
	return strings.Replace(p.localPath, "/tmp/", fmt.Sprintf("%s/%s/", lc.InvokedFunctionArn, lc.AwsRequestID), 1)
}

// Start starts a new profiling session.
// The caller should call the Stop method on the value returned
// to cleanly stop profiling.
func Start(options ...func(*Profile) func(pkg *pkgprofile.Profile)) interface {
	Stop()
} {

	var pkgOptions []func(*pkgprofile.Profile)
	var profile Profile
	for _, option := range options {
		o := option(&profile)
		if o != nil {
			pkgOptions = append(pkgOptions, o)
		}
	}
	pkgOptions = append(pkgOptions, pkgprofile.ProfilePath("/tmp"))

	if profile.region == "" {
		profile.region = "us-east-1"
	}

	if profile.localPath == "" {
		profile.localPath = "/tmp/cpu.pprof"
	}

	profile.p = pkgprofile.Start(pkgOptions...)
	return &profile
}
