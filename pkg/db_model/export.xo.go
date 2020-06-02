// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"database/sql"
	"errors"
	"time"
)

type Export struct {
	Id        int          `json:"id"`
	DiskId    string       `json:"disk_id"`
	CvkName   string       `json:"cvk_name"`
	Iqn       string       `json:"iqn"`
	CvkLun    int          `json:"cvk_lun"`
	Status    int          `json:"status"`
	CreateAt  time.Time    `json:"created_at"`
	UpdateAt  time.Time    `json:"updated_at"`
	DeleteAt  sql.NullTime `json:"delete_at"`
	IsDeleted int          `json:"is_deleted"`

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Disk exists in the database.
func (d *Export) Exists() bool {
	return d._exists
}

// Deleted provides information if the Disk has been deleted from the database.
func (d *Export) Deleted() bool {
	return d._deleted
}

// Save saves the Export to the database.
func (d *Export) Save(db XODB) error {
	if d.Exists() {
		return d.Update(db)
	}
	return d.Insert(db)
}

// Insert inserts the Export to the database.
func (d *Export) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if d._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query
	const sqlStr = `INSERT INTO export (` +
		`disk_id, cvk_name, iqn, cvk_lun, status, created_at, updated_at, deleted_at, is_deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlStr, d.DiskId, d.CvkName, d.Iqn, d.CvkLun, d.Status, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted)
	_, err = db.Exec(sqlStr, d.DiskId, d.CvkName, d.Iqn, d.CvkLun, d.Status, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted)
	if err != nil {
		return err
	}

	// set existence
	d._exists = true

	return nil
}

// Update updates the Export in the database.
func (d *Export) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !d._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if d._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlStr = `UPDATE export SET ` +
		`disk_id = ?, cvk_name = ?, iqn = ?, cvk_lun = ?, status = ?, created_at = ?, updated_at = ?, deleted_at = ?, is_deleted = ? ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlStr, d.DiskId, d.CvkName, d.Iqn, d.CvkLun, d.Status, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted, d.Id)
	_, err = db.Exec(sqlStr, d.DiskId, d.CvkName, d.Iqn, d.CvkLun, d.Status, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted, d.Id)
	return err
}

// Delete deletes the Export from the database.
func (d *Export) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !d._exists {
		return nil
	}

	// if deleted, bail
	if d._deleted {
		return nil
	}

	// sql query
	const sqlStr = `UPDATE FROM export SET delete_at = ?, is_deleted = ? WHERE id = ?`

	// run query
	XOLog(sqlStr, d.DeleteAt, d.IsDeleted, d.Id)
	_, err = db.Exec(sqlStr, d.DeleteAt, d.IsDeleted, d.Id)
	if err != nil {
		return err
	}

	// set deleted
	d._deleted = true

	return nil
}

func GetExport(db XODB, DiskId, CvkName string) (*Export, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, cvk_name, iqn, cvk_lun, status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM export ` +
		`WHERE is_deleted = 0 AND disk_id = ? AND cvk_name = ?`

	// run query
	XOLog(sqlStr, DiskId, CvkName)
	d := Export{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, DiskId, CvkName).
		Scan(&d.Id, &d.DiskId, &d.CvkName, &d.Iqn, &d.CvkLun, &d.Status, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

//get lun not 254
func GetExport3(db XODB, DiskId, CvkName string) (*Export, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, cvk_name, iqn, cvk_lun, status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM export ` +
		`WHERE is_deleted = 0 AND cvk_lun != 254 AND disk_id = ? AND cvk_name = ?`

	// run query
	XOLog(sqlStr, DiskId, CvkName)
	d := Export{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, DiskId, CvkName).
		Scan(&d.Id, &d.DiskId, &d.CvkName, &d.Iqn, &d.CvkLun, &d.Status, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetExport4(db XODB, lun int, DiskId, CvkName string) (*Export, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, cvk_name, iqn, cvk_lun, status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM export ` +
		`WHERE is_deleted = 0 AND disk_id = ? AND cvk_name = ? AND cvk_lun = ?`

	// run query
	XOLog(sqlStr, DiskId, CvkName, lun)
	d := Export{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, DiskId, CvkName, lun).
		Scan(&d.Id, &d.DiskId, &d.CvkName, &d.Iqn, &d.CvkLun, &d.Status, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetExports(db XODB, diskID string) ([]*Export, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, cvk_name, iqn, cvk_lun, status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM export ` +
		`WHERE is_deleted = 0 AND disk_id = ?`

	// run query
	XOLog(sqlStr, diskID)
	q, err := db.Query(sqlStr, diskID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Export{}
	for q.Next() {
		d := Export{
			_exists: true,
		}

		// scan
		err = q.Scan(&d.Id, &d.DiskId, &d.CvkName, &d.Iqn, &d.CvkLun, &d.Status, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &d)
	}

	return res, nil

}
