// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// GetUserID retrieves the id of a named Gitea user - returns NullID if no such user.
func (accessor *DefaultAccessor) GetUserID(userName string) (int64, error) {
	if strings.Trim(userName, " ") == "" {
		return NullID, nil
	}

	var id int64 = NullID
	err := accessor.db.QueryRow(`SELECT id FROM user WHERE lower_name = ? or email = ?`, userName, userName).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of user %s", userName)
		return NullID, err
	}

	return id, nil
}

// GetUserEMailAddress retrieves the email address of a given user
func (accessor *DefaultAccessor) GetUserEMailAddress(userName string) (string, error) {
	var emailAddress string = ""
	err := accessor.db.QueryRow(`SELECT email FROM user WHERE lower_name = ?`, userName).Scan(&emailAddress)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving email address of user %s", userName)
		return "", err
	}

	return emailAddress, nil
}

// getUserRepoURL retrieves the URL of the current repository for the current user
func (accessor *DefaultAccessor) getUserRepoURL() string {
	rootURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", rootURL, accessor.userName, accessor.repoName)
}

// MatchUser retrieves the name of the user best matching a user name or email address
func (accessor *DefaultAccessor) MatchUser(userName string, userEmail string) (string, error) {
	var matchedUserName = ""
	lcUserName := strings.ToLower(userName)
	err := accessor.db.QueryRow(`
		SELECT lower_name FROM user 
		WHERE lower_name = ?
		OR full_name = ?
		OR email = ?`, lcUserName, userName, userEmail).Scan(&matchedUserName)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "trying to match user name %s, email %s", userName, userEmail)
		return "", err
	}

	return matchedUserName, nil
}
