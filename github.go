package social

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
)

func Github(token, clientID, secret string) (u User, err errors.E) {

	link := fmt.Sprintf("https://github.com/login/oauth/access_token?scope=user:email&client_id=%s&client_secret=%s&code=%s&accept=json", clientID, secret, token)

	code, b, er := fasthttp.Post(nil, link, nil)
	if er != nil {
		err = errors.Connection(er)
		return
	}
	if code != 200 {
		err = errors.Connection(code)
		return
	}
	u.Email = jsons.String(b, "email")
	u.Name = jsons.String(b, "name")
	u.Key = jsons.String(b, "id")
	return
}
