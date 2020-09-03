package github

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/shurcooL/githubv4"
	"strings"
)

func dataSourceGithubUsers() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			// Input
			IGNORE_MISSING: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			USER_LOGINS: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Output
			USERS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						USER_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						USER_LOGIN: {
							Type:     schema.TypeString,
							Computed: true,
						},
						USER_IS_SITE_ADMIN: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						USER_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},

		Read: dataSourceGithubUsersRead,
	}
}

func dataSourceGithubUsersRead(d *schema.ResourceData, meta interface{}) error {
	data, err := usersResourceData(d)
	if err != nil {
		return err
	}

	h := sha1.New()
	upns := make([]string, 0)
	users := make([]interface{}, 0)
	for _, login := range data.UserLogins {
		var query struct {
			User User `graphql:"user(login: $login)"`
		}
		variables := map[string]interface{}{
			"login": githubv4.String(login),
		}

		ctx := context.Background()
		client := meta.(*Organization).Client
		err := client.Query(ctx, &query, variables)
		if err != nil {
			if !data.IgnoreMissing && strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
				return err
			}
		}

		if query.User != (User{}) {
			user := map[string]interface{}{
				USER_ID:            query.User.ID,
				USER_LOGIN:         query.User.Login,
				USER_IS_SITE_ADMIN: query.User.IsSiteAdmin,
				USER_NAME:          query.User.Name,
			}

			users = append(users, user)
			upns = append(upns, query.User.ID.(string))
		}
	}

	err = d.Set(USERS, users)
	if err != nil {
		return err
	}

	if _, err := h.Write([]byte(strings.Join(upns, "-"))); err != nil {
		return fmt.Errorf("unable to compute hash: %v", err)
	}

	d.SetId("users#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))

	return nil
}
