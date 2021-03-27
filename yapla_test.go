// Copyright 2021 Iglou.eu
// license that can be found in the LICENSE file

package yapla

import (
	"fmt"
	"testing"
)

var api *Api

func tError(r bool, s string, t *testing.T) {
	if r {
		t.Errorf("%s", s)
		t.Fail()
	}
}

func TestExpireToTime(t *testing.T) {
	ts := ""
	ti, er := expireToTime(ts)
	tError(er == nil, "empty time string expect an error", t)

	ts = "fake"
	ti, er = expireToTime(ts)
	tError(er == nil, "bad time string expect an error", t)

	ts = "Feb 3, 2013 at 7:54pm"
	ti, er = expireToTime(ts)
	tError(er == nil, "other time format expect an error", t)

	ts = "2006-01-02T15:04:05Z"
	ti, er = expireToTime(ts)
	tError(er != nil, fmt.Sprint(er), t)
	tError(ti.String() != "2006-01-02 15:04:05 +0000 UTC", ti.String(), t)

	ts = "2006-01-02T15:04:05"
	ti, er = expireToTime(ts)
	tError(er != nil, fmt.Sprint(er), t)
	tError(ti.String() != "2006-01-02 15:04:05 +0000 UTC", ti.String(), t)

	ts = "2006-01-02 15:04:05"
	ti, er = expireToTime(ts)
	tError(er != nil, fmt.Sprint(er), t)
	tError(ti.String() != "2006-01-02 15:04:05 +0000 UTC", ti.String(), t)
}

func TestNewSession(t *testing.T) {
	var er error

	api, er = NewSession("")
	tError(er == nil, "empty token expect an error", t)

	api, er = NewSession("fake")
	tError(er == nil, "bad token expect an error", t)

	api, er = NewSession(
		"HP1ST252NFKX6Z6RVJ4RKEU23WS2QXSTQHTVYA1JAFWYX306",
		Config{},
	)
	tError(er == nil, "empty url expect an error", t)

	api, er = NewSession(
		"HP1ST252NFKX6Z6RVJ4RKEU23WS2QXSTQHTVYA1JAFWYX306",
		Config{
			URL: "https://duckduckgo.com",
		},
	)
	tError(er == nil, "not yapla api endpoint expect an erro", t)

	api, er = NewSession("HP1ST252NFKX6Z6RVJ4RKEU23WS2QXSTQHTVYA1JAFWYX306")
	tError(er != nil, fmt.Sprint(er), t)
}

func TestLogin(t *testing.T) {
	rep, er := api.LoginMember("moncompte@macompagnie.com", "monp4ssW0R4!")
	tError(er != nil, fmt.Sprint(er), t)
	tError(rep.Result, fmt.Sprint(rep), t)

	rep, er = api.LoginContact("moncompte@macompagnie.com", "monp4ssW0R4!")
	tError(er != nil, fmt.Sprint(er), t)
	tError(rep.Result, fmt.Sprint(rep), t)
}
