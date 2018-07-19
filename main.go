package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/andygrunwald/go-jira"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	baseURL := flag.String("baseurl", "", "a string")
	username := flag.String("username", "", "a string")
	token := flag.String("token", "", "a string")
	id := flag.String("id", "", "a string")

	flag.Parse()

	var password string
	if token == nil || *token == "" {
		fmt.Println("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Setting password failed")
			return
		}
		password = string(bytePassword)
	}

	issueURL := fmt.Sprintf("URL: %s/browse/%s", *baseURL, *id)

	fmt.Println(issueURL)

	fmt.Println("Thanks! Sending request...")

	var tp jira.BasicAuthTransport
	if password != "" {
		tp = jira.BasicAuthTransport{
			Username: *username,
			Password: password,
		}
	} else {
		tp = jira.BasicAuthTransport{
			Username: *username,
			Password: *token,
		}
	}

	jiraClient, err := jira.NewClient(tp.Client(), *baseURL)
	if err != nil {
		fmt.Println("Connecting to Jira failed")
		fmt.Println(err)
		return
	}

	issue, _, err := jiraClient.Issue.Get(*id, nil)
	if err != nil {
		fmt.Println("Failed to get issue")
		fmt.Println(err)
		return
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
