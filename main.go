package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

// Bitbucket API details
const (
	BitbucketBaseURL   = "https://api.bitbucket.org/2.0"
	BitbucketUsername  = "Srijan"
	BitbucketRepoSlug  = "test_repo"      // e.g., "repo-name"
	BitbucketWorkspace = "chapssrijan619" // e.g., "team-name"
	AuthToken          = "ATCTT3xFfGN0yxNahao5bDKDs637xeRNHuwUKWw-etSBRMS8e-N_YdaBR3sZP4lvxHiOKu2oc3g4e1wbrvPBKjM_RTnUckIjX5J_4bZ_EVkkLbzaD1Ir6OWhUH_mZKOROQsJOAJUhnpkkSY3vrhdR4a7FvpxH29I4mjyHHMGyFCGd4l4snzIiqQ=76287C43"
)

// Struct to represent a Pull Request (PR)
type PR struct {
	Title  string `json:"title"`
	State  string `json:"state"`
	Author struct {
		DisplayName string `json:"display_name"`
		Username    string `json:"username"`
	} `json:"author"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	Description interface{} `json:"description"` // Change to interface{} to handle both string and object
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

type BitbucketPRResponse struct {
	Values []PR `json:"values"`
}

// Fetches PRs from Bitbucket
func fetchBitbucketPRs() []PR {
	client := resty.New()
	client.SetAuthToken(AuthToken)

	resp, err := client.R().
		SetResult(&BitbucketPRResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=ALL", BitbucketBaseURL, BitbucketWorkspace, BitbucketRepoSlug))
	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}

	return resp.Result().(*BitbucketPRResponse).Values
}

func main() {
	// Fetch PRs
	prs := fetchBitbucketPRs()

	// Create and run the UI
	app := CreateApp(prs)
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
