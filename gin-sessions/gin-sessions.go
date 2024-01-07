package gin_sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Storage interface {
	Get(session *Session)
	Set(session Session) error
	Delete(key string)
	Clear()
}

type Sessions struct {
	CookieName    string
	CookieHashKey string
	MaxAge        int
	CookieStorage Storage
}

type Session struct {
	Cookie        string
	UserID        string
	Authenticated bool
}

func (s *Sessions) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			session Session
		)

		session.Cookie, _ = c.Cookie(s.CookieName)

		s.CookieStorage.Get(&session)

		c.Set("session", session)
		c.Next()
	}
}

func (s *Sessions) NewSession(c *gin.Context, userID string) {
	var (
		session Session
	)

	session.Cookie = uuid.NewV5(uuid.NamespaceOID, userID).String()
	session.UserID = userID
	session.Authenticated = true

	if s.CookieStorage.Set(session) != nil {
		panic("session already exists")
	}

	c.Set(s.CookieName, session.Cookie)
}
