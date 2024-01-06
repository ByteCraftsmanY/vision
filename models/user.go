package models

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
	"vision/db"
	"vision/forms"
	"vision/utils/token"
)

var userMetaData = table.Metadata{
	Name: "user",
	Columns: []string{
		"uuid", "is_active",
		"name", "email", "password", "phone",
		"created_at", "updated_at", "deleted_at",
	},
	PartKey: []string{"id"},
	SortKey: nil,
}

var userTable = table.New(userMetaData)

type User struct {
	base
	Name            string            `json:"name,omitempty" db:"name"`
	Phone           string            `json:"phone,omitempty" db:"phone"`
	Email           string            `json:"email,omitempty" db:"email"`
	Password        string            `json:"-" db:"password"`
	OrganizationIDs []gocql.UUID      `json:"organization_ids,omitempty" db:"organization_ids"`
	Extra           map[string]string `json:"extra,omitempty" db:"extra"`
}

type Users struct {
	Count         int     `json:"count,omitempty"`
	Results       []*User `json:"results,omitempty"`
	NextPageToken []byte  `json:"next_page_token,omitempty"`
}

func (u *User) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
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
				) WITH  COMMENT = 'contains info about user';`
	session := db.GetSession()
	err := session.ExecStmt(query)
	if err != nil {
		return err
	}
	query = `CREATE INDEX IF NOT EXISTS user_org_id_idx ON user (organization_id);`
	return session.ExecStmt(query)
}

func (u *User) Add(form *forms.User) (*User, error) {
	user := User{
		Password: form.Password,
		Name:     form.Name,
		Email:    form.Email,
		Phone:    form.Phone,
	}
	user.createInstance()

	if exists, err := user.isExists(); exists {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &user, err
	}
	user.Password = string(hashedPassword)

	session := db.GetSession()
	err = session.Query(userTable.Insert()).BindStruct(user).ExecRelease()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByID(id gocql.UUID) (*User, error) {
	user := User{base: base{UUID: id}}
	session := db.GetSession()
	err := session.Query(userTable.Get()).BindStruct(user).GetRelease(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByCredential(form *forms.Login) (*User, error) {
	user := User{}
	err := error(nil)
	session := db.GetSession()

	query := qb.Select(userTable.Name())
	if strings.TrimSpace(form.Email) != "" {
		user.Email = form.Email
		query = query.Where(qb.Eq("email"))
	} else if strings.TrimSpace(form.Phone) != "" {
		user.Phone = form.Phone
		query = query.Where(qb.Eq("phone"))
	} else {
		return nil, errors.New("email or phone is required")
	}

	err = session.Query(query.ToCql()).BindStruct(form).GetRelease(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) LoginCheck(form *forms.Login) (string, error) {
	user, err := u.GetByCredential(form)
	if err != nil {
		return "", err
	}
	err = VerifyPassword(form.Password, user.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", err
	}
	return token.Generate(user.UUID)
}

func (u *User) isExists() (bool, error) {
	session := db.GetSession()
	phoneQuery := qb.Select(userTable.Name()).Where(qb.Eq("phone"))
	emailQuery := qb.Select(userTable.Name()).Where(qb.Eq("email"))

	userByPhone, userByEmail := User{Phone: u.Phone}, User{Email: u.Email}
	userByPhoneErr, userByEmailErr := error(nil), error(nil)

	var wg sync.WaitGroup
	wg.Add(2)
	go func(user *User) {
		defer wg.Done()
		userByPhoneErr = session.Query(phoneQuery.ToCql()).BindStruct(user).GetRelease(user)
	}(&userByPhone)
	go func(user *User) {
		defer wg.Done()
		userByEmailErr = session.Query(emailQuery.ToCql()).BindStruct(user).GetRelease(user)
	}(&userByEmail)
	wg.Wait()

	if (userByPhoneErr == nil && len(userByPhone.Name) > 0) && (userByEmailErr == nil && len(userByEmail.Name) > 0) {
		return true, errors.New("user is already exists with this email and phone number")
	} else if userByPhoneErr == nil && len(userByPhone.Name) > 0 {
		return true, errors.New("user is already exists with this phone number")
	} else if userByEmailErr == nil && len(userByEmail.Name) > 0 {
		return true, errors.New("user is already exists with this email")
	}
	return false, nil
}
