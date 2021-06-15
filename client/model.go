package client

type OidcTokenRequest struct {
	Client_id     string
	Response_type string `default:"id_token"`
	Grant_type    string `default:"password"`
	Scope         string `default:"openid"`
	Username      string
	Password      string
}
type OidcTokenResponse struct {
	IdToken string `json:"id_token,omitempty"`
}

type ClientPrt struct {
	Url         string
	Method      string
	ContentType string
	Bearer      string
	Password    string
	User        string
	Body        interface{}
}

type ArtifactResponse struct {
	Digest string `json:"digest,omitempty"`
}
