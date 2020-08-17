package github

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/shurcooL/githubv4"
)

const (
	PROTECTION_DISMISSES_STALE_REVIEWS         = "dismiss_stale_reviews"
	PROTECTION_IS_ADMIN_ENFORCED               = "enforce_admins"
	PROTECTION_PATTERN                         = "pattern"
	PROTECTION_REQUIRED_APPROVING_REVIEW_COUNT = "required_approving_review_count"
	PROTECTION_REQUIRED_STATUS_CHECK_CONTEXTS  = "contexts"
	PROTECTION_REQUIRES_APPROVING_REVIEWS      = "required_pull_request_reviews"
	PROTECTION_REQUIRES_CODE_OWNER_REVIEWS     = "require_code_owner_reviews"
	PROTECTION_REQUIRES_COMMIT_SIGNATURES      = "require_signed_commits"
	PROTECTION_REQUIRES_STATUS_CHECKS          = "required_status_checks"
	PROTECTION_REQUIRES_STRICT_STATUS_CHECKS   = "strict"
	PROTECTION_RESTRICTS_PUSHES                = "push_restrictions"
	PROTECTION_RESTRICTS_REVIEW_DISMISSALS     = "dismissal_restrictions"
)

type Actor struct {
	ID   githubv4.ID
	Name githubv4.String
}

type BranchProtectionRule struct {
	Repository struct {
		ID   githubv4.String
		Name githubv4.String
	}
	PushAllowances struct {
		Nodes []struct {
			Actor struct {
				// `App` is not supported (at least for GitHub App Installation tokens)
				// Seem to be unable to provide the necessary permissions.
				Team Actor `graphql:"... on Team"`
				User Actor `graphql:"... on User"`
			}
		}
	} `graphql:"pushAllowances(first: 100)"`
	ReviewDismissalAllowances struct {
		Nodes []struct {
			Actor struct {
				Team Actor `graphql:"... on Team"`
				User Actor `graphql:"... on User"`
			}
		}
	} `graphql:"reviewDismissalAllowances(first: 100)"`
	DismissesStaleReviews        githubv4.Boolean
	ID                           githubv4.ID
	IsAdminEnforced              githubv4.Boolean
	Pattern                      githubv4.String
	RequiredApprovingReviewCount githubv4.Int
	RequiredStatusCheckContexts  []githubv4.String
	RequiresApprovingReviews     githubv4.Boolean
	RequiresCodeOwnerReviews     githubv4.Boolean
	RequiresCommitSignatures     githubv4.Boolean
	RequiresStatusChecks         githubv4.Boolean
	RequiresStrictStatusChecks   githubv4.Boolean
	RestrictsPushes              githubv4.Boolean
	RestrictsReviewDismissals    githubv4.Boolean
}

type BranchProtectionResourceData struct {
	BranchProtectionRuleID       string
	DismissesStaleReviews        bool
	IsAdminEnforced              bool
	Pattern                      string
	PushActorIDs                 []string
	RepositoryID                 string
	RequiredApprovingReviewCount int
	RequiredStatusCheckContexts  []string
	RequiresApprovingReviews     bool
	RequiresCodeOwnerReviews     bool
	RequiresCommitSignatures     bool
	RequiresStatusChecks         bool
	RequiresStrictStatusChecks   bool
	RestrictsPushes              bool
	RestrictsReviewDismissals    bool
	ReviewDismissalActorIDs      []string
}

func branchProtectionResourceData(d *schema.ResourceData, meta interface{}) (BranchProtectionResourceData, error) {
	data := BranchProtectionResourceData{}

	if v, ok := d.GetOk(REPOSITORY_ID); ok {
		data.RepositoryID = v.(string)
	}

	if v, ok := d.GetOk(PROTECTION_PATTERN); ok {
		data.Pattern = v.(string)
	}

	if v, ok := d.GetOk(PROTECTION_IS_ADMIN_ENFORCED); ok {
		data.IsAdminEnforced = v.(bool)
	}

	if v, ok := d.GetOk(PROTECTION_REQUIRES_COMMIT_SIGNATURES); ok {
		data.RequiresCommitSignatures = v.(bool)
	}

	if v, ok := d.GetOk(PROTECTION_REQUIRES_APPROVING_REVIEWS); ok {
		vL := v.([]interface{})
		if len(vL) > 1 {
			return BranchProtectionResourceData{},
				fmt.Errorf("error multiple %s declarations", PROTECTION_REQUIRES_APPROVING_REVIEWS)
		}
		for _, v := range vL {
			if v == nil {
				break
			}

			data.RequiresApprovingReviews = true

			m := v.(map[string]interface{})
			if v, ok := m[PROTECTION_REQUIRED_APPROVING_REVIEW_COUNT]; ok {
				data.RequiredApprovingReviewCount = v.(int)
			}
			if v, ok := m[PROTECTION_DISMISSES_STALE_REVIEWS]; ok {
				data.DismissesStaleReviews = v.(bool)
			}
			if v, ok := m[PROTECTION_REQUIRES_CODE_OWNER_REVIEWS]; ok {
				data.RequiresCodeOwnerReviews = v.(bool)
			}
			if v, ok := m[PROTECTION_RESTRICTS_REVIEW_DISMISSALS]; ok {
				reviewDismissalActorIDs := make([]string, 0)
				vL := v.(*schema.Set).List()
				for _, v := range vL {
					reviewDismissalActorIDs = append(reviewDismissalActorIDs, v.(string))
				}
				if len(reviewDismissalActorIDs) > 0 {
					data.ReviewDismissalActorIDs = reviewDismissalActorIDs
					data.RestrictsReviewDismissals = true
				}
			}
		}
	}

	if v, ok := d.GetOk(PROTECTION_REQUIRES_STATUS_CHECKS); ok {
		vL := v.([]interface{})
		if len(vL) > 1 {
			return BranchProtectionResourceData{},
				fmt.Errorf("error multiple %s declarations", PROTECTION_REQUIRES_STATUS_CHECKS)
		}
		for _, v := range vL {
			if v == nil {
				break
			}

			m := v.(map[string]interface{})
			if v, ok := m[PROTECTION_REQUIRES_STRICT_STATUS_CHECKS]; ok {
				data.RequiresStrictStatusChecks = v.(bool)
			}

			data.RequiredStatusCheckContexts = expandNestedSet(m, PROTECTION_REQUIRED_STATUS_CHECK_CONTEXTS)
			if len(data.RequiredStatusCheckContexts) > 0 {
				data.RequiresStatusChecks = true
			}
		}
	}

	if v, ok := d.GetOk(PROTECTION_RESTRICTS_PUSHES); ok {
		pushActorIDs := make([]string, 0)
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			pushActorIDs = append(pushActorIDs, v.(string))
		}
		if len(pushActorIDs) > 0 {
			data.PushActorIDs = pushActorIDs
			data.RestrictsPushes = true
		}
	}

	return data, nil
}

