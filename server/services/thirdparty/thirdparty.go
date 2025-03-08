package thirdparty

type ThirdParty interface {
	RegisterUser(username string, displayname string, password string) (userID string, err error)
}
