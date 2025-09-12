package domain

type RespRedirect struct {
	URL  string `json:"url,omitempty"`
	Resp Resp   `json:"resp,omitempty"`
}

type RedirectLoginInfo struct {
	Platform string
}
