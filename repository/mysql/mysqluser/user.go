package mysqluser

import (
	"database/sql"
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"gameApp/repository/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)
	_, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
		//fmt.Errorf("can't scan query result: %w", err)
	}

	return false, nil
}

func (d *DB) Register(u entity.User) (entity.User, error) {
	res, err := d.conn.Conn().Exec(`insert into users(name, phone_number ,password ,role) values(?, ?, ?, ?)`, u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	// error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		// TODO - log unexpected err for better observability
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserById(userID uint) (entity.User, error) {
	const op = "mysql.GetUserById"
	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE id = ?`, userID)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var createdAt time.Time
	user := entity.User{}
	var roleStr string
	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &roleStr)

	user.Role = entity.MapRoleString(roleStr)

	return user, err

}
