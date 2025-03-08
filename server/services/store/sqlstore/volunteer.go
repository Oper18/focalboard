package sqlstore

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (s *SQLStore) volunteersFromRows(rows *sql.Rows) ([]*model.Volunteer, error) {
	volunteers := []*model.Volunteer{}

	for rows.Next() {
		var volunteer model.Volunteer
		var files pq.StringArray

		err := rows.Scan(
			&volunteer.ID,
			&volunteer.Name,
			&volunteer.Contact,
			&files,
			&volunteer.Description,
			&volunteer.CreateAt,
			&volunteer.UpdateAt,
			&volunteer.DeleteAt,
		)
		if err != nil {
			return nil, err
		}
		volunteer.Files = []string(files)

		volunteers = append(volunteers, &volunteer)
	}

	return volunteers, nil
}

func (s *SQLStore) createVolunteer(db sq.BaseRunner, volunteer *model.Volunteer) (*model.Volunteer, error) {
	now := utils.GetMillis()
	volunteer.CreateAt = now
	volunteer.UpdateAt = now
	volunteer.DeleteAt = 0

	var filesArray pq.StringArray
	if volunteer.Files != nil {
		filesArray = pq.StringArray(volunteer.Files)
	} else {
		filesArray = pq.StringArray{}
	}

	query := s.getQueryBuilder(db).Insert(s.tablePrefix+"volunteers").
		Columns("name", "contact", "files", "description", "create_at", "update_at", "delete_at").
		Values(volunteer.Name, volunteer.Contact, filesArray, volunteer.Description, volunteer.CreateAt, volunteer.UpdateAt, volunteer.DeleteAt).
		Suffix("RETURNING id")

	var lastInsertID int64
	err := query.QueryRow().Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	volunteer.ID = lastInsertID
	return volunteer, err
}

func (s *SQLStore) updateVolunteer(db sq.BaseRunner, volunteer *model.Volunteer) (*model.Volunteer, error) {
	now := utils.GetMillis()
	volunteer.UpdateAt = now

	query := s.getQueryBuilder(db).Update(s.tablePrefix+"volunteers").
		Set("name", volunteer.Name).
		Set("contact", volunteer.Contact).
		Set("description", volunteer.Description).
		Set("update_at", volunteer.UpdateAt)

	if volunteer.Files != nil {
		filesArray := pq.StringArray(volunteer.Files)
		query = query.Set("files", filesArray)
	}
	query = query.Where(sq.Eq{"id": volunteer.ID})

	result, err := query.Exec()
	if err != nil {
		return nil, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowCount < 1 {
		return nil, model.NewErrNotFound("volunteer")
	}

	return volunteer, nil
}

func (s *SQLStore) getVolunteer(db sq.BaseRunner, volunteerID int64) (*model.Volunteer, error) {
	volunteers, err := s.getVolunteersByCondition(db, sq.Eq{fmt.Sprintf("%svolunteers.id", s.tablePrefix): volunteerID}, 1, 0)
	if err != nil {
		return nil, err
	}

	if len(volunteers) == 0 {
		return nil, model.NewErrNotFound("volunteer")
	}

	return volunteers[0], nil
}

func (s *SQLStore) getVolunteersList(db sq.BaseRunner, contact string, limit, offset uint64) ([]*model.Volunteer, error) {
	var condition interface{}
	if contact != "" {
		condition = sq.Eq{fmt.Sprintf("%svolunteers.contact", s.tablePrefix): contact}
	}
	return s.getVolunteersByCondition(db, condition, limit, offset)
}

func (s *SQLStore) getVolunteersByCondition(db sq.BaseRunner, condition interface{}, limit, offset uint64) ([]*model.Volunteer, error) {
	var volunteers []*model.Volunteer
	query := s.getQueryBuilder(db).
		Select(
			fmt.Sprintf("%svolunteers.id", s.tablePrefix),
			fmt.Sprintf("%svolunteers.name", s.tablePrefix),
			fmt.Sprintf("%svolunteers.contact", s.tablePrefix),
			fmt.Sprintf("%svolunteers.files", s.tablePrefix),
			fmt.Sprintf("%svolunteers.description", s.tablePrefix),
			fmt.Sprintf("%svolunteers.create_at", s.tablePrefix),
			fmt.Sprintf("%svolunteers.update_at", s.tablePrefix),
			fmt.Sprintf("%svolunteers.delete_at", s.tablePrefix),
		).
		From(s.tablePrefix + "volunteers").
		Where(sq.Eq{"delete_at": 0}).
		Where(condition)

	if limit != 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(10)
	}

	if offset != 0 {
		query = query.Offset(offset)
	}

	rows, err := query.Query()
	if err != nil {
		s.logger.Error(`getVolunteersList ERROR`, mlog.Err(err))
		return nil, err
	}
	defer s.CloseRows(rows)

	volunteers, err = s.volunteersFromRows(rows)
	if err != nil {
		return nil, err
	}

	if len(volunteers) == 0 {
		return nil, model.NewErrNotFound("volunteers")
	}

	return volunteers, nil
}
