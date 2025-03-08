package app

import (
	"crypto/rand"
	"encoding/base64"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/google/uuid"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/auth"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mozillazg/go-slugify"
	"github.com/pkg/errors"
)

// RegisterVolunteer creates a new volunteer if the provided data is valid.
func (a *App) RegisterVolunteer(name, contact, description string) (*model.Volunteer, error) {
	if contact == "" {
		volunteers, err := a.store.GetVolunteersList(contact, 1, 0)
		if err != nil && !model.IsErrNotFound(err) {
			return nil, err
		}
		if len(volunteers) > 0 {
			return nil, errors.New("Volunteer with such contact already exist")
		}
	}
	if name == "" {
		return nil, errors.New("The name is required")
	}

	volunteer := model.Volunteer{
		Name:        name,
		Contact:     contact,
		Description: description,
	}
	registeredVolunteer, err := a.store.CreateVolunteer(&volunteer)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create the new volunteer")
	}
	return registeredVolunteer, nil
}

func (a *App) GetVolunteer(volunteerID int64) (*model.Volunteer, error) {
	volunteer, err := a.store.GetVolunteer(volunteerID)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get the volunteer")
	}
	return volunteer, nil
}

func (a *App) GetVolunteersList(contact string, limit uint64, offset uint64) ([]*model.Volunteer, error) {
	volunteers, err := a.store.GetVolunteersList(contact, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get the volunteers list")
	}
	return volunteers, nil
}

func (a *App) SaveVolunteerFiles(volunteer *model.Volunteer, files []*multipart.FileHeader) (*model.Volunteer, error) {
	volunteerID := strconv.FormatInt(volunteer.ID, 10)
	volunteerDir := filepath.Join("volunteers_files", volunteerID)
	err := os.MkdirAll(volunteerDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create dir for volunteer")
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			a.logger.Error("Failed to read file", mlog.Err(err))
			continue
		}
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)

		uuidFileName := uuid.New().String() + ext
		filePath := filepath.Join(volunteerDir, uuidFileName)
		_, err = a.filesBackend.WriteFile(file, filePath)
		if err != nil {
			a.logger.Error("Failed to save file", mlog.Err(err))
			continue
		}
		volunteer.Files = append(volunteer.Files, uuidFileName)
	}
	volunteer, err = a.store.UpdateVolunteer(volunteer)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to update volunteer")
	}
	return volunteer, nil
}

func (a *App) CreateUserFromVolunteer(volunteerID int64) (*model.User, error) {
	volunteer, err := a.store.GetVolunteer(volunteerID)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get the volunteer")
	}
	info := whatlanggo.Detect(volunteer.Name)
	var username string
	if info.Lang == whatlanggo.Rus {
		// If it's Russian, transliterate it
		username = slugify.Slugify(volunteer.Name)
	} else {
		// If it's not Russian, just remove spaces and convert to lowercase
		username = strings.ToLower(strings.ReplaceAll(volunteer.Name, " ", ""))
	}
	username = strings.ReplaceAll(username, " ", "_")
	nameParts := strings.Fields(volunteer.Name)
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}
	bytes := make([]byte, 12)
	_, err = rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	password := base64.URLEncoding.EncodeToString(bytes)[:12]
	role, err := a.store.GetRoleByName("volunteer")
	if err != nil {
		return nil, errors.Wrap(err, "Role not found")
	}

	userID, err := a.thirdParty.RegisterUser(username, volunteer.Contact, password)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create the new matrix user")
	}

	userModel := model.User{
		ID:           utils.NewID(utils.IDTypeUser),
		Username:     username,
		Email:        volunteer.Contact,
		Password:     auth.HashPassword(password),
		MfaSecret:    "",
		AuthService:  a.config.AuthMode,
		AuthData:     "",
		RoleID:       role.ID,
		MatrixUserID: userID,
		FirstName:    firstName,
		LastName:     lastName,
	}
	createdUser, err := a.store.CreateUser(&userModel)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create the new user")
	}
	return createdUser, nil
}
