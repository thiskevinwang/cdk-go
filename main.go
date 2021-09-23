package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
type TokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

func main() {
	console_url := "https://console.aws.amazon.com/"
	sign_in_url := "https://signin.aws.amazon.com/federation"

	args := os.Args[1:]
	name := "stranger_danger"
	if len(args) >= 1 {
		name = args[0]
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)
	creds, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
		Name: aws.String(name),
		PolicyArns: []types.PolicyDescriptorType{
			types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess"),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("AccessKeyId:", aws.ToString(creds.Credentials.AccessKeyId))
	// fmt.Println("SecretAccessKey:", aws.ToString(creds.Credentials.SecretAccessKey))
	// fmt.Println("SessionToken:", aws.ToString(creds.Credentials.SessionToken))

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
	// fmt.Println(string(b))

	params := url.Values{}
	params.Add("Action", "getSigninToken")
	params.Add("SessionDuration", "43200")
	params.Add("Session", string(b))

	baseUrl.RawQuery = params.Encode()

	uri := baseUrl.String()
	// fmt.Println(uri)

	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		params := url.Values{}
		params.Add("Action", "login")
		params.Add("Issuer", "arbitrary_issuer")
		params.Add("Destination", console_url)

		data := TokenResponse{}
		json.Unmarshal(bodyBytes, &data)
		// fmt.Printf("SigninToken: %s", data.SigninToken)
		params.Add("SigninToken", data.SigninToken)

		baseUri, err := url.Parse(sign_in_url)
		if err != nil {
			fmt.Println("Malformed URL: ", err.Error())
			return
		}
		baseUri.RawQuery = params.Encode()

		console_uri := baseUri.String()
		// fmt.Println(console_uri)
		fmt.Println("Attempting to open AWS console...")
		openbrowser(console_uri)
	} else {
		fmt.Println("Non 200 status code")
	}
	if err != nil {
		fmt.Println(err)
	}

}

// %7B = {
// %22 = "
// %3A = :

// Example @ https://aws.amazon.com/blogs/security/enable-your-federated-users-to-work-in-the-aws-management-console-for-up-to-12-hours/

// https://signin.aws.amazon.com/federation
// ?Action=getSigninToken
// &Session=%7B%22sessionId%22%3A%22ASIAEXAMPLEMD
// LUUAEYQ%22%2C%22sessionKey%22%3A%22tpSl9thxr2PkEXAMPLETAnVLVGdwC5zXtGDr
// %2FqWi%22%2C%22sessionToken%22%3A%22AQoDYXdz%EXAMPLE
//&SessionDuration=43200
