package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/audit"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (a *API) registerPublicVolunteerRoutes(r *mux.Router) {
	// personal-server specific routes. These are not needed in plugin mode.
	r.HandleFunc("/register_volunteer", a.handleRegisterVolunteer).Methods("POST")
	r.HandleFunc("/upload_volunteer_files", a.handleUploadVolunteerFiles).Methods("POST")
}

func (a *API) registerVolunteerRoutes(r *mux.Router) {
	// personal-server specific routes. These are not needed in plugin mode.
	r.HandleFunc("/volunteers", a.sessionRequired(a.handleGetVolunteers)).Methods("GET")
	r.HandleFunc("/volunteer/{volunteerID}/create_user", a.sessionRequired(a.handleCreateUserFromVolunteer)).Methods("POST")
}

func (a *API) handleRegisterVolunteer(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /register_volunteer register_volunteer
	//
	// Register new volunteer
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Register volunteer request
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RegisterVolunteerRequest"
	// responses:
	//   '200':
	//     description: success
	//   '500':
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	var volunteerRegisterData model.Volunteer
	err = json.Unmarshal(requestBody, &volunteerRegisterData)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}
	volunteerRegisterData.Name = strings.TrimSpace(volunteerRegisterData.Name)
	volunteerRegisterData.Contact = strings.TrimSpace(volunteerRegisterData.Contact)

	auditRec := a.makeAuditRecord(r, "register_volunteer", audit.Fail)
	defer a.audit.LogRecord(audit.LevelAuth, auditRec)
	auditRec.AddMeta("name", volunteerRegisterData.Name)

	registeredVolunteer, err := a.app.RegisterVolunteer(volunteerRegisterData.Name, volunteerRegisterData.Contact, volunteerRegisterData.Description)
	if err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := jwt.MapClaims{
		"sub":     registeredVolunteer.ID,
		"exp":     expirationTime.Unix(),
		"name":    registeredVolunteer.Name,
		"contact": registeredVolunteer.Contact,
	}
	tokenString, err := a.app.EncodeJWTToken(claims)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}
	response := map[string]string{"access_token": tokenString}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}
	jsonBytesResponse(w, http.StatusOK, responseJSON)
	auditRec.Success()
}

func (a *API) handleUploadVolunteerFiles(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /volunteer/files uploadVolunteerFiles
	//
	// Upload a binary file, attached to a root block
	//
	// ---
	// consumes:
	// - multipart/form-data
	// produces:
	// - application/json
	// parameters:
	// - name: file
	//   in: formData
	//   type: file
	//   description: Array of files to upload
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//   '404':
	//     description: board not found
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		a.errorResponse(w, r, model.NewErrBadRequest("missing token"))
		return
	}
	headerParts := strings.Split(tokenString, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		a.errorResponse(w, r, model.NewErrBadRequest("invalid authorization header format"))
	}

	claims, err := a.app.DecodeJWTToken(headerParts[1])
	if err != nil {
		a.errorResponse(w, r, model.NewErrUnauthorized(fmt.Sprintf("%v", err)))
		return
	}
	sub := int64(claims["sub"].(float64))
	volunteer, err := a.app.GetVolunteer(sub)
	if err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(fmt.Sprintf("%v", err)))
		return
	}
	if volunteer.Name != claims["name"] || volunteer.Contact != claims["contact"] {
		a.errorResponse(w, r, model.NewErrBadRequest(fmt.Sprintf("%v", err)))
		return
	}

	// Parse the multipart form with a maximum memory of 100MB
	err = r.ParseMultipartForm(100 << 20) // 100MB
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Retrieve the files from the form-data
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		a.errorResponse(w, r, model.NewErrBadRequest("no files uploaded"))
		return
	}

	auditRec := a.makeAuditRecord(r, "upload_files", audit.Fail)
	defer a.audit.LogRecord(audit.LevelAuth, auditRec)

	_, err = a.app.SaveVolunteerFiles(volunteer, files)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	response := map[string]string{"message": "success"}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}
	jsonBytesResponse(w, http.StatusOK, responseJSON)
	auditRec.Success()
}

func (a *API) handleGetVolunteers(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /volunteers getVolunteers
	//
	// Returns the volunteers array
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: limit
	//   in: path
	//   description: Limit of records
	//   required: true
	//   type: number
	// - name: offset
	//   in: path
	//   description: Amount of skip records
	//   required: true
	//   type: number
	// - name: contact
	//   in: path
	//   description: Volunteer's contact for search
	//   required: true
	//   type: string
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Volunteer"
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"

	query := r.URL.Query()
	limitQuery := query.Get("limit")
	limit, err := strconv.ParseUint(limitQuery, 10, 64)
	if err != nil {
		limit = 10
	}
	offsetQuery := query.Get("offset")
	offset, err := strconv.ParseUint(offsetQuery, 10, 64)
	if err != nil {
		offset = 0
	}
	contact := query.Get("contact")
	userID := getUserID(r)

	if !a.permissions.HasPermissionTo(userID, model.PermissionManageSystem) {
		a.errorResponse(w, r, model.NewErrPermission("access denied to view volunteers"))
		return
	}

	auditRec := a.makeAuditRecord(r, "getVolunteers", audit.Fail)
	defer a.audit.LogRecord(audit.LevelModify, auditRec)
	auditRec.AddMeta("limit", limit)
	auditRec.AddMeta("offset", offset)

	volunteers, err := a.app.GetVolunteersList(contact, limit, offset)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.logger.Debug("GetVolunteers",
		mlog.String("contact", contact),
		mlog.Uint("limit", limit),
		mlog.Uint("offset", offset),
		mlog.Int("membersCount", len(volunteers)),
	)

	data, err := json.Marshal(volunteers)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// response
	jsonBytesResponse(w, http.StatusOK, data)

	auditRec.Success()
}

func (a *API) handleCreateUserFromVolunteer(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /volunteer/{volunteerID}/create_user createUserFromVolunteer
	//
	// Returns the user created from a volunteer
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: volunteerID
	//   in: path
	//   description: Board ID
	//   required: true
	//   type: string
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"

	volunteerIDStr := mux.Vars(r)["volunteerID"]
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest("Invalid volunteer ID"))
		return
	}
	userID := getUserID(r)

	if !a.permissions.HasPermissionTo(userID, model.PermissionManageSystem) {
		a.errorResponse(w, r, model.NewErrPermission("access denied to view volunteers"))
		return
	}

	auditRec := a.makeAuditRecord(r, "createUserFromVolunteer", audit.Fail)
	defer a.audit.LogRecord(audit.LevelModify, auditRec)
	auditRec.AddMeta("volunteerID", volunteerID)

	user, err := a.app.CreateUserFromVolunteer(volunteerID)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.logger.Debug("CreateUserFromVolunteer",
		mlog.Int("volunteerID", volunteerID),
	)

	data, err := json.Marshal(user)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// response
	jsonBytesResponse(w, http.StatusOK, data)

	auditRec.Success()
}
