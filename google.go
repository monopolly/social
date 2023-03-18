package social

import (
	"fmt"

	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
	"github.com/valyala/fasthttp"
)

// google Bearer token auth
func Google(token string) (u *User, err errors.E) {
	// curl -H 'Authorization: Bearer $ACCESS_TOKEN' https://www.googleapis.com/oauth2/v3/tokeninfo
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI("https://www.googleapis.com/oauth2/v3/userinfo")
	req.Header.SetMethod("GET")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	er := fasthttp.Do(req, resp)
	if er != nil {
		err = errors.Access(er)
		return
	}

	b := resp.Body()
	fmt.Println(string(b))
	u = new(User)
	u.Email = jsons.String(b, "email")
	u.Key = u.Email
	u.Token = u.Email
	u.Name = jsons.String(b, "name")
	u.Image = jsons.String(b, "picture")
	u.Verified = jsons.Bool(b, "email_verified")
	u.Lang = jsons.String(b, "locale")
	return
}

/*
func Gmail(token string, appid string) (u User, err errors.E) {
	t, er := jwt.Parse(token, func(t *jwt.Token) (c interface{}, er error) {
		fmt.Println(t)
		return
	})
	if er != nil {
		fmt.Println(er)
		err = errors.Parse("jwt")
		return
	}

	b := jsons.Marshal(t)
	u.Email = jsons.String(b, "email")
	if u.Email == "" {
		err = errors.Email()
		return
	}

	u.Key = u.Email
	u.Name = jsons.String(b, "name")

	if !strings.Contains(jsons.String(b, "iss"), "accounts.google.com") {
		err = errors.Valid("iss")
		return
	}

	if time.Now().Unix() > jsons.Int64(b, "exp") {
		err = errors.Valid("exp")
		return
	}

	return
}
*/
