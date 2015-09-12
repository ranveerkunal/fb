package fb

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-martini/martini"
	"github.com/ranveerkunal/admin"
	"github.com/ranveerkunal/weblogger"
)

func (fbs *Stub) Subscribe(s *admin.Status) {
	v := url.Values{}
	v.Add("object", "user")
	v.Add("fields", "events")
	v.Add("verify_token", fbs.verifyToken)
	v.Add("callback_url", s.AppDomain+"fbsub/cb")
	v.Add(fbs.accessToken[0], fbs.accessToken[1])
	s.Log["FB Subscription URL"] = admin.NewLog(true, &admin.FBSubscriptionURL{
		U: SubURL,
		V: v,
	})
	res, err := http.PostForm(SubURL, v)
	ok := err == nil && res.StatusCode == http.StatusOK
	s.Log["FB Subscription Response"] = admin.NewLog(ok, res)
}

func (fbs *Stub) VerifySub(s *admin.Status, w http.ResponseWriter, r *http.Request) (int, []byte) {
	ok := r.FormValue("hub.mode") == "subscribe" && r.FormValue("hub.verify_token") == fbs.verifyToken
	s.Log["FB Verification Request"] = admin.NewLog(ok, r)
	if !ok {
		return http.StatusBadRequest, []byte("hub.mode or verify_token mismatch")
	}
	return http.StatusOK, []byte(r.FormValue("hub.challenge"))
}

type SubRequestEntry struct {
	Id            string   `json:"id"`
	ChangedFields []string `json:"changed_fields"`
}

type SubRequest struct {
	Object string            `json:"object"`
	Entry  []SubRequestEntry `json:"entry"`
}

func (fbs *Stub) ProcessSub(sb *SubRequest, wlog weblogger.Logger, w http.ResponseWriter, r *http.Request) (int, error) {
	wlog.Remotef("Entry: %v", sb)
	return http.StatusOK, nil
}

func HubSign(c martini.Context, r *http.Request, res http.ResponseWriter) {
	sig := r.Header.Get("X-Hub-Signature")[5:]
	if sig == "" {
		http.Error(res, "Bad Signed Request", http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(res, "Bad Signed Request", http.StatusBadRequest)
	}
	r.Body.Close()

	mac := hmac.New(sha1.New, []byte(*ClientSecret))
	mac.Write(body)
	s := hex.EncodeToString(mac.Sum(nil))
	if s != sig {
		http.Error(res, "Bad Signed Request", http.StatusBadRequest)
	}

	fmt.Println(string(body))
	sb := &SubRequest{}
	json.Unmarshal(body, sb)
	c.Map(sb)
}
