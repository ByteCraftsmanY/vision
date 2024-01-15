package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
	"vision/daos"
	"vision/types"
)

type AuthRepositoryInterface interface {
	FindByID(id *types.ID) (*daos.Auth, error)
	Find(data *daos.Auth) (*daos.Auth, error)
	Save(auth *daos.Auth) error
	SaveWithTTL(auth *daos.Auth, duration time.Duration) error
	Delete(auth *daos.Auth) error
}

type authRepository struct {
	session *gocqlx.Session
	table   *table.Table
}

func NewAuthRepository(session *gocqlx.Session) AuthRepositoryInterface {
	auth := new(daos.Auth)
	metaData := auth.GetTableMetaData()
	return &authRepository{
		session: session,
		table:   table.New(metaData),
	}
}

func (r *authRepository) FindByID(id *types.ID) (*daos.Auth, error) {
	auth := daos.Auth{
		//Base: daos.Base{ID: id},
	}
	query := r.session.Query(r.table.Select()).BindStruct(auth)
	defer query.Release()

	err := query.Get(auth)
	return &auth, err
}

func (r *authRepository) Find(auth *daos.Auth) (*daos.Auth, error) {
	query := r.session.Query(r.table.Select()).BindStruct(auth)
	defer query.Release()

	err := query.Get(auth)
	return auth, err
}

func (r *authRepository) Save(auth *daos.Auth) error {
	query := r.session.Query(r.table.Insert()).BindStruct(auth)
	defer query.Release()
	return query.Exec()
}

func (r *authRepository) SaveWithTTL(auth *daos.Auth, duration time.Duration) error {
	query := r.session.Query(r.table.InsertBuilder().TTL(duration).ToCql()).BindStruct(auth)
	defer query.Release()
	return query.Exec()
}

func (r *authRepository) Delete(auth *daos.Auth) error {
	query := r.session.Query(r.table.Delete()).BindStruct(auth)
	defer query.Release()

	return query.Exec()
}
