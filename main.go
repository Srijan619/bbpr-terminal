package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"simple-git-terminal/types"
	"simple-git-terminal/util"
)

// Bitbucket API details
const (
	BitbucketBaseURL   = "https://api.bitbucket.org/2.0"
	BitbucketUsername  = "Srijan"
	BitbucketRepoSlug  = "test_repo"
	BitbucketWorkspace = "chapssrijan619"
)

func getAuthToken() string {
	token := os.Getenv("BITBUCKET_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Environment variable BITBUCKET_AUTH_TOKEN is not set")
	}
	return token
}

// Fetches PRs from Bitbucket
func fetchBitbucketPRs() []types.PR {
	client := resty.New()
	client.SetAuthToken(getAuthToken())

	resp, err := client.R().
		SetResult(&types.BitbucketPRResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=ALL", BitbucketBaseURL, BitbucketWorkspace, BitbucketRepoSlug))
	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}

	prs := resp.Result().(*types.BitbucketPRResponse).Values

	for i := range prs {
		prs[i] = util.SanitizePR(prs[i])
	}

	return prs
}

func fetchBitbucketDiffStat(id int) string {
	client := resty.New()
	client.SetAuthToken(getAuthToken())

	// Fetching the diff for the given pull request ID
	resp, err := client.R().
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/diff", BitbucketBaseURL, BitbucketWorkspace, BitbucketRepoSlug, id))
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
func fetchBitbucketActivities(id int) []types.Activity {
	client := resty.New()
	client.SetAuthToken(getAuthToken())

	resp, err := client.R().
		SetResult(&types.BitbucketActivityResponse{}).
		Get(fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/activity", BitbucketBaseURL, BitbucketWorkspace, BitbucketRepoSlug, id))
	if err != nil {
		log.Fatalf("Error fetching activities: %v", err)
	}
	activityResponse := resp.Result().(*types.BitbucketActivityResponse)
	return activityResponse.Values
}

func main() {
	// Open or create the log file
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)
	// Optionally add log flags (e.g., timestamp)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Log a test message to verify
	log.Printf("Application started")

	prs := fetchBitbucketPRs()

	app := CreateApp(prs)
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	log.Printf("Application ended")
}
