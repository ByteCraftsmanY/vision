package models

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
	"time"
	"vision/db"
	"vision/forms"
	"vision/utils/token"
)

var userMetaData = table.Metadata{
	Name:    "user",
	Columns: []string{"id", "name", "email", "password", "phone", "address", "city", "state", "country", "active"},
	PartKey: []string{"id"},
	SortKey: nil,
}

var userTable = table.New(userMetaData)

type User struct {
	ID       gocql.UUID `json:"id,omitempty"`
	Name     string     `json:"name,omitempty"`
	Email    string     `json:"email,omitempty"`
	Password string     `json:"-"`
	Phone    string     `json:"phone,omitempty"`
	Address  string     `json:"address,omitempty"`
	City     string     `json:"city,omitempty"`
	State    string     `json:"state,omitempty"`
	Country  string     `json:"country,omitempty"`
	Active   bool       `json:"active,omitempty"`
}

func (u *User) SignUp(form *forms.UserSignUp) (*User, error) {
	user := User{
		ID:       gocql.UUIDFromTime(time.Now()),
		Password: form.Password,
		Name:     form.Name,
		Email:    form.Email,
		Phone:    form.Phone,
		Country:  form.Country,
		Active:   true,
	}

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
	user := User{ID: id}
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
	return token.Generate(user.ID)
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
