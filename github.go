package social

import (
	"github.com/valyala/fasthttp"

	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
)

const githubAgent = "monopolly-social"

// Github exchanges the oauth code for an access token and reads the profile.
func Github(token, clientID, secret string) (u User, err errors.E) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 1. code -> access token. github answers with a form encoded body
	// unless the json content type is asked for by header.
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)
	args.Set("client_id", clientID)
	args.Set("client_secret", secret)
	args.Set("code", token)

	req.SetRequestURI("https://github.com/login/oauth/access_token")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.SetUserAgent(githubAgent)
	req.SetBody(args.QueryString())

	b, err := githubDo(req, resp)
	if err != nil {
		return
	}

	// github answers 200 with an error field when the code is bad
	if e := jsons.String(b, "error"); e != "" {
		err = errors.Token(e)
		err.Set("origin", jsons.String(b, "error_description"))
		err.AddPoint()
		return
	}

	access := jsons.String(b, "access_token")
	if access == "" {
		err = errors.Token("access_token not found")
		err.AddPoint()
		return
	}

	// 2. access token -> profile
	req.Reset()
	resp.Reset()
	req.SetRequestURI("https://api.github.com/user")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+access)
	req.Header.SetUserAgent(githubAgent)

	b, err = githubDo(req, resp)
	if err != nil {
		return
	}

	u.ID = jsons.Int64(b, "id")
	u.Key = jsons.String(b, "id")
	u.Token = access
	u.Name = jsons.String(b, "name")
	if u.Name == "" {
		u.Name = jsons.String(b, "login")
	}
	u.Email = jsons.String(b, "email")
	u.Image = jsons.String(b, "avatar_url")
	u.Source = string(b)

	// 3. the profile hides the email when the user keeps it private,
	// the user:email scope gives it through a separate endpoint
	if u.Email == "" {
		req.Reset()
		resp.Reset()
		req.SetRequestURI("https://api.github.com/user/emails")
		req.Header.SetMethod(fasthttp.MethodGet)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+access)
		req.Header.SetUserAgent(githubAgent)

		b, err = githubDo(req, resp)
		if err != nil {
			return
		}

		for _, v := range jsons.Array(b, "@this") {
			raw := []byte(v.Raw)
			if !jsons.Bool(raw, "primary") {
				continue
			}
			u.Email = jsons.String(raw, "email")
			u.Verified = jsons.Bool(raw, "verified")
			break
		}
	} else {
		u.Verified = true
	}

	if u.Email == "" {
		err = errors.Server("email not found")
		err.AddPoint()
	}
	return
}

func githubDo(req *fasthttp.Request, resp *fasthttp.Response) (b []byte, err errors.E) {
	if er := fasthttp.DoTimeout(req, resp, httpTimeout); er != nil {
		err = errors.Connection(er)
		err.AddPoint()
		return
	}
	if code := resp.StatusCode(); code != fasthttp.StatusOK {
		err = errors.Connection(code)
		err.Set("origin", string(resp.Body()))
		err.AddPoint()
		return
	}
	return resp.Body(), nil
}
