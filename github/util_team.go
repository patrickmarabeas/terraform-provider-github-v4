package github

import "github.com/shurcooL/githubv4"

const (
	TEAM_CHILD_TEAMS = "child_teams"
	TEAM_DESCRIPTION = "description"
	TEAM_ID          = "team_id"
	TEAM_MEMBERS     = "members"
	TEAM_NAME        = "name"
	TEAM_PARENT_TEAM = "parent_team"
	TEAM_PRIVACY     = "privacy"
	TEAM_SLUG        = "slug"
)

type Team struct {
	ChildTeams struct {
		Nodes []struct {
			ID   githubv4.ID
			Slug githubv4.String
		}
		PageInfo PageInfo
	} `graphql:"childTeams(first: $childTeamFirst, after: $childTeamCursor, immediateOnly: $immediateOnly)"`
	Members struct {
		Edges []struct {
			Node User
			Role githubv4.TeamMemberRole
		}
		PageInfo PageInfo
	} `graphql:"members(first: $membersFirst, after: $membersCursor)"`
	ParentTeam struct {
		ID   githubv4.ID
		Slug githubv4.String
	}
	Description githubv4.String
	ID          githubv4.ID
	Name        githubv4.String
	Privacy     githubv4.TeamPrivacy
}
