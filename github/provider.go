package github

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			PROVIDER_BASE_URL: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_BASE_URL", "https://api.github.com/"),
				Description: "The GitHub Root API URL.",
			},
			PROVIDER_ORGANIZATION: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_ORGANIZATION", nil),
				Description: "The target GitHub organization to manage.",
			},
			PROVIDER_TOKEN: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_TOKEN", nil),
				Description: "The GitHub access token.",
				Sensitive:   true,
			},
			PROVIDER_APP: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PROVIDER_APP_PEM: {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("GITHUB_APP_PEM", nil),
							Description: "The GitHub App PEM string.",
							Sensitive:   true,
						},
						PROVIDER_APP_ID: {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("GITHUB_APP_ID", nil),
							Description: "The GitHub App ID.",
						},
						PROVIDER_APP_INSTALLATION_ID: {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("GITHUB_APP_INSTALLATION_ID", nil),
							Description: "The GitHub App installation instance ID.",
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"github_branch_protection": resourceGithubBranchProtection(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"github_user": dataSourceGithubUser(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		var (
			baseURL      = d.Get(PROVIDER_BASE_URL).(string)
			organization = d.Get(PROVIDER_ORGANIZATION).(string)
			token        = d.Get(PROVIDER_TOKEN).(string)
			appPEM       = ""
			appID        = ""
			appInstID    = ""
		)

		if v, ok := d.GetOk(PROVIDER_APP); ok {
			vL := v.([]interface{})
			if len(vL) > 1 {
				return nil, fmt.Errorf("error: multiple %s declarations", "app")
			}
			for _, v := range vL {
				if v == nil {
					break
				}

				m := v.(map[string]interface{})
				if v, ok := m[PROVIDER_APP_PEM]; ok {
					appPEM = v.(string)
				}
				if v, ok := m[PROVIDER_APP_ID]; ok {
					appID = v.(string)
				}
				if v, ok := m[PROVIDER_APP_INSTALLATION_ID]; ok {
					appInstID = v.(string)
				}
			}
		}

		config := Config{
			BaseURL:        baseURL,
			Organization:   organization,
			Token:          token,
			Pem:            appPEM,
			AppID:          appID,
			InstallationID: appInstID,
		}

		meta, err := config.Clients()
		if err != nil {
			return nil, err
		}

		meta.(*Organization).StopContext = p.StopContext()

		return meta, nil
	}
}
