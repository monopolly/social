package social

import (
	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
	"github.com/valyala/fasthttp"
)

// Facebook reads the profile behind a graph api access token.
func Facebook(token string) (u User, err errors.E) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://graph.facebook.com/v21.0/me?fields=id,name,first_name,last_name,email,picture")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Authorization", "Bearer "+token)

	if er := fasthttp.DoTimeout(req, resp, httpTimeout); er != nil {
		err = errors.Connection(er)
		err.AddPoint()
		return
	}

	b := resp.Body()

	if e := jsons.String(b, "error", "message"); e != "" {
		err = errors.Access(e)
		err.Set("origin", jsons.String(b, "error", "type"))
		err.AddPoint()
		return
	}

	if code := resp.StatusCode(); code != fasthttp.StatusOK {
		err = errors.Access(code)
		err.Set("origin", string(b))
		err.AddPoint()
		return
	}

	u.ID = jsons.String(b, "id")
	u.Key = jsons.String(b, "id")
	u.Token = token
	u.Name = jsons.String(b, "first_name")
	if u.Name == "" {
		u.Name = jsons.String(b, "name")
	}
	u.Family = jsons.String(b, "last_name")
	u.Email = jsons.String(b, "email")
	u.Image = jsons.String(b, "picture", "data", "url")
	u.Source = string(b)

	// facebook only ever hands out a verified email
	u.Verified = u.Email != ""
	return
}
