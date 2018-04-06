package main

import "flag"
import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"log"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"time"
	"github.com/google/go-querystring/query"
	"net/http"
	"io"
	"os"
	"encoding/hex"
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

const serviceName = "rds"
var region = "ap-southeast-1"
var host = "rds.ap-southeast-1.amazonaws.com"
var accessKeyId = ""
var secretKey = ""
var dateStamp = time.Now().UTC().Format("20060102")

type QueryOptions struct {
	Algorithm   	string 	`url:"X-Amz-Algorithm"`
	Credential 		string  `url:"X-Amz-Credential"`
	Expires    		int    	`url:"X-Amz-Expires"`
	Date    		string  `url:"X-Amz-Date"`
	SignedHeaders   string  `url:"X-Amz-SignedHeaders"`
}

func downloadFile(path string, url string) (err error) {
	println("Downloading to " + path)
	out, err := os.Create(path)
	if err != nil  {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}
	return nil
}

func sign(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func getSignatureKey() []byte {
	signatureKey := sign([]byte("AWS4" + secretKey), []byte(dateStamp))
	signatureKey = sign(signatureKey, []byte(region))
	signatureKey = sign(signatureKey, []byte(serviceName))
	signatureKey = sign(signatureKey, []byte("aws4_request"))
	return signatureKey
}

func buildSignedUrl(instanceIdentifier string, logFile string) string {
	requestDate := time.Now().UTC().Format("20060102T150405Z")
	canonicalUri := fmt.Sprintf("/v13/downloadCompleteLogFile/%s/%s", instanceIdentifier, logFile)
	canonicalHeaders := fmt.Sprintf("host:%s\n", host)
	signedHeaders := "host"
	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, serviceName)
	queryParams := QueryOptions{
		Algorithm: algorithm,
		Credential: accessKeyId + "/" + credentialScope,
		Date: requestDate,
		Expires: 30,
		SignedHeaders: signedHeaders,
	}
	queryParamValues, _ := query.Values(queryParams)
	canonicalQueryString := queryParamValues.Encode()

	payloadHasher := sha256.New()
	payloadHasher.Write([]byte(""))
	payloadHash := fmt.Sprintf("%x", payloadHasher.Sum(nil))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", "GET", canonicalUri, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)
	requestHasher := sha256.New()
	requestHasher.Write([]byte(canonicalRequest))

	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%x",  algorithm, requestDate, credentialScope, requestHasher.Sum(nil))
	signingKey := getSignatureKey()
	signature := hex.EncodeToString(sign(signingKey, []byte(stringToSign)))
	canonicalQueryString += "&X-Amz-Signature=" + signature
	return "https://" + host + canonicalUri + "?" + canonicalQueryString
}

func main()  {
	instanceIdentifier := flag.String("instance", "foo-master", "Your instance identifier")
	logFile := flag.String("log-file", "all", "Log file to download")
	destination := flag.String("destination", "./logs", "Destination to store the log files")
	flag.Parse()
	currentSession := session.Must(session.NewSession())
	region = *currentSession.Config.Region
	currentCredentials := currentSession.Config.Credentials
	credentials, err := currentCredentials.Get()
	if err != nil {
		log.Fatal(err)
		return
	}
	accessKeyId = credentials.AccessKeyID
	secretKey = credentials.SecretAccessKey
	host = fmt.Sprintf("rds.%s.amazonaws.com", region)

	os.MkdirAll(*destination, os.ModePerm)


	if *logFile == "all" {
		rdsService := rds.New(currentSession)
		input := &rds.DescribeDBLogFilesInput{
			DBInstanceIdentifier: aws.String(*instanceIdentifier),
		}

		result, err := rdsService.DescribeDBLogFiles(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case rds.ErrCodeDBInstanceNotFoundFault:
					fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		for _, logFile := range result.DescribeDBLogFiles {
			fileUrl := buildSignedUrl(*instanceIdentifier, *logFile.LogFileName)
			_, fileName := filepath.Split(*logFile.LogFileName)
			downloadFile(filepath.Join(*destination, fileName) , fileUrl)
		}

	} else {
		fileUrl := buildSignedUrl(*instanceIdentifier, *logFile)
		_, fileName := filepath.Split(*logFile)
		downloadFile(filepath.Join(*destination, fileName) , fileUrl)
	}
}