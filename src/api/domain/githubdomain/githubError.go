package githubdomain

//GithubErrorResponse holds details about the errors received back from Github
type GithubErrorResponse struct {
	StatusCode       int           `json:"status_code"`
	Message          string        `json:"message"`
	DocumentationURL string        `json:"documentation_url"`
	Errors           []GithubError `json:"errors"`
}

func (r GithubErrorResponse) Error() string {
	return r.Message
}

//GithubError holds details about a single error from Github
type GithubError struct {
	Resource string `json:"resource"`
	Code     string `json:"code"`
	Field    string `json:"field"`
	Message  string `json:"message"`
}
