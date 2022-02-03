package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// *****SET THESE PARAMETERS*****
const ORG_NAME string = ""
const PAT string = ""
const NAME string = ""
const EMAIL string = ""

//Adds the Google Scorecard workflow to all repositores under the given organization
func main() {
	//Get github user client
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PAT},
	)
	tokenClient := oauth2.NewClient(context, tokenService)
	client := github.NewClient(tokenClient)

	//Get repositories under organization
	lops := &github.RepositoryListByOrgOptions{Type: "all"}
	repos, _, err := client.Repositories.ListByOrg(context, ORG_NAME, lops)
	err_check(err, "List Org Repos Error")

	//Convert to list of repository names
	repoNames := []string{}
	for _, repo := range repos {
		repoNames = append(repoNames, *repo.Name)
	}

	fmt.Println(repoNames)
	fileContent, _ := ioutil.ReadFile("scorecards.yml")

	for _, repoName := range repoNames {
		//TODO: SKIP REPO IF WORKFLOW ALREADY EXISTS
		//Add .yml workflow file to the repository
		opts := &github.RepositoryContentFileOptions{
			Message:   github.String("Adding workflow"),
			Content:   fileContent,
			Branch:    github.String("main"),
			Committer: &github.CommitAuthor{Name: github.String(NAME), Email: github.String(EMAIL)},
		}
		_, _, err = client.Repositories.CreateFile(context, ORG_NAME, repoName, ".github/workflows/scorecards.yml", opts)
		err_check(err, "CreateFile Error")

		//Wait for workflow file to finish creating
		time.Sleep(time.Second)

		//Trigger the workflow
		ref := github.CreateWorkflowDispatchEventRequest{
			Ref: "main",
		}
		_, err = client.Actions.CreateWorkflowDispatchEventByFileName(context, ORG_NAME, "workflow-test", "scorecards.yml", ref)
		err_check(err, "Trigger Workflow Error")
	}
}

func err_check(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}
