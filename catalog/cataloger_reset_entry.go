package catalog

import (
	"context"

	"github.com/treeverse/lakefs/db"
)

func (c *cataloger) ResetEntry(ctx context.Context, repository, branch string, path string) error {
	if err := Validate(ValidateFields{
		{Name: "repository", IsValid: ValidateRepositoryName(repository)},
		{Name: "branch", IsValid: ValidateBranchName(branch)},
		{Name: "path", IsValid: ValidatePath(path)},
	}); err != nil {
		return err
	}
	_, err := c.db.Transact(func(tx db.Tx) (interface{}, error) {
		branchID, err := getBranchID(tx, repository, branch, LockTypeShare)
		if err != nil {
			return nil, err
		}
		res, err := tx.Exec(`DELETE FROM entries WHERE branch_id = $1 AND path = $2 AND min_commit = 0`, branchID, path)
		if err != nil {
			return nil, err
		}
		if affected, err := res.RowsAffected(); err != nil {
			return nil, err
		} else if affected != 1 {
			return nil, ErrEntryNotFound
		}
		return nil, nil
	}, c.txOpts(ctx)...)
	return err
}