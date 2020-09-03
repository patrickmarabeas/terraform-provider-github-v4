package github

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/shurcooL/githubv4"
)

const (
	IGNORE_MISSING     = "ignore_missing"
	USER_ID            = "user_id"
	USER_IS_SITE_ADMIN = "is_site_admin"
	USER_LOGIN         = "login"
	USER_LOGINS        = "logins"
	USER_NAME          = "name"
	USER_ROLE          = "role"
	USER_PERMISSION    = "permission"
	USERS              = "users"
)

type User struct {
	ID          githubv4.ID
	IsSiteAdmin githubv4.Boolean
	Login       githubv4.String
	Name        githubv4.String
}

type UsersResourceData struct {
	IgnoreMissing bool
	UserLogins    []string
}

func usersResourceData(d *schema.ResourceData) (UsersResourceData, error) {
	data := UsersResourceData{}

	if v, ok := d.GetOk(IGNORE_MISSING); ok {
		data.IgnoreMissing = v.(bool)
	}

	if v, ok := d.GetOk(USER_LOGINS); ok {

		users := make([]string, 0)
		vL := v.([]interface{})
		for _, v := range vL {
			users = append(users, v.(string))
		}
		data.UserLogins = users
	}

	return data, nil
}
