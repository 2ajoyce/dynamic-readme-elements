package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"io"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func HandleVersion(c *gin.Context) {
	info, _ := debug.ReadBuildInfo()
	vcsRevision := getCommitHashFromSettings(info.Settings)
	var tags []string
	if vcsRevision != "" {
		slog.Info(fmt.Sprintf("vcs.revision=%s", vcsRevision))
		var err error
		tags, err = GetTagsForCommit(info.Main.Path, vcsRevision)
		if err != nil {
			slog.Error("Error Retrieving Version Tags", "error", err)
		}
	} else {
		slog.Error("No Commit Hash Found")
	}
	slog.Info(fmt.Sprintf("Tags:%+v", tags))
	c.JSON(http.StatusOK, gin.H{"tags": tags, "revision": vcsRevision})
}

// getCommitHashFromSettings takes a slice of debug.ModuleBuildInfo and iterates over it,
// looking for a key of "vcs.revision".
// When found, it returns the corresponding value as a string.
// If "vcs.revision" is not found in the settings, it returns an empty string.
func getCommitHashFromSettings(settings []debug.BuildSetting) string {
	var commitHash string
	for _, setting := range settings {
		if setting.Key == "vcs.revision" {
			commitHash = setting.Value
			break
		}
	}
	return commitHash
}

// getOwnerAndRepo takes a string path of format "github.com/owner/repo/version"
// and returns the "owner" and "repo" as string values.
// Returns an error if the path format is invalid.
func getOwnerAndRepo(path string) (string, string, error) {
	split := strings.Split(path, "/")
	if len(split) < 3 {
		return "", "", fmt.Errorf("invalid path format")
	}
	return split[1], split[2], nil
}

// GetTagsForCommit takes a string path (format "github.com/owner/repo/version")
// and a commit SHA string as arguments.
// It makes a GET request to https://github.com/{owner}/{repo}/branch_commits/{commitSha},
// parses the response body as HTML, and returns the tags associated with the commit.
// Returns an error if the HTTP request fails, or if the HTML parsing fails.
func GetTagsForCommit(path, commitSha string) ([]string, error) {
	owner, repo, err := getOwnerAndRepo(path)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/branch_commits/%s", owner, repo, commitSha))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close the response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	tags := extractTags(doc)
	return tags, nil
}

// extractTags takes an *html.Node and traverses the node tree,
// extracting the text content of <a> elements within the second <ul> element it encounters.
// Returns a slice of strings containing the extracted tags.
func extractTags(n *html.Node) []string {
	var ulCounter = 0
	var result []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			ulCounter++
			if ulCounter == 2 {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "li" {
						for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
							if cc.Type == html.ElementNode && cc.Data == "a" {
								result = append(result, cc.FirstChild.Data)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return result
}
