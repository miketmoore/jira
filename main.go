package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/andygrunwald/go-jira"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("Required args: base URL, username, issue ID")
		return
	}

	baseURL := args[0]
	username := args[1]
	id := args[2]

	issueURL := fmt.Sprintf("URL: %s/browse/%s", baseURL, id)

	fmt.Println(issueURL)

	fmt.Println("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Setting password failed")
		return
	}
	password := string(bytePassword)
	fmt.Println("Thanks! Sending request...")

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		fmt.Println("Connecting to Jira failed")
		fmt.Println(err)
		return
	}

	issue, _, err := jiraClient.Issue.Get(id, nil)
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
	fmt.Println()
	fmt.Printf("Summary: \t\t%s\n", issue.Fields.Summary)
	fmt.Printf("Status: \t\t%s\n", issue.Fields.Status.Name)
	fmt.Printf("Type: \t\t\t%s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: \t\t%s\n", issue.Fields.Priority.Name)
	fmt.Printf("Assignee: \t\t%s\n", issue.Fields.Assignee.DisplayName)
	fmt.Printf("Creator: \t\t%s\n", issue.Fields.Creator.DisplayName)

	fmt.Printf("Original Estimate: \t%s\n", getHumanReadableDuration(issue.Fields.TimeOriginalEstimate))
	fmt.Printf("Estimate: \t\t%s\n", getHumanReadableDuration(issue.Fields.TimeEstimate))
	fmt.Printf("Time Spent: \t\t%s\n", getHumanReadableDuration(issue.Fields.TimeSpent))

	if issue.Fields.Sprint != nil {
		fmt.Printf("Sprint ID: %d\n", issue.Fields.Sprint.ID)
		fmt.Printf("Sprint Name: %s\n", issue.Fields.Sprint.Name)
	}

	fmt.Printf("Total Comments: %d\n", len(issue.Fields.Comments.Comments))
	fmt.Println()
}
