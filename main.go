package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

type Session struct {
	Id    string `json:"sessionId"`
	Key   string `json:"sessionKey"`
	Token string `json:"sessionToken"`
}

func main() {
	console_url := "https://console.aws.amazon.com/"
	sign_in_url := "https://signin.aws.amazon.com/federation"

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)
	creds, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
		Name: aws.String("foobar"),
		PolicyArns: []types.PolicyDescriptorType{
			types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess"),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("AccessKeyId:", aws.ToString(creds.Credentials.AccessKeyId))
	fmt.Println("SecretAccessKey:", aws.ToString(creds.Credentials.SecretAccessKey))
	fmt.Println("SessionToken:", aws.ToString(creds.Credentials.SessionToken))

	baseUrl, err := url.Parse(sign_in_url)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
	}

	session := Session{
		Id:    aws.ToString(creds.Credentials.AccessKeyId),
		Key:   aws.ToString(creds.Credentials.SecretAccessKey),
		Token: aws.ToString(creds.Credentials.SessionToken),
	}
	b, err := json.Marshal(session)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	params := url.Values{}
	params.Add("Action", "getSigninToken")
	// params.Add("SessionDuration", "43200")
	params.Add("SessionType", "json")
	params.Add("Session", url.QueryEscape(string(b)))

	baseUrl.RawQuery = params.Encode()

	uri := baseUrl.String()
	fmt.Println(uri)

	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)

		_params := url.Values{}
		_params.Add("Action", "login")
		_params.Add("Issuer", "foobar")
		_params.Add("Destination", console_url)

		uri := fmt.Sprintf("", sign_in_url)
		openbrowser(uri)
	} else {
		fmt.Println("Non 200 status code")
	}
	if err != nil {
		fmt.Println(err)
	}

}
