package gitea

import "log"

// AddComment adds a comment to Gitea
func (accessor *Accessor) AddComment(issueID int64, authorID int64, comment string, time int64) int64 {
	_, err := accessor.db.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, content, created_unix, updated_unix)
			VALUES ( 0, $1, $2, $3, $4, $4 )`,
		issueID, authorID, comment, time)
	if err != nil {
		log.Fatal(err)
	}

	var commentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&commentID)
	if err != nil {
		log.Fatal(err)
	}

	return commentID
}
