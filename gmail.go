package social

import (
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwt"

	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
)

func Gmail(token string, appid string) (u User, err errors.E) {
	t, er := jwt.ParseString(token)
	if er != nil {
		err = errors.Parse("jwt")
		return
	}

	er = jwt.Validate(t)
	if er != nil {
		err = errors.Valid("jwt")
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
