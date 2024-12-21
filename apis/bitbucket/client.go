package bitbucket

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/url"
	"os"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"strings"
)

// Bitbucket API details
const (
	BitbucketBaseURL                = "https://api.bitbucket.org/2.0"
	BitbucketEnvTokenName           = "BITBUCKET_AUTH_TOKEN"
	BitbucketEnvAppPasswordName     = "BITBUCKET_APP_PASSWORD"
	BitbucketEnvAppPasswordUsername = "BITBUCKET_APP_USERNAME"
)

func getAuthToken(tokenString string) string {
	token := os.Getenv(tokenString)
	if token == "" {
		log.Printf("Environment variable %s is not set will try using basic authentication with app password", tokenString)
	}
	return token
}

// Helper function to create a Resty client with authentication
func createClient() *resty.Client {
	client := resty.New()

	authToken := getAuthToken(BitbucketEnvTokenName)
	if authToken != "" {
		client.SetAuthToken(authToken)
	} else {
		username := os.Getenv(BitbucketEnvAppPasswordUsername)
		appPassword := os.Getenv(BitbucketEnvAppPasswordName)

		if username != "" && appPassword != "" {
			client.SetBasicAuth(username, appPassword)
		} else {
			log.Fatalf("Error: Missing authentication credentials. Please check your environment variables.")
		}
	}

	return client
}

func FetchPR(id int) *types.PR {
	client := createClient()
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d", BitbucketBaseURL, state.Workspace, state.Repo, id)

	resp, err := client.R().
		SetResult(&types.PR{}).
		Get(url)

	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode(), string(resp.Body()))
	}

	pr := resp.Result().(*types.PR)
	return pr
}

// Make query using BuildQuery method....
func FetchPRsByQuery(query string) []types.PR {
	client := createClient()
	encodedQuery := url.QueryEscape(query) // This will properly encode the query string
	fields := url.QueryEscape("+values.participants,-values.description,-values.summary")
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?pagelen=25&fields=%s&q=%s&page=1",
		BitbucketBaseURL, state.Workspace, state.Repo, fields, encodedQuery)
	url = strings.ReplaceAll(url, "+", "%20") // TODO: Some weird encoding issue..

	log.Printf("Fetching quer...%v", url)
	resp, err := client.R().
		SetResult(&types.BitbucketPRResponse{}).
		Get(url)

	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
		return nil
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode(), string(resp.Body()))
		return nil
	}

	prs := resp.Result().(*types.BitbucketPRResponse).Values
	// for i := range prs {
	// 	prs[i] = util.SanitizePR(prs[i])
	// }

	return prs
}

// TODO: This can be done better as there is an internal filter system that bitbucket provides...check search method
func FetchPRsByState(prState string) []types.PR {
	client := createClient()
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=%s", BitbucketBaseURL, state.Workspace, state.Repo, prState)
	log.Printf("Fetching PRs with state: %s", prState)

	resp, err := client.R().
		SetResult(&types.BitbucketPRResponse{}).
		Get(url)

	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode(), string(resp.Body()))
	}

	prs := resp.Result().(*types.BitbucketPRResponse).Values
	// for i := range prs {
	// 	prs[i] = util.SanitizePR(prs[i])
	// }

	return prs
}

func FetchBitbucketPRs() []types.PR {
	return FetchPRsByState("ALL")
}

func FetchBitbucketOpenPRs() []types.PR {
	return FetchPRsByState("OPEN")
}

func FetchBitbucketMergedPRs() []types.PR {
	return FetchPRsByState("MERGED")
}

func FetchBitbucketDeclinedPRs() []types.PR {
	return FetchPRsByState("DECLINED")
}

