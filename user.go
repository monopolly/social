package social

import "time"

// deadline for every provider request, none of them may hang forever
const httpTimeout = 10 * time.Second

type User struct {
	ID       any    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Family   string `json:"family,omitempty"`
	Key      string `json:"key,omitempty"`
	Token    string `json:"token,omitempty"`
	Email    string `json:"email,omitempty"`
	Image    string `json:"image,omitempty"`
	Lang     string `json:"lang,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Source   string `json:"source,omitempty"`
}
