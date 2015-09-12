package fb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (fbs *Stub) fetch(id uint64, fields string) (json.RawMessage, error) {
	v := url.Values{}
	v.Add(fbs.accessToken[0], fbs.accessToken[1])
	v.Add("fields", fields)
	res, err := http.Get(fmt.Sprintf("%s/%s?%s", EndPoint, id, v.Encode()))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (fbs *Stub) FetchUser(userID uint64) (json.RawMessage, error) {
	return fbs.fetch(userID, "name,email,gender,birthday,locale,location,first_name,last_name")
}

func (fbs *Stub) FetchPlace(placeID uint64) (json.RawMessage, error) {
	return fbs.fetch(placeID, "")
}

func (fbs *Stub) FetchEvent(eventID uint64) (json.RawMessage, error) {
	return fbs.fetch(eventID, "")
}
