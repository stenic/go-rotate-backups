package drivers

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
)

func init() {
	AddDriver("s3", &S3Driver{})
}

type S3Driver struct {
	BaseDriver

	bucket string
}

func (d *S3Driver) getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func (d *S3Driver) Init() error {
	d.bucket = os.Getenv("GRB_S3_BUCKET")

	if d.bucket == "" {
		return fmt.Errorf("you need to set 'GRB_S3_BUCKET' for a target bucket")
	}

	svc := sts.New(d.getSession())
	identity, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	logrus.Debugf("Connected to s3 as %s", identity.String())

	return err
}

func (d *S3Driver) ListDirs(path string) ([]string, error) {
	res := []string{}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	items, err := d.listRaw(path)
	if err != nil {
		return res, err
	}

	for _, i := range items {
		res = append(res, strings.Split(strings.Replace(i, path, "", 1), "/")[0])
	}

	return res, err
}

func (d *S3Driver) listRaw(path string) ([]string, error) {
	res := []string{}

	svc := s3.New(d.getSession())

	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket:  aws.String(d.bucket),
		Prefix:  aws.String(path),
		MaxKeys: aws.Int64(20),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, i := range page.Contents {
			res = append(res, *i.Key)
		}
		return true
	})

	logrus.Tracef("Listing %s:%s -> %v", d.bucket, path, res)

	return res, err
}

func (d *S3Driver) Mkdir(path string) error {
	// Not needed in s3
	return nil
}

func (d *S3Driver) Delete(src string) error {

	items, err := d.listRaw(src)
	if err != nil {
		return err
	}

	svc := s3.New(d.getSession())
	for _, item := range items {
		logrus.Tracef("Deleting %s:%s", d.bucket, item)
		if _, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(item),
		}); err != nil {
			return err
		}
	}

	return err
}

func (d *S3Driver) Copy(src, dst string) (int64, error) {
	uploader := s3manager.NewUploader(d.getSession())

	f, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %q, %v", src, err)
	}

	logrus.Tracef("Uploading %s:%s", d.bucket, src)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(dst),
		Body:   f,
	})

	info, _ := f.Stat()

	return info.Size(), err
}
