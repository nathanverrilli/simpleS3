package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"simpleS3/misc"
	"strings"
	"time"
)

// myFatal() -- flushes and closes log
// file prior to exiting with error status
func myFatal() {
	_ = xLogBuffer.Flush()
	_ = xLogFile.Close()
	os.Exit(1)
}

// exitError() -- AWS claims this function is needed
// but I think AWS is mistaken
/**************************
func exitErrorf(msg string, args ...interface{}) {
	xLog.Printf(msg+"\n", args...)
	myFatal()
}
*************************/

const AWS_REGION = "us-west-000"
const AWS_ENDPOINT = "s3.us-west-000.backblazeb2.com"

func main() {

	initLog()
	defer misc.DeferError(xLogFile.Close)
	defer misc.DeferError(xLogBuffer.Flush)
	initFlags()

	var awsConfig aws.Config
	awsConfig.Region = aws.String(AWS_REGION)
	awsConfig.Endpoint = aws.String(AWS_ENDPOINT)

	// use the BACKBLAZE profile in the ~/.aws/credentials file which we're not using :-)
	_ = os.Setenv("AWS_PROFILE", "BACKBLAZE")

	// from the working directory, which is where the program runs in the ide
	// change this as appropriate ...
	// _ = os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./.aws/credentials")

	// where I keep the real credentials relative to the project source
	_ = os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "../.aws/credentials")

	// might not need this, since we're specifying it in the config?
	_ = os.Setenv("AWS_DEFAULT_REGION", AWS_REGION)

	sess, err := session.NewSession(&awsConfig)
	if nil != err {
		xLog.Fatalf("could not create session for region %s because %s",
			AWS_ENDPOINT, err.Error())
	}

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if nil != err {
		xLog.Printf("could not list in region %s endpoint %s because %s\n",
			AWS_REGION, AWS_ENDPOINT, err.Error())
		myFatal()
	}

	var sb strings.Builder
	sb.WriteString("\tfound buckets: \n")
	for ix, bucket := range result.Buckets {
		sb.WriteString(
			fmt.Sprintf("%3d:\t%s created on %s\n",
				(ix + 1),
				aws.StringValue(bucket.Name),
				aws.TimeValue(bucket.CreationDate)))
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s", sb.String())

	// list the objects in the bucket
	var lbo s3.ListObjectsV2Input
	lbo.Bucket = aws.String(*result.Buckets[0].Name)
	resp, err := svc.ListObjectsV2(&lbo)
	if nil != err {
		xLog.Printf("could not list object in bucket %s because %s\n",
			result.Buckets[0].Name, err.Error())
		myFatal()
	}

	sb.Reset()
	sb.WriteString("list of objects in bucket ")
	sb.WriteString(*result.Buckets[0].Name)
	fmt.Println(sb.String())
	for ix, item := range resp.Contents {
		sb.Reset()
		sb.WriteString(fmt.Sprintf("\n****** ITEM #:\t%d\n", ix+1))
		sb.WriteString(fmt.Sprintf("         name:\t%s\n", *item.Key))
		sb.WriteString(fmt.Sprintf("last modified:\t%s\n", item.LastModified.UTC().Format(time.RFC850)))
		sb.WriteString(fmt.Sprintf("         size:\t%d\n", *item.Size))
		sb.WriteString(fmt.Sprintf("storage class:\t%s\n\n", *item.StorageClass))
		fmt.Print(sb.String())
	}

	// download (a) file
	newFileName := "dwnld_" + *resp.Contents[0].Key
	downloadFile, err := os.OpenFile(newFileName,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Printf("Error opening download file %s because %s",
			newFileName, err.Error())
	}
	defer misc.DeferError(downloadFile.Close)
	// do not buffer this; the downloader gets
	// chunks and may get them out of order.

	var goi s3.GetObjectInput
	goi.Bucket = aws.String(*result.Buckets[0].Name)
	goi.Key = aws.String(*resp.Contents[0].Key)
	dl := s3manager.NewDownloader(sess)
	byteCount, err := dl.Download(downloadFile, &goi)

	if nil != err {
		xLog.Printf("unable to download item %s in bucket %s because %s",
			*result.Buckets[0].Name, *resp.Contents[0].Key, err.Error())
		myFatal()
	}

	if FlagDebug || FlagVerbose {
		xLog.Printf("downloaded %s as %s, in %d bytes",
			*resp.Contents[0].Key, newFileName, byteCount)
	}

}
