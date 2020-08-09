package github

import "github.com/shurcooL/githubv4"

const (
	USER_ID            = "user_id"
	USER_IS_SITE_ADMIN = "is_site_admin"
	USER_LOGIN         = "login"
	USER_NAME          = "name"
	USER_ROLE          = "role"
	USER_PERMISSION    = "permission"
)

type User struct {
	ID          githubv4.ID
	IsSiteAdmin githubv4.Boolean
	Login       githubv4.String
	Name        githubv4.String
}
