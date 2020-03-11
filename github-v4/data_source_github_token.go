package github

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	TOKEN = "token"
)

func dataSourceGithubToken() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			// Computed
			TOKEN: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Read: resourceGithubTokenRead,
	}
}

func resourceGithubTokenRead(d *schema.ResourceData, meta interface{}) error {
	d.Set(TOKEN, meta.(*Organization).Token)

	d.SetId(fmt.Sprintf("%s/token", meta.(*Organization).Name))

	return nil
}
