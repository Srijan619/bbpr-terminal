package types

type PR struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	State       string      `json:"state"`
	Author      Author      `json:"author"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	Description interface{} `json:"description"` // To handle both string and object
	Links       struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
}

type Activity struct {
	PullRequest PR           `json:"pull_request"`
	Update      UpdateDetail `json:"update,omitempty"`
	Approval    Approval     `json:"approval,omitempty"`
}

type Author struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
}

type UpdateDetail struct {
	State       string            `json:"state"`
	Draft       bool              `json:"draft"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Reviewers   []interface{}     `json:"reviewers"`
	Changes     map[string]Change `json:"changes"`
	Reason      string            `json:"reason"`
	Author      AuthorDetail      `json:"author"`
	Date        string            `json:"date"`
	Destination BranchDetail      `json:"destination"`
	Source      BranchDetail      `json:"source"`
}

type Change struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type AuthorDetail struct {
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
	Date        string       `json:"date"`
	User        AuthorDetail `json:"user"`
	PullRequest PR           `json:"pullrequest"`
}

type BitbucketPRResponse struct {
	Values []PR `json:"values"`
}

type BitbucketActivityResponse struct {
	Values []Activity `json:"values"`
}
