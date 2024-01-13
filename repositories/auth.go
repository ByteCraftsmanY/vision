package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
	"vision/entities"
	"vision/types"
)

type AuthRepositoryInterface interface {
	FindByID(id *types.ID) (*entities.Auth, error)
	Find(data *entities.Auth) (*entities.Auth, error)
	Save(auth *entities.Auth) error
	SaveWithTTL(auth *entities.Auth, duration time.Duration) error
	Delete(auth *entities.Auth) error
}

type authRepository struct {
	session *gocqlx.Session
	table   *table.Table
}

func NewAuthRepository(session *gocqlx.Session) AuthRepositoryInterface {
	auth := new(entities.Auth)
	metaData := auth.GetTableMetaData()
	return &authRepository{
		session: session,
		table:   table.New(metaData),
	}
}

func (r *authRepository) FindByID(id *types.ID) (*entities.Auth, error) {
	auth := entities.Auth{
		//Base: entities.Base{ID: id},
	}
	query := r.session.Query(r.table.Select()).BindStruct(auth)
	defer query.Release()

	err := query.Get(auth)
	return &auth, err
}

func (r *authRepository) Find(auth *entities.Auth) (*entities.Auth, error) {
	query := r.session.Query(r.table.Select()).BindStruct(auth)
	defer query.Release()

	err := query.Get(auth)
	return auth, err
}

func (r *authRepository) Save(auth *entities.Auth) error {
	query := r.session.Query(r.table.Insert()).BindStruct(auth)
	defer query.Release()
	return query.Exec()
}

func (r *authRepository) SaveWithTTL(auth *entities.Auth, duration time.Duration) error {
	query := r.session.Query(r.table.InsertBuilder().TTL(duration).ToCql()).BindStruct(auth)
	defer query.Release()
	return query.Exec()
}

func (r *authRepository) Delete(auth *entities.Auth) error {
	query := r.session.Query(r.table.Delete()).BindStruct(auth)
	defer query.Release()

	return query.Exec()
}
