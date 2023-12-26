package postgres

import (
	"log"

	"github.com/jmoiron/sqlx"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

type UserGateway struct {
	Db *sqlx.DB
}

func (i *UserGateway) Create(user_schama *entities.CreateUserSchema) (*entities.User, error) {
	stmt := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, token, token_expires_at, status) 
		VALUES(:username, :email, :password_hash, :first_name, :last_name, :token, :token_expires_at, :status)`

	rows, err := i.Db.NamedQuery(stmt, user_schama)

	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Username: user_schama.Username,
		Email:    user_schama.Email,
		UserPassword: entities.UserPassword{
			PasswordHash: user_schama.PasswordHash,
		},
		FirstName: user_schama.FirstName,
		LastName:  user_schama.LastName,
		Role:      user_schama.Role,
		Status:    user_schama.Status,
	}

	if rows.Next() {
		rows.Scan(user.Id)
	}

	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Println(err)
	}

	return user, nil
}

func (i *UserGateway) GetById(user_id value_objects.ID) (*entities.User, error) {
	stmt := `SELECT * FROM users WHERE id=$1`
	user := &entities.User{}

	if err := i.Db.Get(user, stmt, user_id); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (i *UserGateway) GetByEmail(email value_objects.EmailAddress) (*entities.User, error) {
	stmt := `SELECT * FROM users WHERE email=$1`
	user := &entities.User{}

	if err := i.Db.Get(user, stmt, email); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (i *UserGateway) GetByUsername(username string) (*entities.User, error) {
	stmt := `SELECT * FROM users WHERE username=$1`
	user := &entities.User{}

	if err := i.Db.Get(user, stmt, username); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (i *UserGateway) GetByRegToken(reg_token string) (*entities.User, error) {
	stmt := `
		SELECT * FROM users 
		WHERE token=$1 AND token_expires_at > CURRENT_TIMESTAMP AND status=$2`
	user := &entities.User{}

	if err := i.Db.Get(user, stmt, reg_token, entities.UserUnconfirmed); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (i *UserGateway) UpdateFromEntity(user *entities.User) error {
	query := `
		UPDATE users 
		SET username=:username, email=:email, password_hash=:password_hash, 
			first_name=:first_name, last_name=:last_name, status=:status
		WHERE id=:id
	`
	if _, err := i.Db.NamedExec(query, user); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *UserGateway) MarkAsDeleted(userId *value_objects.ID) error {
	stmt := `UPDATE users SET status=$1 WHERE id=$2`

	if _, err := i.Db.Exec(stmt, entities.UserDeleted, userId); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
