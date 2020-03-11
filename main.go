package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/c-bata/go-prompt"
)

var repos, repoMap, _ = githubGetUserRepos(os.Getenv("GITHUB_ACCESS_TOKEN"))

var promptState struct {
	subcommand string

	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return promptState.LivePrefix, promptState.IsEnable
}

func executor(in string) {
	if promptState.subcommand == "" {
		if in == "" {
			promptState.IsEnable = false
			promptState.LivePrefix = in
			return
		} else if in == "list" {
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "Name\tDescription\tLanguage\tURL\t# of Open Issues")
			for _, repo := range repos {
				fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v", *repo.Name, nilableString(repo.Description, 50), nilableString(repo.Language, 0), *repo.HTMLURL, *repo.OpenIssuesCount))
			}
			fmt.Fprintln(w)
			w.Flush()
		} else if in == "search" {
			promptState.LivePrefix = in + " > "
			promptState.IsEnable = true
			promptState.subcommand = "search"
		}
	} else if promptState.subcommand == "search" {

		if in == "" {
			promptState.IsEnable = false
			promptState.LivePrefix = in
			return
		} else if !validRepoName(in) {
			fmt.Println("Invalid repo name")
		} else {
			promptState.LivePrefix = "repo - " + in + " > "
			promptState.IsEnable = true
			promptState.subcommand = in
		}
	} else if validRepoName(promptState.subcommand) {
		if in == "" {
			promptState.IsEnable = false
			promptState.LivePrefix = in
			return
		} else if in == "info" {
			repo := repoMap[promptState.subcommand]
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 2, 8, 2, '\t', 0)
			fmt.Fprintln(w, "Name\tDescription\tLanguage\tURL\t# of Open Issues")
			fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v", *repo.Name, nilableString(repo.Description, 50), nilableString(repo.Language, 0), *repo.HTMLURL, *repo.OpenIssuesCount))
			fmt.Fprintln(w)
			w.Flush()
		} else if in == "open" {
			repo := repoMap[promptState.subcommand]
			openbrowser(*repo.HTMLURL)
		}
	}

}

func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	if promptState.subcommand == "" {
		s = []prompt.Suggest{{Text: "list", Description: "list all of your GitHub profiles"},
			{Text: "search", Description: "search for a specific GitHub repo"}}
	} else if promptState.subcommand == "search" {
		s = repoToSuggest()
	} else if validRepoName(promptState.subcommand) {
		s = []prompt.Suggest{{Text: "info", Description: "Get repository info"},
			{Text: "open", Description: "Open repository in default browser"}}
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("live-prefix-example"),
	)
	p.Run()
}

func repoToSuggest() []prompt.Suggest {
	var suggestions []prompt.Suggest
	for _, repo := range repos {
		item := prompt.Suggest{Text: *repo.Name, Description: nilableString(repo.Description, 0)}
		suggestions = append(suggestions, item)
	}
	return suggestions
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func nilableString(s *string, maxLength int) string {
	if s == nil {
		return ""
	} else if len(*s) > maxLength && maxLength != 0 {
		return strings.Replace(*s, `"`, `'`, -1)[:maxLength]
	}
	return strings.Replace(*s, `"`, `'`, -1)
}

func validRepoName(name string) bool {
	if _, ok := repoMap[name]; ok {
		return true
	}
	return false
}
