package types

const (
	StateApproved         = "approved"
	StateRequestedChanges = "changes_requested"
	StateDeclined         = "declined"
)

type State string

type PR struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	State       string      `json:"state"`
	Author      Author      `json:"author"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	Description interface{} `json:"description"`
	Links       struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Commit Commit `json:"commit"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
	Reviewers    []Reviewer    `json:"reviewers"`
	Participants []Participant `json:"participants"`
}

type Commit struct {
	Hash  string `json:"hash"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Type string `json:"type"`
}

type Reviewer struct {
	DisplayName string `json:"display_name"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Type      string `json:"type"`
	UUID      string `json:"uuid"`
	AccountID string `json:"account_id"`
	Nickname  string `json:"nickname"`
}

type Participant struct {
	Type           string      `json:"type"`
	User           *User       `json:"user"`
	Role           string      `json:"role"`
	Approved       bool        `json:"approved"`
	State          State       `json:"state"`
	ParticipatedOn interface{} `json:"participated_on"`
}

type Activity struct {
	PullRequest      PR              `json:"pull_request"`
	Update           UpdateDetail    `json:"update,omitempty"`
	Approval         Approval        `json:"approval,omitempty"`
	ChangesRequested ChangeRequested `json:"changes_requested,omitempty"`
}

type Author struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
}

type UpdateDetail struct {
	State       string        `json:"state"`
	Draft       bool          `json:"draft"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Reviewers   []interface{} `json:"reviewers"`
	Changes     Changes       `json:"changes"`
	Reason      string        `json:"reason"`
	Author      User          `json:"author"`
	Date        string        `json:"date"`
	Destination BranchDetail  `json:"destination"`
	Source      BranchDetail  `json:"source"`
}

type Changes struct {
	Reviewers struct {
		Added []Reviewer `json:"added"`
	} `json:"reviewers"`
	Description struct {
		New string `json:"new"`
		Old string `json:"old"`
	} `json:"description"`
	Title struct {
		New string `json:"new"`
		Old string `json:"old"`
	} `json:"title"`
}

type ChangeRequested struct {
	Date        string `json:"date"`
	User        User   `json:"user"`
	PullRequest PR     `json:"pull_request"`
}

type User struct {
	DisplayName string `json:"display_name"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Type      string `json:"type"`
	UUID      string `json:"uuid"`
	AccountID string `json:"account_id"`
	Nickname  string `json:"nickname"`
}

type BranchDetail struct {
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
	Commit struct {
		Hash  string `json:"hash"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
		Type string `json:"type"`
	} `json:"commit"`
	Repository struct {
		Type     string `json:"type"`
		FullName string `json:"full_name"`
		Links    struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		Name string `json:"name"`
		UUID string `json:"uuid"`
	} `json:"repository"`
}

type Approval struct {
	Date        string `json:"date"`
	User        User   `json:"user"`
	PullRequest PR     `json:"pullrequest"`
}

type BitbucketPRResponse struct {
	Values []PR `json:"values"`
}

type BitbucketActivityResponse struct {
	Values []Activity `json:"values"`
}

type DiffstatResponse struct {
	Values  []DiffstatEntry `json:"values"`
	PageLen int             `json:"pagelen"`
	Size    int             `json:"size"`
	Page    int             `json:"page"`
}

type DiffstatEntry struct {
	Type         string    `json:"type"` // Typically "diffstat"
	LinesAdded   int       `json:"lines_added"`
	LinesRemoved int       `json:"lines_removed"`
	Status       string    `json:"status"` // e.g., "added", "removed", "modified"
	Old          *DiffFile `json:"old,omitempty"`
	New          *DiffFile `json:"new,omitempty"`
}

type DiffFile struct {
	Path        string `json:"path"`
	Type        string `json:"type"` // e.g., "commit_file"
	EscapedPath string `json:"escaped_path"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}
