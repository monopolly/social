package social

import (
	"context"
	"fmt"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/monopolly/errors"
)

func Apple(token, privateKey, teamID, clientID, keyID string) (u User, err errors.E) {

	var er error
	secret, er := apple.GenerateClientSecret(privateKey, teamID, clientID, keyID)
	if er != nil {
		//fmt.Println("error generating secret: " + err.Error())
		err = errors.Internal(er)
		err.SetCodeLine()
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
		//fmt.Println(err.Error())
		err = errors.Token(er)
		return
	}

	if resp.Error != "" {
		err = errors.Internal(resp.Error)
		return
	}

	// Get the unique user ID
	u.ID, er = apple.GetUniqueID(resp.IDToken)
	if er != nil {
		//fmt.Println("failed to get unique ID: " + )
		err = errors.Internal(er)
		return
	}

	// Get the email
	claim, er := apple.GetClaims(resp.IDToken)
	if err != nil {
		err = errors.Internal(er)
		return
	}

	u.Email = fmt.Sprint((*claim)["email"])
	if u.Email == "" {
		err = errors.Email("email not found")
		return
	}
	/* emailVerified := (*claim)["email_verified"]
	isPrivateEmail := (*claim)["is_private_email"]
	*/
	return
}
