package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/olekukonko/tablewriter"
)

type jsonConfig struct {
	BaseURL  string `json:"baseurl"`
	Username string `json:"username"`
	APIToken string `json:"apitoken"`
}

const configEnvName = "JIRACONFIG"

func main() {

	issueID := flag.String("issueid", "", "issue id")
	configpathflag := flag.String("config", "", "Jira API Config (optional)")
	configpathenv := os.Getenv(configEnvName)
	flag.Parse()

	var configpath string
	if configpathenv == "" {
		if configpathflag == nil || *configpathflag == "" {
			// neither throw error
			fmt.Println("You must specify your configuration file path.")
			fmt.Println("You can do this either by setting the $JIRACONFIG environment variable, or by using the -config flag.")
			os.Exit(1)
		}
		configpath = *configpathflag
	} else {
		configpath = configpathenv
	}

	flag.Parse()

	if issueID == nil || *issueID == "" {
		fmt.Println("-issueid flag is required")
		os.Exit(1)
	}

	var config jsonConfig

	file, err := ioutil.ReadFile(configpath)
	if err != nil {
		fmt.Println("Failed to load config")
		fmt.Println(err)
		os.Exit(1)
	}

	json.Unmarshal(file, &config)

	issueURL := fmt.Sprintf("URL: %s/browse/%s", config.BaseURL, *issueID)

	fmt.Println(issueURL)

	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.APIToken,
	}

	jiraClient, err := jira.NewClient(tp.Client(), config.BaseURL)
	if err != nil {
		fmt.Println("Connecting to Jira failed")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(*issueID)
	issue, _, err := jiraClient.Issue.Get(*issueID, nil)
	if err != nil {
		fmt.Println("Failed to get issue")
		fmt.Println(err)
		os.Exit(1)
	}

	printIssueDetails(issue)
}

func getHumanReadableDuration(seconds int) string {
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 8

	final := []string{}

	if days > 0 {
		final = append(final, fmt.Sprintf("%dd", days))
	} else if hours%60 > 0 {
		final = append(final, fmt.Sprintf("%dhr", hours%60))
	}

	if minutes%60 > 0 {
		final = append(final, fmt.Sprintf("%dm", minutes%60))
	}

	if seconds%60 > 0 {
		final = append(final, fmt.Sprintf("%ds", seconds%60))
	}

	return strings.Join(final, " ")
}

func printIssueDetails(issue *jira.Issue) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})

	data := [][]string{
		[]string{"Summary", issue.Fields.Summary},
		[]string{"Status", issue.Fields.Status.Name},
		[]string{"Type", issue.Fields.Type.Name},
		[]string{"Priority", issue.Fields.Priority.Name},
		[]string{"Assignee", issue.Fields.Assignee.DisplayName},
		[]string{"Creator", issue.Fields.Creator.DisplayName},
		[]string{"Original Estimate", getHumanReadableDuration(issue.Fields.TimeOriginalEstimate)},
		[]string{"Estimate", getHumanReadableDuration(issue.Fields.TimeEstimate)},
		[]string{"Time Spent", getHumanReadableDuration(issue.Fields.TimeSpent)},
		[]string{"Total Comments", fmt.Sprintf("%d", len(issue.Fields.Comments.Comments))},
	}

	table.AppendBulk(data)

	// Customize output for markdown
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
}
