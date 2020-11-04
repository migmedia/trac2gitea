// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// getIssueParticipantID retrieves the id of the given issue participant, returns NullID if no such participant
func (accessor *DefaultAccessor) getIssueParticipantID(issueID int64, userID int64) (int64, error) {
	var issueParticipantID = NullID
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_user WHERE issue_id = ? AND uid = ?
		`, issueID, userID).Scan(&issueParticipantID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id for participant %d in issue %d", userID, issueID)
		return NullID, err
	}

	return issueParticipantID, nil
}

// updateIssueParticipant updates an existing issue participant
func (accessor *DefaultAccessor) updateIssueParticipant(issueParticipantID int64, issueID int64, userID int64) error {
	_, err := accessor.db.Exec(`UPDATE issue_user SET issue_id=?, uid=? WHERE id=?`,
		issueID, userID, issueParticipantID)
	if err != nil {
		err = errors.Wrapf(err, "updating participant %d in issue %d", userID, issueID)
		return err
	}

	log.Debug("updated participant %d in issue %d (id %d)", userID, issueID, issueParticipantID)

	return nil
}

// insertIssueParticipant creates a new issue participant
func (accessor *DefaultAccessor) insertIssueParticipant(issueID int64, userID int64) error {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_user(issue_id, uid, is_read, is_mentioned) VALUES (?, ?, 1, 0)`,
		issueID, userID)
	if err != nil {
		err = errors.Wrapf(err, "adding participant %d in issue %d", userID, issueID)
		return err
	}

	log.Debug("added participant %d in issue %d", userID, issueID)

	return nil
}

// AddIssueParticipant adds a participant to a Gitea issue.
func (accessor *DefaultAccessor) AddIssueParticipant(issueID int64, userID int64) error {
	issueParticipantID, err := accessor.getIssueParticipantID(issueID, userID)
	if err != nil {
		return err
	}

	if issueParticipantID == NullID {
		return accessor.insertIssueParticipant(issueID, userID)
	}

	if accessor.overwrite {
		err = accessor.updateIssueParticipant(issueParticipantID, issueID, userID)
		if err != nil {
			return err
		}
	} else {
		log.Debug("issue %d already has participant %d - ignored", issueID, userID)
	}

	return nil
}
