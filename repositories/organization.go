package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/daos"
	"vision/types"
)

type OrganizationRepositoryInterface interface {
	FindAll() ([]*daos.Organization, error)
	FindByID(id *types.ID) (*daos.Organization, error)
	Save(organization *daos.Organization) error
	Delete(organization *daos.Organization) error
}

type organizationRepository struct {
	session *gocqlx.Session
	table   *table.Table
}

func NewOrganizationRepository(session *gocqlx.Session) OrganizationRepositoryInterface {
	organization := new(daos.Organization)
	metaData := organization.GetTableMetaData()
	return &organizationRepository{
		session: session,
		table:   table.New(metaData),
	}
}

func (o *organizationRepository) FindAll() ([]*daos.Organization, error) {
	organization := new(daos.Organization)
	organizations := make([]*daos.Organization, 0)

	query := o.session.Query(o.table.SelectAll()).BindStruct(organization)
	defer query.Release()

	err := query.Select(&organizations)
	return organizations, err
}

func (o *organizationRepository) FindByID(id *types.ID) (*daos.Organization, error) {
	organization := &daos.Organization{Base: daos.Base{ID: id}}

	query := o.session.Query(o.table.Select()).BindStruct(organization)
	defer query.Release()

	err := query.Get(organization)
	return organization, err
}

func (o *organizationRepository) Save(organization *daos.Organization) error {
	query := o.session.Query(o.table.Insert()).BindStruct(organization)
	defer query.Release()

	return query.Exec()
}

func (o *organizationRepository) Delete(organization *daos.Organization) error {
	query := o.session.Query(o.table.Delete()).BindStruct(organization)
	defer query.Release()

	return query.Exec()
}

/*func (s *OrgService) CreateTable() error {
	session := db.GetSession()
	query := `CREATE TABLE IF NOT EXISTS organization (
				    id uuid PRIMARY KEY,
				    is_active boolean,
				    type varchar,
				    name text,
				    address text,
				    contact varchar,
				    cctv_count int,
				    user_count int,
				    bill_amount bigint,
				    associated_user_id uuid,
				    created_at timestamp,
					updated_at timestamp,
					deleted_at timestamp
				) WITH COMMENT = 'contains info about organization_id';`
	err := session.ExecStmt(query)
	if err != nil {
		return err
	}
	query = `CREATE INDEX IF NOT EXISTS org_user_id_idx ON organization (admin_id);`
	return session.ExecStmt(query)
}*/
