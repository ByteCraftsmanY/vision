package services

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/db"
	"vision/forms"
	"vision/models"
)

var (
	cctvMetaData = table.Metadata{
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
	cctvTable = table.New(cctvMetaData)
)

type CCTVService struct{}

func (s *CCTVService) CreateTable() error {
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

func (s *CCTVService) Add(form *forms.CCTV) (*models.CCTV, error) {
	cctv := models.CCTV{
		Username:       form.UserName,
		Password:       form.Password,
		URL:            form.URL,
		OrganizationID: form.OrganizationID,
		Extra:          form.Extra,
	}
	cctv.Initialize()
	session := db.GetSession()
	err := session.Query(cctvTable.Insert()).BindStruct(cctv).ExecRelease()
	return &cctv, err
}

func (s *CCTVService) GetByID(id gocql.UUID) (*models.CCTV, error) {
	cctv := models.CCTV{Base: models.Base{UUID: id}}
	session := db.GetSession()
	err := session.Query(cctvTable.Get()).BindStruct(cctv).GetRelease(&cctv)
	if err != nil {
		return nil, err
	}
	return &cctv, nil
}

func (s *CCTVService) GetAll(form *forms.CCTVPagination) (*models.CCTVList, error) {
	cctvList := make([]*models.CCTV, 0)
	session := db.GetSession()
	query := session.Query(cctvTable.SelectAll()).BindStruct(&models.CCTV{}).PageState(form.NextPageToken).PageSize(form.Size)
	defer query.Release()

	iter := query.Iter()
	err := iter.Select(&cctvList)
	return &models.CCTVList{
		Count:         len(cctvList),
		Results:       cctvList,
		NextPageToken: iter.PageState(),
	}, err
}

func (s *CCTVService) Update(form *forms.CCTVEdit) (*models.CCTV, error) {
	cctv := models.CCTV{
		Base:           models.Base{UUID: form.UUID},
		Username:       form.UserName,
		Password:       form.Password,
		URL:            form.URL,
		OrganizationID: form.OrganizationID,
		Extra:          form.Extra,
	}
	cctv.Update()
	columns := []string{"username", "password", "url", "organization_id", "extra"}
	session := db.GetSession()
	err := session.Query(cctvTable.Update(columns...)).BindStruct(&cctv).ExecRelease()
	return &cctv, err
}

func (s *CCTVService) Remove(id gocql.UUID) error {
	cctv := models.CCTV{Base: models.Base{UUID: id}}
	cctv.Delete()
	session := db.GetSession()
	err := session.Query(cctvTable.Delete()).BindStruct(cctv).ExecRelease()
	return err
}
