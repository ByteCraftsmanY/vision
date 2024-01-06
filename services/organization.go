package services

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/db"
	"vision/forms"
	"vision/models"
)

var (
	orgMetaData = table.Metadata{
		Name: "organization",
		Columns: []string{
			"uuid", "is_active",
			"type", "admin_id",
			"name", "contact", "address",
			"cctv_count", "user_count", "bill_amount",
			"created_at", "updated_at", "deleted_at",
		},
		PartKey: []string{"uuid"},
		SortKey: nil,
	}
	orgTable = table.New(orgMetaData)
)

type OrgService struct{}

func (s *OrgService) CreateTable() error {
	session := db.GetSession()
	query := `CREATE TABLE IF NOT EXISTS organization (	
				    uuid uuid PRIMARY KEY, 
				    is_active boolean,
				    type varchar, 
				    name text, 
				    address text,
				    contact varchar, 
				    cctv_count int,
				    user_count int,
				    bill_amount bigint,
				    admin_id uuid,
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
}

func (s *OrgService) Add(form *forms.OrganizationNew) (*models.Org, error) {
	organization := models.Org{
		Name:    form.Name,
		Contact: form.Contact,
		Type:    form.Type,
		Address: form.Address,
		AdminID: form.AdminID,
	}
	organization.Initialize()
	session := db.GetSession()
	err := session.Query(orgTable.Insert()).BindStruct(organization).ExecRelease()
	return &organization, err
}

func (s *OrgService) GetByID(id gocql.UUID) (*models.Org, error) {
	organization := models.Org{
		Base: models.Base{UUID: id},
	}
	session := db.GetSession()
	err := session.Query(orgTable.Get()).BindStruct(organization).GetRelease(&organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (s *OrgService) GetAll(form *forms.OrganizationPagination) (*models.OrgList, error) {
	organizationList := make([]*models.Org, 0)
	organization := models.Org{
		Contact: form.Contact,
		Type:    form.Type,
		AdminID: form.AdminID,
	}
	session := db.GetSession()
	query := session.Query(orgTable.SelectAll()).BindStruct(&organization).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&organizationList)
	return &models.OrgList{
		Count:         len(organizationList),
		Results:       organizationList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *OrgService) Update(form *forms.OrganizationEdit) (*models.Org, error) {
	organization := models.Org{
		Base: models.Base{UUID: form.OrganizationID},
	}
	organization.Update()
	columns := []string{"name", "contact", "admin_id"}
	session := db.GetSession()
	err := session.Query(orgTable.Update(columns...)).BindStruct(&organization).ExecRelease()
	return &organization, err
}

func (s *OrgService) Remove(id gocql.UUID) error {
	organization := models.Org{Base: models.Base{UUID: id}}
	organization.Delete()
	session := db.GetSession()
	err := session.Query(orgTable.Delete()).BindStruct(organization).ExecRelease()
	return err
}
