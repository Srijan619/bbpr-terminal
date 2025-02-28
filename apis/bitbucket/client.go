package bitbucket

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Bitbucket API details
const (
	BitbucketBaseURL                = "https://api.bitbucket.org/2.0"
	BitbucketEnvTokenName           = "BITBUCKET_AUTH_TOKEN"
	BitbucketEnvAppPasswordName     = "BITBUCKET_APP_PASSWORD"
	BitbucketEnvAppPasswordUsername = "BITBUCKET_APP_USERNAME"
)

var client *resty.Client

func getAuthToken(tokenString string) string {
	token := os.Getenv(tokenString)
	if token == "" {
		log.Printf("[CLIENT] Environment variable %s is not set will try using basic authentication with app password", tokenString)
	}
	return token
}

// Helper function to create a Resty client with authentication
func createClient() *resty.Client {
	if client != nil {
		return client
	}

	client = resty.New()

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
func FetchPRsByQuery(query string) ([]types.PR, types.Pagination) {
	client := createClient()
	encodedQuery := url.QueryEscape(query) // Properly encode the query string
	fields := url.QueryEscape("+values.participants,-values.description,-values.summary")

	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?fields=%s&q=%s&page=%d",
		BitbucketBaseURL, state.Workspace, state.Repo, fields, encodedQuery, state.Pagination.Page)
	url = strings.ReplaceAll(url, "+", "%20") // TODO: Fix encoding issue

	log.Printf("[CLIENT] Fetching PRs with query...%v", url)
	resp, err := client.R().
		SetResult(&types.BitbucketPRResponse{}).
		Get(url)
	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
		return nil, types.Pagination{}
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode(), string(resp.Body()))
		return nil, types.Pagination{}
	}

	response := resp.Result().(*types.BitbucketPRResponse)
	return response.Values, response.Pagination
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

func FetchCurrentUser() *types.User {
	client := createClient()

	resp, err := client.R().
		SetResult(&types.User{}).
		Get(fmt.Sprintf("%s/user", BitbucketBaseURL))
	if err != nil {
		log.Fatalf("Error fetching user: %v", err)
	}
	userResponse := resp.Result().(*types.User)
	if userResponse == nil {
		log.Printf("No active user, probably using API token which does not give access to current active user..")
	} else {
		log.Printf("Current active user => %v", userResponse)
	}

	return userResponse
}

func BuildQuery(searchTerm string) string {
	var filters []string

	stateFilter := buildStateFilter()
	if stateFilter != "" {
		filters = append(filters, stateFilter)
	}

	// Add author filter if IAmReviewing is true
	if state.CurrentUser.UUID != "" && state.PRStatusFilter.IAmAuthor {
		authorFilter := fmt.Sprintf("author.uuid=\"%s\"", state.CurrentUser.UUID)
		filters = append(filters, authorFilter)
	}

	if state.CurrentUser.UUID != "" && state.PRStatusFilter.IAmReviewer {
		reviewersFilter := fmt.Sprintf("reviewers.uuid=\"%s\"", state.CurrentUser.UUID)
		filters = append(filters, reviewersFilter)
	}
	// Add search term filter
	if searchTerm != "" {
		searchFilter := fmt.Sprintf("(description~\"%s\" OR title~\"%s\")", searchTerm, searchTerm)
		filters = append(filters, searchFilter)
	}

	// Combine all filters with AND
	finalQuery := strings.Join(filters, " AND ")

	log.Printf("Final built query => %s", finalQuery)
	return finalQuery
}

func buildStateFilter() string {
	// Initialize state filters array
	var stateFilters []string

	// Add individual state filters (Open, Merged, Declined)
	if state.PRStatusFilter.Merged {
		stateFilters = append(stateFilters, "state=\"MERGED\"")
	}
	if state.PRStatusFilter.Declined {
		stateFilters = append(stateFilters, "state=\"DECLINED\"")
	}
	if state.PRStatusFilter.Open {
		stateFilters = append(stateFilters, "state=\"OPEN\"")
	}

	// Combine the state filters into a single string with OR
	if len(stateFilters) > 0 {
		return strings.Join(stateFilters, " OR ")
	}

	return ""
}
