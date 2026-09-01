package social

import (
	"context"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/monopolly/errors"
)

func Apple(token, privateKey, teamID, clientID, keyID string) (u *User, err errors.E) {
	u = new(User)

	secret, er := apple.GenerateClientSecret(privateKey, teamID, clientID, keyID)
	if er != nil {
		err = errors.Server(er)
		err.AddPoint()
		return
	}

	// Generate a new validation client
	client := apple.New()

	vReq := apple.AppValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         token,
	}

	var resp apple.ValidationResponse

	// Do the verification
	er = client.VerifyAppToken(context.Background(), vReq, &resp)
	if er != nil {
		err = errors.Token(er)
		err.AddPoint()
		return
	}

	if resp.Error != "" {
		err = errors.Server(resp.Error)
		err.Set("origin", resp.ErrorDescription)
		err.AddPoint()
		return
	}

	// Get the unique user ID
	u.ID, er = apple.GetUniqueID(resp.IDToken)
	if er != nil {
		err = errors.Server(er)
		err.AddPoint()
		return
	}

	// Get the email
	claim, er := apple.GetClaims(resp.IDToken)
	if er != nil {
		err = errors.Server(er)
		err.AddPoint()
		return
	}

	u.Email, _ = (*claim)["email"].(string)
	if u.Email == "" {
		err = errors.Server("email not found")
		err.AddPoint()
		return
	}

	// apple sends email_verified either as a bool or as a quoted string
	switch v := (*claim)["email_verified"].(type) {
	case bool:
		u.Verified = v
	case string:
		u.Verified = v == "true"
	}

	return
}
