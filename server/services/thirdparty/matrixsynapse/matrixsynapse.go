package matrixsynapse

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mattermost/focalboard/server/services/config"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

type MatrixSynapse struct {
	Config *config.Configuration
	Logger mlog.LoggerIFace
}

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type RegisterRequest struct {
	Nonce       string `json:"nonce"`
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
	Password    string `json:"password"`
	Admin       bool   `json:"admin"`
	Mac         string `json:"mac"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	HomeServer  string `json:"home_server"`
	DeviceID    string `json:"device_id"`
}

func New(config *config.Configuration, logger mlog.LoggerIFace) (*MatrixSynapse, error) {
	return &MatrixSynapse{Config: config, Logger: logger}, nil
}

func (ms *MatrixSynapse) getNonce() (*NonceResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/_synapse/admin/v1/register", ms.Config.MatrixSynapseUrl))
	if err != nil {
		ms.Logger.Error("MatrixSynapse, getNonce GET request failed", mlog.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			ms.Logger.Error("MatrixSynapse, getNonce failed to read response body", mlog.Err(err))
			return nil, err
		}
		ms.Logger.Error("MatrixSynapse, getNonce received non-OK HTTP status", mlog.Int("status", resp.StatusCode), mlog.String("body", string(bodyBytes)))
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var nonceResponse NonceResponse
	err = json.NewDecoder(resp.Body).Decode(&nonceResponse)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, getNonce decode response body failed", mlog.Err(err))
		return nil, err
	}
	return &nonceResponse, nil
}

func (ms *MatrixSynapse) createUser(requestBody RegisterRequest) (userID string, err error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, createUser marshal request body failed", mlog.Err(err))
		return userID, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/_synapse/admin/v1/register", ms.Config.MatrixSynapseUrl),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, createUser POST request failed", mlog.Err(err))
		return userID, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, createUser failed to read response body", mlog.Err(err))
		return userID, err
	}

	// Log the response body
	ms.Logger.Debug("MatrixSynapse, createUser response body",
		mlog.String("body", string(bodyBytes)),
		mlog.Int("status", resp.StatusCode),
	)

	if resp.StatusCode != http.StatusOK {
		ms.Logger.Error("MatrixSynapse, createUser received non-OK HTTP status", mlog.Int("status", resp.StatusCode))
		return userID, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var registerResponse RegisterResponse
	err = json.Unmarshal(bodyBytes, &registerResponse)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, createUser decode response body failed", mlog.Err(err))
		return userID, err
	}

	userID = registerResponse.UserID

	return userID, nil
}

func (ms *MatrixSynapse) generateHMAC(nonce, username, password string, admin bool, secret string) (string, error) {
	// Convert admin bool to string
	adminStr := "notadmin"
	if admin {
		adminStr = "admin"
	}

	// Create the message
	message := fmt.Sprintf("%s\x00%s\x00%s\x00%s", nonce, username, password, adminStr)

	ms.Logger.Debug("MatrixSynapse, generateHMAC input",
		mlog.String("message", message),
		mlog.String("secret", secret))

	// Create the HMAC
	h := hmac.New(sha1.New, []byte(secret))
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write to HMAC: %w", err)
	}

	// Get the result and encode as hex
	result := hex.EncodeToString(h.Sum(nil))

	ms.Logger.Debug("MatrixSynapse, generateHMAC result", mlog.String("hmac", result))

	return result, nil
}

func (ms *MatrixSynapse) RegisterUser(username string, displayname string, password string) (userID string, err error) {
	if ms.Config.MatrixSynapseSecret == "" {
		return "", nil
	}
	nonce, err := ms.getNonce()
	if err != nil {
		ms.Logger.Error("MatrixSynapse, RegisterUser get nonce failed", mlog.Err(err))
		return userID, err
	}

	mac, err := ms.generateHMAC(nonce.Nonce, username, password, false, ms.Config.MatrixSynapseSecret)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, RegisterUser generate HMAC failed", mlog.Err(err))
		return userID, err
	}

	registerRequest := RegisterRequest{
		Nonce:       nonce.Nonce,
		Username:    username,
		Displayname: displayname,
		Password:    password,
		Admin:       false,
		Mac:         mac,
	}

	userID, err = ms.createUser(registerRequest)
	if err != nil {
		ms.Logger.Error("MatrixSynapse, RegisterUser create user failed", mlog.Err(err))
		return userID, err
	}

	return userID, nil
}
