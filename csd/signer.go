package csd

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Authorizer struct {
	authorizationHeaders string
	date                 string
	payloadHash          []byte
}

func signer(c *AuthInfo, r *http.Request) *Authorizer {
	currentTime := time.Now().UTC()

	const (
		awsService         = "execute-api"
		requestQuery       = ""
		requestContentType = "application/json"
		awsRegion          = "us-east-1"
		signedHeaders      = "content-type;host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
		dateFmt            = "20060102"
		timeFmt            = "20060102T150405Z"
	)

	getBody, _ := r.GetBody()
	body, _ := io.ReadAll(getBody)

	requestPayloadhash := hashSHA256([]byte(string(body)))

	var requestHeaders = fmt.Sprintf("content-type:%v\nhost:%v\nx-amz-content-sha256:%x\nx-amz-date:%s\nx-amz-security-token:%s", requestContentType, r.Host, requestPayloadhash, currentTime.Format(timeFmt), c.SessionToken)
	var cononicalRequest = fmt.Sprintf("%v\n%v\n%v\n%v\n\n%v\n%x", r.Method, r.URL.Path, requestQuery, requestHeaders, signedHeaders, requestPayloadhash)

	var stringToSign = fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/%s/%s/aws4_request\n%x", currentTime.Format(timeFmt), currentTime.Format(dateFmt), awsRegion, awsService, hashSHA256([]byte(cononicalRequest)))

	kDate := hmacSHA256([]byte("AWS4"+c.SecretAccessKey), []byte(currentTime.Format(dateFmt)))
	kRegion := hmacSHA256(kDate, []byte(awsRegion))
	kService := hmacSHA256(kRegion, []byte(awsService))
	signingKey := hmacSHA256(kService, []byte("aws4_request"))
	requestSignature := hmacSHA256(signingKey, []byte(stringToSign))

	var signature = fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=%s;x-amz-date, Signature=%x", c.AccessKeyId, currentTime.Format(dateFmt), awsRegion, awsService, signedHeaders, requestSignature)

	p := Authorizer{authorizationHeaders: signature, date: currentTime.Format(timeFmt), payloadHash: requestPayloadhash}

	return &p
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func hmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
