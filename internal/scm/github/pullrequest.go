package github

import (
	"fmt"

	"github.com/google/go-github/v39/github"

	"github.com/lindell/multi-gitter/internal/scm"
)

func convertPullRequest(pr *github.PullRequest) pullRequest {
	return pullRequest{
		ownerName:   pr.GetBase().GetUser().GetLogin(),
		repoName:    pr.GetBase().GetRepo().GetName(),
		branchName:  pr.GetHead().GetRef(),
		prOwnerName: pr.GetHead().GetUser().GetLogin(),
		prRepoName:  pr.GetHead().GetRepo().GetName(),
		number:      pr.GetNumber(),
		guiURL:      pr.GetHTMLURL(),
	}
}

func convertGraphQLPullRequest(pr graphqlPR) pullRequest {
	combinedStatus := pr.Commits.Nodes[0].Commit.StatusCheckRollup.State
	status := scm.PullRequestStatusUnknown

	if pr.Merged {
		status = scm.PullRequestStatusMerged
	} else if combinedStatus == nil {
		status = scm.PullRequestStatusSuccess
	} else {
		switch *combinedStatus {
		case graphqlPullRequestStatePending:
			status = scm.PullRequestStatusPending
		case graphqlPullRequestStateSuccess:
			status = scm.PullRequestStatusSuccess
		case graphqlPullRequestStateFailure, graphqlPullRequestStateError:
			status = scm.PullRequestStatusError
		}
	}

	return pullRequest{
		ownerName:   pr.BaseRepository.Owner.Login,
		repoName:    pr.BaseRepository.Name,
		branchName:  pr.HeadRefName,
		prOwnerName: pr.HeadRepository.Owner.Login,
		prRepoName:  pr.HeadRepository.Name,
		number:      pr.Number,
		guiURL:      pr.URL,
		status:      status,
	}
}

type pullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	number      int
	guiURL      string
	status      scm.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.number)
}

func (pr pullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.guiURL
}
