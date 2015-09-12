package fb

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ranveerkunal/admin"
)

const (
	EndPoint = "https://graph.facebook.com/v2.0"
)

var (
	ClientID     = flag.String("fbClientID", os.Getenv("FB_CLIENT_ID"), "")
	ClientSecret = flag.String("fbClientSecret", os.Getenv("FB_CLIENT_SECRET"), "")
	VerifyToken  = flag.String("verifyToken", os.Getenv("VERIFY_TOKEN"), "")
	SubURL       = fmt.Sprintf("%s/%s/subscriptions", EndPoint, *ClientID)
)

type SignedRequest struct {
	Code      string `json:"code"`
	Algorithm string `json:"algorithm"`
	IssuedAt  uint64 `json:"issued_at"`
	UserId    string `json:"user_id"`
}

type Stub struct {
	verifyToken string
	accessToken []string
}

func NewStub() *Stub {
	return &Stub{
		verifyToken: *VerifyToken,
	}
}

func (fbs *Stub) OAuth(s *admin.Status) error {
	v := url.Values{}
	v.Add("client_id", *ClientID)
	v.Add("client_secret", *ClientSecret)
	v.Add("grant_type", "client_credentials")
	res, err := http.Get(fmt.Sprintf("%s/oauth/access_token?%s", EndPoint, v.Encode()))
	ok := err == nil && res.StatusCode == http.StatusOK
	s.Log["FB Access Token"] = admin.NewLog(ok, res)
	if !ok {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	fbs.accessToken = strings.Split(string(body), "=")
	res.Body.Close()
	return err
}
