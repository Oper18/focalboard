package sqlpermissions

import (
	"github.com/mattermost/focalboard/server/services/permissions"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

type Service struct {
	store  permissions.Store
	logger mlog.LoggerIFace
}

func New(store permissions.Store, logger mlog.LoggerIFace) *Service {
	return &Service{
		store:  store,
		logger: logger,
	}
}

func (s *Service) HasPermissionTo(userID string, permission *mmModel.Permission) bool {
	user, err := s.store.HasUserPermissionTo(userID, permission.Id)
	if err != nil {
		s.logger.Error(`SQL HasPermissionTo ERROR`, mlog.Err(err))
		return false
	}
	if user == nil {
		return false
	}
	return true
}

func (s *Service) HasPermissionToTeam(userID, teamID string, permission *mmModel.Permission) bool {
	return s.HasPermissionTo(userID, permission)
}

func (s *Service) HasPermissionToChannel(userID, channelID string, permission *mmModel.Permission) bool {
	return s.HasPermissionTo(userID, permission)
}

func (s *Service) HasPermissionToBoard(userID, boardID string, permission *mmModel.Permission) bool {
	return s.HasPermissionTo(userID, permission)
}
