package core

type PrEvent struct {
	Action      string      `json:"action" required:"true"`
	Before      string      `json:"before"`
	After       string      `json:"after"`
	PullRequest PullRequest `json:"pull_request" validate:"required"`
}

type PullRequest struct {
	Number    int    `json:"number" validate:"required"`
	HtmlUrl   string `json:"html_url" validate:"required"`
	StatusUrl string `json:"statuses_url" validate:"required"`
	Head      Commit `json:"head" validate:"required"`
	Base      Commit `json:"base" validate:"required"`
}

type Commit struct {
	Ref  string     `json:"ref" validate:"required"`
	Sha  string     `json:"sha" validate:"required"`
	User User       `json:"user" validate:"required"`
	Repo Repository `json:"repo" validate:"required"`
}

type User struct {
	Login     string `json:"login" validate:"required"`
	AvatarUrl string `json:"avatar_url" validate:"required"`
}

type Repository struct {
	HtmlUrl string `json:"html_url" validate:"required"`
}

type PushEvent struct {
	Ref    string `json:"ref" validate:"required"`
	Branch string
	Before string     `json:"before" validate:"required"`
	After  string     `json:"after" validate:"required"`
	Repo   Repository `json:"repository" validate:"required"`
}

type BuildEvent struct {
	Commit string
	Status int
}

// BuildStatus
const (
	BuildQueued int = iota
	BuildMachineStarted
	BuildRepoCloned
	BuildScriptCopied

	BuildSucceed
	BuildFailed
	BuildCancelled
)
