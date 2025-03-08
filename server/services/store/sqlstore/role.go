package sqlstore

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (s *SQLStore) rolesFromRows(rows *sql.Rows) ([]*model.Role, error) {
	roles := []*model.Role{}

	for rows.Next() {
		var role model.Role
		var permissions []byte

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&permissions,
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(permissions, &role.Permissions)
		if err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	return roles, nil
}

func (s *SQLStore) getRolesByCondition(db sq.BaseRunner, condition interface{}) ([]*model.Role, error) {
	query := s.getQueryBuilder(db).
		Select(
			fmt.Sprintf("%sroles.id", s.tablePrefix),
			fmt.Sprintf("%sroles.name", s.tablePrefix),
			fmt.Sprintf("json_agg(%spermissions.name) as permissions", s.tablePrefix),
		).
		From(s.tablePrefix + "roles").
		Join(fmt.Sprintf("%srole_permissions ON %srole_permissions.role_id = %sroles.id", s.tablePrefix, s.tablePrefix, s.tablePrefix)).
		Join(fmt.Sprintf("%spermissions ON %spermissions.id = %srole_permissions.permission_id", s.tablePrefix, s.tablePrefix, s.tablePrefix))

	if condition != nil {
		query = query.Where(condition)
	}
	query = query.GroupBy(
		fmt.Sprintf("%sroles.id", s.tablePrefix),
	)

	rows, err := query.Query()
	if err != nil {
		s.logger.Error(`getRoles ERROR`, mlog.Err(err))
		return nil, err
	}
	defer s.CloseRows(rows)

	roles, err := s.rolesFromRows(rows)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, model.NewErrNotFound("role")
	}

	return roles, nil
}

func (s *SQLStore) getRoles(db sq.BaseRunner) ([]*model.Role, error) {
	return s.getRolesByCondition(db, nil)
}

func (s *SQLStore) getRoleByName(db sq.BaseRunner, name string) (*model.Role, error) {
	roles, err := s.getRolesByCondition(db, sq.Eq{fmt.Sprintf("%sroles.name", s.tablePrefix): name})
	if err != nil {
		return nil, err
	}
	if len(roles) < 1 {
		return nil, model.NewErrNotFound("role")
	}
	return roles[0], err
}
