package repos

import (
	"log"

	"github.com/kallabs/sso-api/src/internal/app"
	"github.com/kallabs/sso-api/src/internal/app/valueobject"
	"github.com/kallabs/sso-api/src/internal/interfaces/db"
)

type UserRepo struct {
	db db.DB
}

func NewUserRepo(db db.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(obj app.User) (*app.User, error) {
	stmt := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, token, token_expires_at, status) 
		VALUES(:username, :email, :password_hash, :first_name, :last_name, :token, :token_expires_at, :status)`

	rows, err := r.db.Db().NamedQuery(stmt, obj)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(obj.Id)
	}

	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Println(err)
	}

	return &obj, nil
}

func (r *UserRepo) Get(userId *valueobject.ID) (*app.User, error) {
	stmt := `SELECT * FROM users WHERE id=?`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, userId); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindByEmail(email valueobject.EmailAddress) (*app.User, error) {
	stmt := `SELECT * FROM users WHERE email=?`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, email); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindByUsername(username string) (*app.User, error) {
	query := `SELECT * FROM users WHERE username=?`
	user := &app.User{}

	if err := r.db.Db().Get(user, query, username); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindByToken(token string) (*app.User, error) {
	stmt := `
		SELECT * FROM users 
		WHERE token=? AND token_expires_at > NOW() AND status=?`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, token, app.UserUnconfirmed); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) Delete(userId *valueobject.ID) error {
	stmt := `UPDATE users SET status=? WHERE id=?`

	if _, err := r.db.Db().Exec(stmt, app.UserDeleted, userId); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *UserRepo) Update(obj app.User) error {
	query := `
		UPDATE users 
		SET username=:username, email=:email, password_hash=:password_hash, 
			first_name=:first_name, last_name=:last_name, status=:status
		WHERE id=:id
	`
	if _, err := r.db.Db().NamedExec(query, obj); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
