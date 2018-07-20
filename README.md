# Jira

## Install

The easiest way to install and use this tool is to download the binary from the [latest release](https://github.com/miketmoore/jira/releases/tag/v1.1.0).

## Configure

After downloading the binary, you will need to create a configuration file:

```
{
  "baseurl": "",
  "username": "",
  "apitoken": ""
}
```

Now, you can set the `$JIRACONFIG` environment variable to point to this config, or use the `-config` flag.

## Use

```
jira -issueid=YOUR_ISSUE_ID
```
