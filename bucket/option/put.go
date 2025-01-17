package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// The PutObjectInput type is an adapter to change a parameter in
// s3.PutObjectInput.
type PutObjectInput func(req *s3.PutObjectInput)

// SSEKMSKeyID returns a PutObjectInput that changes a SSE-KMS Key ID.
func SSEKMSKeyID(keyID string) PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.SSEKMSKeyId = aws.String(keyID)
		req.ServerSideEncryption = aws.String("aws:kms")
	}
}

// SSES3 returns a PutObjectInput that uses SSE-S3 (AES256) in S3.
func SSES3() PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.ServerSideEncryption = aws.String("AES256")
	}
}

// ACLPrivate returns a PutObjectInput that set ACL private.
func ACLPrivate() PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.ACL = aws.String(s3.ObjectCannedACLPrivate)
	}
}

// ACLPublicRead returns a PutObjectInput that set ACL public-read.
func ACLPublicRead() PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}
}

// ContentType returns a PutObjectInput that set Content-Type.
func ContentType(ct string) PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.ContentType = aws.String(ct)
	}
}

// ContentLength returns a PutObjectInput that set Content-Length.
func ContentLength(length int64) PutObjectInput {
	return func(req *s3.PutObjectInput) {
		req.ContentLength = aws.Int64(length)
	}
}
