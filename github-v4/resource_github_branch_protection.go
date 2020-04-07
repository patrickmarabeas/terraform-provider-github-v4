package github

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/shurcooL/githubv4"
)

func resourceGithubBranchProtection() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			// Input
			REPOSITORY_ID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			PROTECTION_PATTERN: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			PROTECTION_IS_ADMIN_ENFORCED: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			PROTECTION_REQUIRES_COMMIT_SIGNATURES: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			PROTECTION_REQUIRES_APPROVING_REVIEWS: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PROTECTION_REQUIRED_APPROVING_REVIEW_COUNT: {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 6),
						},
						PROTECTION_REQUIRES_CODE_OWNER_REVIEWS: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						PROTECTION_DISMISSES_STALE_REVIEWS: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						PROTECTION_RESTRICTS_REVIEW_DISMISSALS: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									PROTECTION_ACTOR_IDS: {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			PROTECTION_REQUIRES_STATUS_CHECKS: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PROTECTION_REQUIRES_STRICT_STATUS_CHECKS: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						PROTECTION_REQUIRED_STATUS_CHECK_CONTEXTS: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			PROTECTION_RESTRICTS_PUSHES: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PROTECTION_ACTOR_IDS: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},

		Create: resourceGithubBranchProtectionCreate,
		Read:   resourceGithubBranchProtectionRead,
		Update: resourceGithubBranchProtectionUpdate,
		Delete: resourceGithubBranchProtectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceGithubBranchProtectionV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceGithubBranchProtectionUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceGithubBranchProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	var mutate struct {
		CreateBranchProtectionRule struct {
			BranchProtectionRule struct {
				ID githubv4.ID
			}
		} `graphql:"createBranchProtectionRule(input: $input)"`
	}
	data, err := branchProtectionResourceData(d, meta)
	if err != nil {
		return err
	}
	input := githubv4.CreateBranchProtectionRuleInput{
		DismissesStaleReviews:        githubv4.NewBoolean(githubv4.Boolean(data.DismissesStaleReviews)),
		IsAdminEnforced:              githubv4.NewBoolean(githubv4.Boolean(data.IsAdminEnforced)),
		Pattern:                      githubv4.String(data.Pattern),
		PushActorIDs:                 githubv4NewIDSlice(githubv4IDSlice(data.PushActorIDs)),
		RepositoryID:                 githubv4.NewID(githubv4.ID(data.RepositoryID)),
		RequiredApprovingReviewCount: githubv4.NewInt(githubv4.Int(data.RequiredApprovingReviewCount)),
		RequiredStatusCheckContexts:  githubv4NewStringSlice(githubv4StringSlice(data.RequiredStatusCheckContexts)),
		RequiresApprovingReviews:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresApprovingReviews)),
		RequiresCodeOwnerReviews:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresCodeOwnerReviews)),
		RequiresCommitSignatures:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresCommitSignatures)),
		RequiresStatusChecks:         githubv4.NewBoolean(githubv4.Boolean(data.RequiresStatusChecks)),
		RequiresStrictStatusChecks:   githubv4.NewBoolean(githubv4.Boolean(data.RequiresStrictStatusChecks)),
		RestrictsPushes:              githubv4.NewBoolean(githubv4.Boolean(data.RestrictsPushes)),
		RestrictsReviewDismissals:    githubv4.NewBoolean(githubv4.Boolean(data.RestrictsReviewDismissals)),
		ReviewDismissalActorIDs:      githubv4NewIDSlice(githubv4IDSlice(data.ReviewDismissalActorIDs)),
	}

	ctx := context.Background()
	client := meta.(*Organization).Client
	err = client.Mutate(ctx, &mutate, input, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s", mutate.CreateBranchProtectionRule.BranchProtectionRule.ID))

	return resourceGithubBranchProtectionRead(d, meta)
}

func resourceGithubBranchProtectionRead(d *schema.ResourceData, meta interface{}) error {
	var query struct {
		Node struct {
			Node BranchProtectionRule `graphql:"... on BranchProtectionRule"`
		} `graphql:"node(id: $id)"`
	}
	variables := map[string]interface{}{
		"id": d.Id(),
	}

	ctx := context.WithValue(context.Background(), "id", d.Id())
	client := meta.(*Organization).Client
	err := client.Query(ctx, &query, variables)

	return err
}

func resourceGithubBranchProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	var mutate struct {
		UpdateBranchProtectionRule struct {
			BranchProtectionRule struct {
				ID githubv4.ID
			}
		} `graphql:"updateBranchProtectionRule(input: $input)"`
	}
	data, err := branchProtectionResourceData(d, meta)
	if err != nil {
		return err
	}
	input := githubv4.UpdateBranchProtectionRuleInput{
		BranchProtectionRuleID:       d.Id(),
		DismissesStaleReviews:        githubv4.NewBoolean(githubv4.Boolean(data.DismissesStaleReviews)),
		IsAdminEnforced:              githubv4.NewBoolean(githubv4.Boolean(data.IsAdminEnforced)),
		Pattern:                      githubv4.NewString(githubv4.String(data.Pattern)),
		PushActorIDs:                 githubv4NewIDSlice(githubv4IDSlice(data.PushActorIDs)),
		RequiredApprovingReviewCount: githubv4.NewInt(githubv4.Int(data.RequiredApprovingReviewCount)),
		RequiredStatusCheckContexts:  githubv4NewStringSlice(githubv4StringSlice(data.RequiredStatusCheckContexts)),
		RequiresApprovingReviews:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresApprovingReviews)),
		RequiresCodeOwnerReviews:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresCodeOwnerReviews)),
		RequiresCommitSignatures:     githubv4.NewBoolean(githubv4.Boolean(data.RequiresCommitSignatures)),
		RequiresStatusChecks:         githubv4.NewBoolean(githubv4.Boolean(data.RequiresStatusChecks)),
		RequiresStrictStatusChecks:   githubv4.NewBoolean(githubv4.Boolean(data.RequiresStrictStatusChecks)),
		RestrictsPushes:              githubv4.NewBoolean(githubv4.Boolean(data.RestrictsPushes)),
		RestrictsReviewDismissals:    githubv4.NewBoolean(githubv4.Boolean(data.RestrictsReviewDismissals)),
		ReviewDismissalActorIDs:      githubv4NewIDSlice(githubv4IDSlice(data.ReviewDismissalActorIDs)),
	}

	ctx := context.WithValue(context.Background(), "id", d.Id())
	client := meta.(*Organization).Client
	err = client.Mutate(ctx, &mutate, input, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s", mutate.UpdateBranchProtectionRule.BranchProtectionRule.ID))

	return resourceGithubBranchProtectionRead(d, meta)
}

func resourceGithubBranchProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	var mutate struct {
		DeleteBranchProtectionRule struct { // Empty struct does not work
			ClientMutationId githubv4.ID
		} `graphql:"deleteBranchProtectionRule(input: $input)"`
	}
	input := githubv4.DeleteBranchProtectionRuleInput{
		BranchProtectionRuleID: d.Id(),
	}

	ctx := context.WithValue(context.Background(), "id", d.Id())
	client := meta.(*Organization).Client
	err := client.Mutate(ctx, &mutate, input, nil)

	return err
}
