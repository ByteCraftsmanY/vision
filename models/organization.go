package models

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/db"
	"vision/forms"
)

var organizationMetaData = table.Metadata{
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

var organizationTable = table.New(organizationMetaData)

type Organization struct {
	base
	Name       string     `json:"name,omitempty"`
	Contact    string     `json:"contact,omitempty"`
	Type       string     `json:"type,omitempty"`
	BillAmount int64      `json:"bill_amount,omitempty"`
	Address    string     `json:"address,omitempty"`
	CCTVCount  int64      `json:"cctv_count,omitempty"`
	AdminID    gocql.UUID `json:"admin_id,omitempty"`
	UserCount  int64      `json:"user_count,omitempty"`
}

type Organizations struct {
	Count         int             `json:"count,omitempty"`
	Results       []*Organization `json:"results,omitempty"`
	NextPageToken []byte          `json:"next_page_token,omitempty"`
}

func (o *Organization) CreateTable() error {
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

func (o *Organization) Add(form *forms.OrganizationNew) (*Organization, error) {
	organization := Organization{
		Name:    form.Name,
		Contact: form.Contact,
		Type:    form.Type,
		Address: form.Address,
		AdminID: form.AdminID,
	}
	organization.createInstance()
	session := db.GetSession()
	err := session.Query(organizationTable.Insert()).BindStruct(organization).ExecRelease()
	return &organization, err
}

func (o *Organization) GetByID(id gocql.UUID) (*Organization, error) {
	organization := Organization{
		base: base{UUID: id},
	}
	session := db.GetSession()
	err := session.Query(organizationTable.Get()).BindStruct(organization).GetRelease(&organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (o *Organization) GetAll(form *forms.OrganizationPagination) (*Organizations, error) {
	organizationList := make([]*Organization, 0)
	organization := Organization{
		Contact: form.Contact,
		Type:    form.Type,
		AdminID: form.AdminID,
	}
	session := db.GetSession()
	query := session.Query(organizationTable.SelectAll()).BindStruct(&organization).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&organizationList)
	return &Organizations{
		Count:         len(organizationList),
		Results:       organizationList,
		NextPageToken: iter.PageState(),
	}, err
}

func (o *Organization) Update(form *forms.OrganizationEdit) (*Organization, error) {
	organization := Organization{
		base: base{UUID: form.OrganizationID},
	}
	organization.updateInstance()
	columns := []string{"name", "contact", "admin_id"}
	session := db.GetSession()
	err := session.Query(organizationTable.Update(columns...)).BindStruct(&organization).ExecRelease()
	return &organization, err
}

func (o *Organization) Remove(id gocql.UUID) error {
	organization := Organization{base: base{UUID: id}}
	organization.deleteInstance()
	session := db.GetSession()
	err := session.Query(organizationTable.Delete()).BindStruct(organization).ExecRelease()
	return err
}
