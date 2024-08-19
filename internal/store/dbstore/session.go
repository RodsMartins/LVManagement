package dbstore

/* import (
	"lvm/internal/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

type NewSessionStoreParams struct {
	DB *gorm.DB
}

func NewSessionStore(params NewSessionStoreParams) *SessionStore {
	return &SessionStore{
		db: params.DB,
	}
}

func (s *SessionStore) CreateSession(session *store.Session) (*store.Session, error) {

	session.SessionID = uuid.New().String()

	result := s.db.Create(session)

	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
} */
