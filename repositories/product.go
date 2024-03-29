package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"vision/daos"
	"vision/types"
)

type ProductRepositoryInterface interface {
	FindAll() []*daos.Product
	FindByID(id *types.ID) (*daos.Product, error)
	Save(product *daos.Product) (*daos.Product, error)
	Delete(product *daos.Product) error
}

type productRepository struct {
	session *gocqlx.Session
	table   *table.Table
}

func NewProductRepository(session *gocqlx.Session) ProductRepositoryInterface {
	product := new(daos.Product)
	metaData := product.GetTableMetaData()
	return &productRepository{
		session: session,
		table:   table.New(metaData),
	}
}

func (p *productRepository) FindAll() []*daos.Product {
	//TODO implement me
	panic("implement me")
}

func (p *productRepository) FindByID(id *types.ID) (*daos.Product, error) {
	product := daos.Product{Base: daos.Base{ID: id}}
	query := p.session.Query(p.table.Get()).BindStruct(product)
	defer query.Release()

	err := query.Get(&product)
	return &product, err
}

func (p *productRepository) Save(product *daos.Product) (*daos.Product, error) {
	query := p.session.Query(p.table.Insert()).BindStruct(product)
	defer query.Release()

	err := query.Exec()
	return product, err
}

func (p *productRepository) Delete(product *daos.Product) error {
	query := p.session.Query(p.table.Delete()).BindStruct(product)
	defer query.Release()

	err := query.Exec()
	return err
}

/*func (s *CCTVService) CreateTable() error {
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
}*/
