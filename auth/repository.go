package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jeri06/majoo/exception"
	"github.com/jeri06/majoo/model"
)

type Repository interface {
	FindByUsername(ctx context.Context, userName string) (auth model.Authorization, err error)
	FindOutletByMerchantId(ctx context.Context, merchantId int64) (arr []int64, err error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) Repository {
	return &authRepository{
		db: db,
	}
}

func (a authRepository) FindByUsername(ctx context.Context, userName string) (auth model.Authorization, err error) {
	var cmd = a.db
	var rows *sql.Rows

	q := `SELECT u.name,u.user_name,u.password, m.id AS merchantId, m.merchant_name
  FROM users u JOIN merchants m ON u.id=m.user_id where user_name=?`

	if rows, err = cmd.QueryContext(ctx, q, userName); err != nil {
		fmt.Println("1", err)
		err = exception.ErrInternalServer

		return
	}

	for rows.Next() {

		err := rows.Scan(
			&auth.Name, &auth.UserName, &auth.Password, &auth.MerchandId, &auth.MerchandName,
		)
		if err != nil {
			fmt.Println(err)
			err = exception.ErrInternalServer

		}

	}
	return
}

func (a authRepository) FindOutletByMerchantId(ctx context.Context, merchantId int64) (arr []int64, err error) {
	var cmd = a.db
	var rows *sql.Rows

	q := `SELECT o.id from outlets o where o.merchant_id=? `

	if rows, err = cmd.QueryContext(ctx, q, merchantId); err != nil {
		fmt.Println(err)
		err = exception.ErrInternalServer

		return
	}

	for rows.Next() {
		var outletId int64
		err := rows.Scan(
			&outletId,
		)

		if err != nil {
			fmt.Println(err)
			err = exception.ErrInternalServer

		}
		arr = append(arr, outletId)
	}
	return
}
