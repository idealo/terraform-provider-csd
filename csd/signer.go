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

func signer(request *http.Request, accessKeyId string, secretAccessKey string, sessionToken string) *Authorizer {
	currentTime := time.Now().UTC()

	const (
		awsService         = "execute-api"
		requestQuery       = ""
		requestContentType = "application/json"
		awsRegion          = "eu-central-1"
		signedHeaders      = "content-type;host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
		dateFmt            = "20060102"         // this is just the format, not a specific time. I know, Go is stupid.
		timeFmt            = "20060102T150405Z" // Same here, look at: https://pkg.go.dev/time
	)

	getBody, _ := request.GetBody()
	body, _ := io.ReadAll(getBody)

	requestPayloadHash := hashSHA256(body)

	var requestHeaders = fmt.Sprintf("content-type:%v\nhost:%v\nx-amz-content-sha256:%x\nx-amz-date:%s\nx-amz-security-token:%s", requestContentType, request.Host, requestPayloadHash, currentTime.Format(timeFmt), sessionToken)
	var canonicalRequest = fmt.Sprintf("%v\n%v\n%v\n%v\n\n%v\n%x", request.Method, request.URL.Path, requestQuery, requestHeaders, signedHeaders, requestPayloadHash)

	var stringToSign = fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/%s/%s/aws4_request\n%x", currentTime.Format(timeFmt), currentTime.Format(dateFmt), awsRegion, awsService, hashSHA256([]byte(canonicalRequest)))

	kDate := hmacSHA256([]byte("AWS4"+secretAccessKey), []byte(currentTime.Format(dateFmt)))
	kRegion := hmacSHA256(kDate, []byte(awsRegion))
	kService := hmacSHA256(kRegion, []byte(awsService))
	signingKey := hmacSHA256(kService, []byte("aws4_request"))
	requestSignature := hmacSHA256(signingKey, []byte(stringToSign))

	var signature = fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=%s;x-amz-date, Signature=%x", accessKeyId, currentTime.Format(dateFmt), awsRegion, awsService, signedHeaders, requestSignature)

	return &Authorizer{
		authorizationHeaders: signature,
		date:                 currentTime.Format(timeFmt),
		payloadHash:          requestPayloadHash,
	}
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
