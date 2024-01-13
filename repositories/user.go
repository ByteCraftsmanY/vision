package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/entities"
	"vision/types"
)

type UserRepositoryInterface interface {
	FindAll() []*entities.User
	FindByID(id *types.ID) (*entities.User, error)
	FindByPhoneOrEmail(user *entities.User) (*entities.User, error)
	Save(user *entities.User) error
	Delete(user *entities.User) error
}

type userRepository struct {
	session *gocqlx.Session
	table   *table.Table
}

func NewUserRepository(session *gocqlx.Session) UserRepositoryInterface {
	user := new(entities.User)
	metaData := user.GetTableMetaData()
	return &userRepository{
		session: session,
		table:   table.New(metaData),
	}
}

func (r *userRepository) FindAll() []*entities.User {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) FindByID(userID *types.ID) (*entities.User, error) {
	user := entities.User{Base: entities.Base{ID: userID}}
	query := r.session.Query(r.table.Get()).BindStruct(user)
	defer query.Release()

	err := query.Select(&user)
	return &user, err
}

func (r *userRepository) Save(user *entities.User) error {
	query := r.session.Query(r.table.Insert()).BindStruct(user)
	defer query.Release()

	return query.ExecRelease()
}

func (r *userRepository) Delete(user *entities.User) error {
	query := r.session.Query(r.table.Delete()).BindStruct(user)
	defer query.Release()

	return query.ExecRelease()
}

func (r *userRepository) FindByPhoneOrEmail(user *entities.User) (*entities.User, error) {
	condition := qb.Eq("phone")
	if len(user.Phone) == 0 {
		condition = qb.Eq("email")
	}
	selectBuilder := qb.Select(r.table.Name()).Where(condition)

	query := r.session.Query(selectBuilder.ToCql()).BindStruct(user)
	defer query.Release()

	err := query.Get(user)
	return user, err
}

/*func (r *User) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
					id              uuid primary key,
					is_active 		boolean,
					username        text,
					password        text,
					url             text,
					organization_ids set<uuid>,
					extra           map<text, text>,
					created_at 		timestamp,
					updated_at 		timestamp,
					deleted_at 		timestamp,
				) WITH  COMMENT = 'contains info about user';`
	session := db.GetSession()
	err := session.ExecStmt(query)
	if err != nil {
		return err
	}
	query = `CREATE INDEX IF NOT EXISTS user_org_id_idx ON user (organization_id);`
	return session.ExecStmt(query)
}*/
