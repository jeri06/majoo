package order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jeri06/majoo/exception"
	"github.com/jeri06/majoo/model"
)

type Repository interface {
	ReportOmzetByMerchant(ctx context.Context, param model.DailyReportMerchantParam, merchantId int, offset, size int64) (result []model.DailyReportMerchant, err error)
	ReportOmzetByMerchantTotalData(ctx context.Context, param model.DailyReportMerchantParam, merchantId int) (total int, err error)
	ReportOmzetOutlets(ctx context.Context, param model.DailyReportOutletParam, merchantId int, outletId int64, offset, size int64) (result []model.DailyReportOutlet, err error)
	ReportOmzetOutletsTotalData(ctx context.Context, param model.DailyReportOutletParam, merchantId int) (total int, err error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) Repository {
	return &orderRepository{
		db: db,
	}
}

func (a orderRepository) ReportOmzetByMerchant(ctx context.Context, param model.DailyReportMerchantParam, merchantId int, offset, size int64) (result []model.DailyReportMerchant, err error) {
	var cmd = a.db
	var rows *sql.Rows

	q := `SELECT  dd.dt as date ,
				m.merchant_name, 
				IFNULL(SUM(t.bill_total),0) AS omzet
				FROM (WITH RECURSIVE t as (
						select ? as dt
					UNION
						SELECT DATE_ADD(t.dt, INTERVAL 1 DAY) FROM t WHERE DATE_ADD(t.dt, INTERVAL 1 DAY) <= ?
				)
			select * FROM t) AS dd  
			left JOIN (SELECT * FROM transactions s WHERE s.merchant_id = ?) t ON  DATE_FORMAT(t.created_at,'%Y-%m-%d')= dd.dt 
				left JOIN merchants m ON m.id = ?
			GROUP BY dd.dt limit ? offset ?`

	if rows, err = cmd.QueryContext(ctx, q, param.StartDate, param.EndDate, merchantId, merchantId, size, offset); err != nil {
		fmt.Println(err.Error())
		err = exception.ErrInternalServer
		return
	}

	for rows.Next() {
		var report model.DailyReportMerchant
		err := rows.Scan(
			&report.Date, &report.MerchantName, &report.Omzet,
		)
		if err != nil {
			fmt.Println(err.Error())
			err = exception.ErrInternalServer

		}

		result = append(result, report)

	}
	return
}
func (a orderRepository) ReportOmzetByMerchantTotalData(ctx context.Context, param model.DailyReportMerchantParam, merchantId int) (total int, err error) {
	var cmd = a.db

	q := `SELECT COUNT(dd.dt)  FROM (WITH RECURSIVE t as (
			select ? as dt
		UNION
			SELECT DATE_ADD(t.dt, INTERVAL 1 DAY) FROM t WHERE DATE_ADD(t.dt, INTERVAL 1 DAY) <= ?
	)
	select * FROM t) AS dd  
	 left JOIN (SELECT * FROM transactions s WHERE s.merchant_id= ? ) t ON DATE_FORMAT(t.created_at,'%Y-%m-%d')= dd.dt
	`

	row := cmd.QueryRowContext(ctx, q, param.StartDate, param.EndDate, merchantId)
	err = row.Scan(&total)
	if err != nil {
		err = exception.ErrInternalServer
		return
	}
	return
}
func (a orderRepository) ReportOmzetOutlets(ctx context.Context, param model.DailyReportOutletParam, merchantId int, outletId int64, offset, size int64) (result []model.DailyReportOutlet, err error) {
	var cmd = a.db
	var rows *sql.Rows

	q := `SELECT  dd.dt as date ,
				m.merchant_name, 
				t.outlet_name,
				IFNULL(SUM(t.bill_total),0) AS omzet
				FROM (WITH RECURSIVE t as (
						select ? as dt
					UNION
						SELECT DATE_ADD(t.dt, INTERVAL 1 DAY) FROM t WHERE DATE_ADD(t.dt, INTERVAL 1 DAY) <= ?
				)
			select * FROM t) AS dd  
			left JOIN (SELECT s.*,o.outlet_name FROM transactions s JOIN outlets o ON s.outlet_id=o.id WHERE s.merchant_id = ? and s.outlet_id =?) t ON  DATE_FORMAT(t.created_at,'%Y-%m-%d')= dd.dt 
				left JOIN merchants m ON m.id = ?
			GROUP BY dd.dt,t.outlet_id limit ? offset ?`

	if rows, err = cmd.QueryContext(ctx, q, param.StartDate, param.EndDate, merchantId, outletId, merchantId, size, offset); err != nil {
		fmt.Println(err.Error())
		err = exception.ErrInternalServer
		return
	}

	for rows.Next() {
		var report model.DailyReportOutlet
		err := rows.Scan(
			&report.Date, &report.MerchantName, &report.OutletName, &report.Omzet,
		)
		if err != nil {
			fmt.Println(err.Error())
			err = exception.ErrInternalServer

		}

		result = append(result, report)

	}
	return
}
func (a orderRepository) ReportOmzetOutletsTotalData(ctx context.Context, param model.DailyReportOutletParam, merchantId int) (total int, err error) {
	var cmd = a.db

	q := `SELECT COUNT(dd.dt)  FROM (WITH RECURSIVE t as (
			select ? as dt
		UNION
			SELECT DATE_ADD(t.dt, INTERVAL 1 DAY) FROM t WHERE DATE_ADD(t.dt, INTERVAL 1 DAY) <= ?
	)
	select * FROM t) AS dd  
	 left JOIN (SELECT * FROM transactions s WHERE s.merchant_id= ? ) t ON DATE_FORMAT(t.created_at,'%Y-%m-%d')= dd.dt
	`

	row := cmd.QueryRowContext(ctx, q, param.StartDate, param.EndDate, merchantId)
	err = row.Scan(&total)
	if err != nil {
		err = exception.ErrInternalServer
		return
	}
	return
}