func setApprovingReviews(d *schema.ResourceData, protection BranchProtectionRule) interface{} {
	if protection.RequiresApprovingReviews == false {
		return nil
	}

	dismissalAllowances := protection.ReviewDismissalAllowances.Nodes
	dismissalActors := make([]interface{}, 0, len(dismissalAllowances))
	for _, d := range dismissalAllowances {
		if d.Actor.Team != (Actor{}) {
			dismissalActors = append(dismissalActors, d.Actor.Team.ID)
		}
		if d.Actor.User != (Actor{}) {
			dismissalActors = append(dismissalActors, d.Actor.Team.ID)
		}
	}

	approvalReviews := []interface{}{
		map[string]interface{}{
			PROTECTION_REQUIRED_APPROVING_REVIEW_COUNT: protection.RequiredApprovingReviewCount,
			PROTECTION_REQUIRES_CODE_OWNER_REVIEWS:     protection.RequiresCodeOwnerReviews,
			PROTECTION_DISMISSES_STALE_REVIEWS:         protection.DismissesStaleReviews,
			PROTECTION_RESTRICTS_REVIEW_DISMISSALS: 	dismissalActors,
		},
	}

	return approvalReviews
}

func setStatusChecks(d *schema.ResourceData, protection BranchProtectionRule) interface{} {
	if protection.RequiresStatusChecks == false {
		return nil
	}

	statusChecks := []interface{}{
		map[string]interface{}{
			PROTECTION_REQUIRES_STRICT_STATUS_CHECKS:  protection.RequiresStrictStatusChecks,
			PROTECTION_REQUIRED_STATUS_CHECK_CONTEXTS: protection.RequiredStatusCheckContexts,
		},
	}

	return statusChecks
}

func setPushes(d *schema.ResourceData, protection BranchProtectionRule) []string {
	if protection.RestrictsPushes == false {
		return nil
	}

	pushAllowances := protection.PushAllowances.Nodes
	pushActors := make([]string, 0, len(pushAllowances))
	for _, p := range pushAllowances {
		if p.Actor.Team != (Actor{}) {
			pushActors = append(pushActors, p.Actor.Team.ID.(string))
		}
		if p.Actor.User != (Actor{}) {
			pushActors = append(pushActors, p.Actor.Team.ID.(string))
		}
	}

	return pushActors
}

func getBranchProtectionID(name string, pattern string, meta interface{}) (githubv4.ID, error) {
	var query struct {
		Node struct {
			Repository struct {
				BranchProtectionRules struct {
					Nodes []struct {
						ID      string
						Pattern string
					}
					PageInfo PageInfo
				} `graphql:"branchProtectionRules(first: $first, after: $cursor)"`
				ID string
			} `graphql:"... on Repository"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":  githubv4.String(meta.(*Organization).Name),
		"name":   githubv4.String(name),
		"first":  githubv4.Int(100),
		"cursor": (*githubv4.String)(nil),
	}

	ctx := context.Background()
	client := meta.(*Organization).Client

	var allRules []struct {
		ID      string
		Pattern string
	}
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}

		allRules = append(allRules, query.Node.Repository.BranchProtectionRules.Nodes...)

		if !query.Node.Repository.BranchProtectionRules.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Node.Repository.BranchProtectionRules.PageInfo.EndCursor)
	}

	var id string
	for i := range allRules {
		if allRules[i].Pattern == pattern {
			id = allRules[i].ID
			break
		}
	}

	return id, nil
}
