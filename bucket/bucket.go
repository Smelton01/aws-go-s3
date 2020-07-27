package bucket

import (
	"io"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/nabeken/aws-go-s3/bucket/option"
)

// A Bucket is an S3 bucket which holds properties such as bucket name and SSE things for S3 Bucket.
type Bucket struct {
	S3   s3iface.S3API
	Name *string
}

// New returns Bucket instance with bucket name name.
func New(s s3iface.S3API, name string) *Bucket {
	return &Bucket{
		S3:   s,
		Name: aws.String(name),
	}
}

// GetObject returns the s3.GetObjectOutput.
func (b *Bucket) GetObject(key string, opts ...option.GetObjectInput) (*s3.GetObjectOutput, error) {
	req := &s3.GetObjectInput{
		Bucket: b.Name,
		Key:    aws.String(key),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.GetObject(req)
}

// GetObjectReader returns a reader assosiated with body. A caller of this MUST close the reader when it finishes reading.
func (b *Bucket) GetObjectReader(key string, opts ...option.GetObjectInput) (io.ReadCloser, error) {
	resp, err := b.GetObject(key, opts...)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// GetObjectRequest generates a "aws/request.Request" representing the client's request for the GetObject operation.
func (b *Bucket) GetObjectRequest(key string, opts ...option.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	req := &s3.GetObjectInput{
		Bucket: b.Name,
		Key:    aws.String(key),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.GetObjectRequest(req)
}

// HeadObject retrieves an object metadata for key.
func (b *Bucket) HeadObject(key string, opts ...option.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	req := &s3.HeadObjectInput{
		Bucket: b.Name,
		Key:    aws.String(key),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.HeadObject(req)
}

// ExistsObject returns true if key does not exist on bucket.
func (b *Bucket) ExistsObject(key string, opts ...option.HeadObjectInput) (bool, error) {
	_, err := b.HeadObject(key, opts...)
	if err == nil {
		return true, nil
	}

	if s3err, ok := err.(awserr.RequestFailure); ok && s3err.StatusCode() == http.StatusNotFound {
		// actually key does not exist
		return false, nil
	}

	// in some error situation
	return false, err
}

// PutObject puts an object with reading data from reader.
func (b *Bucket) PutObject(key string, rs io.ReadSeeker, opts ...option.PutObjectInput) (*s3.PutObjectOutput, error) {
	req := &s3.PutObjectInput{
		Bucket: b.Name,
		Key:    aws.String(key),
		Body:   rs,
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.PutObject(req)
}

// DeleteObject deletes an object for key.
func (b *Bucket) DeleteObject(key string) (*s3.DeleteObjectOutput, error) {
	req := &s3.DeleteObjectInput{
		Bucket: b.Name,
		Key:    aws.String(key),
	}

	return b.S3.DeleteObject(req)
}

// DeleteObjects deletes each object for the given identifiers.
// A maximum of 1000 objects can be deleted at a time with this method.
func (b *Bucket) DeleteObjects(identifiers []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error) {
	req := &s3.DeleteObjectsInput{
		Bucket: b.Name,
		Delete: &s3.Delete{
			Objects: identifiers,
		},
	}

	return b.S3.DeleteObjects(req)
}

// ListObjects lists objects that has prefix.
func (b *Bucket) ListObjects(prefix string, opts ...option.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	req := &s3.ListObjectsInput{
		Bucket: b.Name,
		Prefix: aws.String(prefix),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.ListObjects(req)
}

// ListObjectsV2PagesWithContext will page through objects with the given prefix.
func (b *Bucket) ListObjectsV2PagesWithContext(
	ctx aws.Context,
	prefix string,
	pageFunc func(*s3.ListObjectsV2Output, bool) bool,
	opts ...option.ListObjectsV2Input,
) error {
	req := &s3.ListObjectsV2Input{
		Bucket: b.Name,
		Prefix: aws.String(prefix),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.ListObjectsV2PagesWithContext(ctx, req, pageFunc)
}

// ListObjectVersionsPagesWithContext will page through all versions of all objects with the given prefix.
func (b *Bucket) ListObjectVersionsPagesWithContext(
	ctx aws.Context,
	prefix string,
	pageFunc func(*s3.ListObjectVersionsOutput, bool) bool,
	opts ...option.ListObjectVersionsInput,
) error {
	req := &s3.ListObjectVersionsInput{
		Bucket: b.Name,
		Prefix: aws.String(prefix),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.ListObjectVersionsPagesWithContext(ctx, req, pageFunc)
}

// CopyObject copies an object within the bucket.
func (b *Bucket) CopyObject(dest, src string, opts ...option.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	req := &s3.CopyObjectInput{
		Bucket:     b.Name,
		Key:        aws.String(dest),
		CopySource: aws.String(aws.StringValue(b.Name) + "/" + url.QueryEscape(src)),
	}

	for _, f := range opts {
		f(req)
	}

	return b.S3.CopyObject(req)
}
