package sqlstore

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/focalboard/server/model"
)

func (s *SQLStore) hasUserPermissionTo(db sq.BaseRunner, userID string, permissionID string) (*model.User, error) {
	return s.getUserWithPermission(db, userID, permissionID)
}
