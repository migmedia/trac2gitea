// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/import/issue"
	"stevejefferson.co.uk/trac2gitea/import/wiki"
	"stevejefferson.co.uk/trac2gitea/log"
)

var dbOnly bool
var wikiOnly bool
var wikiPush bool
var verbose bool
var wikiConvertPredefineds bool
var tracRootDir string
var giteaRootDir string
var giteaUser string
var giteaRepo string
var giteaWikiRepoURL string
var giteaWikiRepoToken string
var giteaWikiRepoDir string
var giteaDefaultAssignee string
var giteaDefaultAuthor string
var giteaDefaultWikiAuthor string

func parseArgs() {
	defaultAssigneeParam := pflag.String("default-assignee", "",
		"`username` to assign tickets to when trac assignee is not found in Gitea - defaults to <gitea-user>")
	defaultAuthorParam := pflag.String("default-author", "",
		"`username` to attribute content to when trac author is not found in Gitea - defaults to <gitea-user>")
	defaultWikiAuthorParam := pflag.String("default-wiki-author", "",
		"`username` to attribute Wiki content to when trac author is not found in Gitea - defaults to <gitea-user>")

	wikiURLParam := pflag.String("wiki-url", "",
		"URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git")
	wikiTokenParam := pflag.String("wiki-token", "",
		"password/token for accessing wiki repository (ignored if wiki-url provided)")
	wikiDirParam := pflag.String("wiki-dir", "",
		"directory into which to checkout (clone) wiki repository - defaults to cwd")
	wikiConvertPredefinedsParam := pflag.Bool("wiki-convert-predefined", false,
		"convert Trac predefined wiki pages - by default we skip these")

	dbOnlyParam := pflag.Bool("db-only", false,
		"convert database only")
	wikiOnlyParam := pflag.Bool("wiki-only", false,
		"convert wiki only")
	wikiNoPushParam := pflag.Bool("no-wiki-push", false,
		"do not push wiki on completion")
	verboseParam := pflag.Bool("verbose", false,
		"verbose output")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo>\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	verbose = *verboseParam
	dbOnly = *dbOnlyParam
	wikiOnly = *wikiOnlyParam
	wikiPush = !*wikiNoPushParam

	if dbOnly && wikiOnly {
		log.Fatal("Cannot generate only database AND only wiki!")
	}
	wikiConvertPredefineds = *wikiConvertPredefinedsParam

	if pflag.NArg() < 4 {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUser = pflag.Arg(2)
	giteaRepo = pflag.Arg(3)

	giteaDefaultAssignee = *defaultAssigneeParam
	if giteaDefaultAssignee == "" {
		giteaDefaultAssignee = giteaUser
	}
	giteaDefaultAuthor = *defaultAuthorParam
	if giteaDefaultAuthor == "" {
		giteaDefaultAuthor = giteaUser
	}
	giteaDefaultWikiAuthor = *defaultWikiAuthorParam
	if giteaDefaultWikiAuthor == "" {
		giteaDefaultWikiAuthor = giteaUser
	}
	giteaWikiRepoURL = *wikiURLParam
	giteaWikiRepoToken = *wikiTokenParam
	giteaWikiRepoDir = *wikiDirParam
}

func main() {
	parseArgs()

	var logLevel = log.INFO
	if verbose {
		logLevel = log.TRACE
	}
	log.SetLevel(logLevel)

	tracAccessor, err := trac.CreateDefaultAccessor(tracRootDir)
	if err != nil {
		log.Fatal("%v\n", err)
	}
	giteaAccessor, err := gitea.CreateDefaultAccessor(
		giteaRootDir, giteaUser, giteaRepo, giteaWikiRepoURL, giteaWikiRepoToken, giteaWikiRepoDir, giteaDefaultAssignee, giteaDefaultAuthor)
	if err != nil {
		log.Fatal("%v\n", err)
	}

	if !wikiOnly {
		issueImporter, err := issue.CreateImporter(tracAccessor, giteaAccessor)
		if err != nil {
			log.Fatal("%v\n", err)
		}

		issueImporter.ImportComponents()
		issueImporter.ImportPriorities()
		issueImporter.ImportSeverities()
		issueImporter.ImportVersions()
		issueImporter.ImportTypes()
		issueImporter.ImportResolutions()
		issueImporter.ImportMilestones()
		issueImporter.ImportTickets()
	}

	if !dbOnly {
		wikiImporter, err := wiki.CreateImporter(tracAccessor, giteaAccessor, giteaDefaultWikiAuthor, wikiConvertPredefineds)
		if err != nil {
			log.Fatal("%v\n", err)
		}

		wikiImporter.ImportWiki(wikiPush)
	}
}