func FetchBitbucketDiffContent(id int, filePath string) (string, error) {
	client := createClient()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/diff?path=%s",
			BitbucketBaseURL,
			state.Workspace,
			state.Repo,
			id,
			filePath,
		))
	if err != nil {
		return "", fmt.Errorf("error fetching diff content: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	return string(resp.Body()), nil
}

// TODO: Same here maybe this endpoint should be made optional for user and just do local diff for faster diff?
func FetchBitbucketDiffstat(id int) []types.DiffstatEntry {
	client := createClient()

	// Fetching the diffstat for the given pull request ID
	resp, err := client.R().
		SetResult(&types.DiffstatResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/diffstat", BitbucketBaseURL, state.Workspace, state.Repo, id))
	if err != nil {
		log.Fatalf("Error fetching diffstat: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Fatalf("Error: Unexpected status code %d", resp.StatusCode())
	}

	response := resp.Result().(*types.DiffstatResponse)
	return response.Values
}

// TODO: Maybe this endpoint should be able optional for end user if they want to use network? It is pretty slow
func FetchBitbucketDiff(id int) string {
	client := createClient()

	// Fetching the diff for the given pull request ID
	resp, err := client.R().
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/diff", BitbucketBaseURL, state.Workspace, state.Repo, id))
	if err != nil {
		log.Fatalf("Error fetching diffstat: %v", err)
	}

	// Check if the response is successful (e.g., status code 200)
	if resp.StatusCode() != 200 {
		log.Fatalf("Error: Unexpected status code %d", resp.StatusCode())
	}

	// Return the raw diff content (response body is the diff)
	return string(resp.Body())
}

// Fetches recent activities from Bitbucket
func FetchBitbucketActivities(id int) []types.Activity {
	client := createClient()

	resp, err := client.R().
		SetResult(&types.BitbucketActivityResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/activity", BitbucketBaseURL, state.Workspace, state.Repo, id))
	if err != nil {
		log.Fatalf("Error fetching activities: %v", err)
	}
	activityResponse := resp.Result().(*types.BitbucketActivityResponse)
	return activityResponse.Values
}

func FetchBitbucketComments(id int) []types.Comment {
	client := createClient()

	resp, err := client.R().
		SetResult(&types.BitbucketCommentsResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/comments", BitbucketBaseURL, state.Workspace, state.Repo, id))
	if err != nil {
		log.Fatalf("Error fetching comments: %v", err)
	}
	response := resp.Result().(*types.BitbucketCommentsResponse)
	return response.Values
}

func UpdateFilteredPRs() {
	var filteredPRs []types.PR
	// Fetch or use cached PRs based on active filters
	if state.PRStatusFilter.Open {
		log.Printf("[sausage] Appending open ones...")
		filteredPRs = append(filteredPRs, FetchPRsByState("OPEN")...)
	}
	if state.PRStatusFilter.Merged {
		log.Printf("[sausage] Appending merged ones...")
		filteredPRs = append(filteredPRs, FetchPRsByState("MERGED")...)
	}
	if state.PRStatusFilter.Declined {
		filteredPRs = append(filteredPRs, FetchPRsByState("DECLINED")...)
	}
	state.SetFilteredPRs(&filteredPRs)
}

func BuildQuery(searchTerm string) string {
	// Add state filters (Open, Merged, Declined)
	stateFilters := []string{}
	filters := state.PRStatusFilter
	if filters.Merged {
		stateFilters = append(stateFilters, "state=\"MERGED\"")
	}
	if filters.Declined {
		stateFilters = append(stateFilters, "state=\"DECLINED\"")
	}
	if filters.Open {
		stateFilters = append(stateFilters, "state=\"OPEN\"")
	}

	// Combine state filters with OR if any exist
	var stateQuery string
	if len(stateFilters) > 0 {
		stateQuery = strings.Join(stateFilters, " OR ")
	}

	// Add search term conditions (title or description contains search term)
	var searchQuery string
	if searchTerm != "" {
		searchQuery = "(description~\"" + searchTerm + "\" OR title~\"" + searchTerm + "\")"
	}

	// Combine the state and search queries with AND
	var finalQuery string
	if stateQuery != "" && searchQuery != "" {
		// Combine state query and search query with AND
		finalQuery = stateQuery + " AND " + searchQuery
	} else if stateQuery != "" {
		// If only the state query exists
		finalQuery = stateQuery
	} else if searchQuery != "" {
		// If only the search query exists
		finalQuery = searchQuery
	}

	return finalQuery
}
