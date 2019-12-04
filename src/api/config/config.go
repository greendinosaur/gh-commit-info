package config

import "os"

const (
	apiGitHubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
	//LogLevel to be used across the application
	LogLevel = "info"
)

var (
	githubAccessToken = os.Getenv(apiGitHubAccessToken)
)

//GetGithubAccessToken returns the access token used to access the particular user account in the Github API
func GetGithubAccessToken() string {
	return githubAccessToken
}
