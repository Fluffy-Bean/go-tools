package gin_sessions

import "errors"

type InMemoryStorage struct {
	sessions map[string]string
}

func (s *InMemoryStorage) Get(session *Session) {
	session.UserID = s.sessions[session.Cookie]
	session.Authenticated = session.UserID != ""
}

func (s *InMemoryStorage) Set(session Session) error {
	if s.sessions[session.Cookie] != "" {
		return errors.New("session already exists")
	}

	s.sessions[session.Cookie] = session.UserID
	return nil
}

func (s *InMemoryStorage) Delete(key string) {
	delete(s.sessions, key)
}

func (s *InMemoryStorage) Clear() {
	s.sessions = make(map[string]string)
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		sessions: make(map[string]string),
	}
}
