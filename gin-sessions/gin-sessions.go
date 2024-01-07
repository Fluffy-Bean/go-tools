package gin_sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type storage interface {
	Get(session *Session)
	Set(session Session) error
	Delete(key string)
	Clear()
}

type Sessions struct {
	Storage storage
	Name    string
	HashKey string
	MaxAge  int
}

type Session struct {
	Cookie        string
	UserID        string
	Authenticated bool
}

func (s *Sessions) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session Session

		session.Cookie, _ = c.Cookie(s.Name)

		s.Storage.Get(&session)

		c.Set("session", session)
		c.Next()
	}
}

func (s *Sessions) NewSession(c *gin.Context, userID string) {
	session := Session{
		UserID:        userID,
		Cookie:        uuid.NewV5(uuid.NamespaceOID, userID).String(),
		Authenticated: true,
	}

	if s.Storage.Set(session) != nil {
		panic("session already exists")
	}

	c.Set(s.Name, session.Cookie)
}

func (s *Sessions) GetSession(c *gin.Context) Session {
	cookie := c.MustGet(s.Name).(Session)
	s.Storage.Get(&cookie)
	return cookie
}

func (s *Sessions) DeleteSession(c *gin.Context) {
	cookie := s.GetSession(c)
	s.Storage.Delete(cookie.Cookie)
	c.SetCookie(s.Name, "", -1, "/", "", false, true)
}

func (s *Sessions) Clear() {
	s.Storage.Clear()
}

func NewMiddleware(storage storage, name string, hashKey string, maxAge int) *Sessions {
	return &Sessions{
		Storage: storage,
		Name:    name,
		HashKey: hashKey,
		MaxAge:  maxAge,
	}
}
