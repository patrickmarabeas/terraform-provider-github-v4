package github

import "github.com/hashicorp/terraform/helper/schema"

func resourceGithubBranchProtectionV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGithubBranchProtectionUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	repositoryName := rawState["repository"].(string)
	repositoryID, err := getRepositoryID(repositoryName, meta)
	if err != nil {
		return nil, err
	}

	branch := rawState["branch"].(string)
	branchProtectionRuleID, err := getBranchProtectionID(repositoryName, branch, meta)
	if err != nil {
		return nil, err
	}

	rawState["id"] = branchProtectionRuleID
	rawState[REPOSITORY_ID] = repositoryID
	rawState[PROTECTION_PATTERN] = branch

	return rawState, nil
}
