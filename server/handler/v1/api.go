package v1

type API struct {
	*Middleware
	*User
	*OAuth2
	*Session
}

func NewAPI(
	middleware *Middleware,
	user *User,
	oAuth2 *OAuth2,
	session *Session,
) *API {
	return &API{
		Middleware:   middleware,
		User:         user,
		OAuth2:       oAuth2,
		Session:      session,
	}
}
