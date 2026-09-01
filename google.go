package social

import (
	"fmt"

	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
	"github.com/valyala/fasthttp"
)

/* {
  "sub": "101886101454223860289",
  "name": "Martin Prestone",
  "given_name": "Martin",
  "family_name": "Prestone",
  "picture": "https://lh3.googleusercontent.com/a/ACg8ocLGZJiSoOOW3gS6dVg2RPBCyWGS4H1qsFbDx8ThRVAy=s96-c",
  "email": "wire.common@gmail.com",
  "email_verified": true,
  "locale": "en"
} */

// google Bearer token auth
func Google(token string) (u *User, err errors.E) {
	// curl -H 'Authorization: Bearer $ACCESS_TOKEN' https://www.googleapis.com/oauth2/v3/userinfo
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://www.googleapis.com/oauth2/v3/userinfo")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	er := fasthttp.DoTimeout(req, resp, httpTimeout)
	if er != nil {
		err = errors.Access(er)
		err.AddPoint()
		return
	}

	b := resp.Body()

	ers := jsons.String(b, "error")
	if ers != "" {
		err = errors.Access(ers)
		err.Set("origin", jsons.String(b, "error_description"))
		err.AddPoint()
		return
	}

	if code := resp.StatusCode(); code != fasthttp.StatusOK {
		err = errors.Access(code)
		err.Set("origin", string(b))
		err.AddPoint()
		return
	}

	// fmt.Println(string(b))
	u = new(User)
	u.ID = jsons.String(b, "sub")
	u.Email = jsons.String(b, "email")
	u.Key = u.Email
	u.Token = u.Email
	u.Name = jsons.String(b, "given_name")
	u.Family = jsons.String(b, "family_name")
	if u.Name == "" {
		u.Name = jsons.String(b, "name")
	}
	u.Image = jsons.String(b, "picture")
	u.Verified = jsons.Bool(b, "email_verified")
	u.Lang = jsons.String(b, "locale")
	u.Source = string(b)
	return
}
