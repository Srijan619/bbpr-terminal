package bitbucket

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"simple-git-terminal/state"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
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

func FetchBitbucketPRs() []types.PR {
	// Create the client (authentication handled inside)
	client := createClient()

	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=ALL", BitbucketBaseURL, state.Workspace, state.Repo)
	log.Printf("URL:::%s", url)
	resp, err := client.R().
		SetResult(&types.BitbucketPRResponse{}).
		Get(url)

	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}
	// Check if the response is successful (status code 200)
	if resp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode(), string(resp.Body()))
	}
	prs := resp.Result().(*types.BitbucketPRResponse).Values

	// Process and sanitize the PRs
	for i := range prs {
		prs[i] = util.SanitizePR(prs[i])
	}

	return prs
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
