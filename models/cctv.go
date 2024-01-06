package models

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/db"
	"vision/forms"
)

var cctvMetaData = table.Metadata{
	Name: "cctv",
	Columns: []string{
		"uuid", "is_active",
		"username", "password", "url",
		"organization_id", "extra",
		"created_at", "updated_at", "deleted_at",
	},
	PartKey: []string{"uuid"},
	SortKey: nil,
}

var cctvTable = table.New(cctvMetaData)

type CCTV struct {
	base
	Username       string            `json:"username,omitempty"`
	Password       string            `json:"-"`
	URL            string            `json:"url,omitempty"`
	OrganizationID gocql.UUID        `json:"organization_id,omitempty"`
	Extra          map[string]string `json:"extra,omitempty"`
}

type CCTVList struct {
	Count         int     `json:"count,omitempty"`
	Results       []*CCTV `json:"results,omitempty"`
	NextPageToken []byte  `json:"next_page_token,omitempty"`
}

func (c *CCTV) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS cctv (
					uuid            uuid primary key,
					is_active 		boolean,
					username        text,
					password        text,
					url             text,
					organization_id uuid,
					extra           map<text, text>,
					created_at 		timestamp,
					updated_at 		timestamp,
					deleted_at 		timestamp,
				) WITH  COMMENT = 'contains info about cctv camera';`
	session := db.GetSession()
	err := session.ExecStmt(query)
	if err != nil {
		return err
	}
	query = `CREATE INDEX IF NOT EXISTS cctv_org_id_idx ON cctv (organization_id);`
	return session.ExecStmt(query)
}

func (c *CCTV) Add(form *forms.CCTV) (*CCTV, error) {
	cctv := CCTV{
		Username:       form.UserName,
		Password:       form.Password,
		URL:            form.URL,
		OrganizationID: form.OrganizationID,
		Extra:          form.Extra,
	}
	cctv.createInstance()
	session := db.GetSession()
	err := session.Query(cctvTable.Insert()).BindStruct(cctv).ExecRelease()
	return &cctv, err
}

func (c *CCTV) GetByID(id gocql.UUID) (*CCTV, error) {
	cctv := CCTV{base: base{UUID: id}}
	session := db.GetSession()
	err := session.Query(cctvTable.Get()).BindStruct(cctv).GetRelease(&cctv)
	if err != nil {
		return nil, err
	}
	return &cctv, nil
}

func (c *CCTV) GetAll(form *forms.CCTVPagination) (*CCTVList, error) {
	cctvList := make([]*CCTV, 0)
	session := db.GetSession()
	query := session.Query(cctvTable.SelectAll()).BindStruct(&CCTV{}).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&cctvList)
	return &CCTVList{
		Count:         len(cctvList),
		Results:       cctvList,
		NextPageToken: iter.PageState(),
	}, err
}

func (c *CCTV) Update(form *forms.CCTVEdit) (*CCTV, error) {
	cctv := CCTV{
		base:           base{UUID: form.UUID},
		Username:       form.UserName,
		Password:       form.Password,
		URL:            form.URL,
		OrganizationID: form.OrganizationID,
		Extra:          form.Extra,
	}
	cctv.updateInstance()
	columns := []string{"username", "password", "url", "organization_id", "extra"}
	session := db.GetSession()
	err := session.Query(cctvTable.Update(columns...)).BindStruct(&cctv).ExecRelease()
	return &cctv, err
}

func (c *CCTV) Remove(id gocql.UUID) error {
	cctv := CCTV{base: base{UUID: id}}
	cctv.deleteInstance()
	session := db.GetSession()
	err := session.Query(cctvTable.Delete()).BindStruct(cctv).ExecRelease()
	return err
}
