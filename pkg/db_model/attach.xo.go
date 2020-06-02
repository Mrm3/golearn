// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"database/sql"
	"errors"
	"time"
)

type Attach struct {
	Id           int          `json:"id"`
	DiskId       string       `json:"disk_id"`
	InstanceId   string       `json:"instance_id"`
	CvkName      string       `json:"cvk_name"`
	CvkNameOrig  string       `json:"cvk_name_orig"`
	AttachStatus int          `json:"attach_status"`
	CreateAt     time.Time    `json:"created_at"`
	UpdateAt     time.Time    `json:"updated_at"`
	DeleteAt     sql.NullTime `json:"delete_at"`
	IsDeleted    int          `json:"is_deleted"`

	// xo fields
	_exists, _deleted bool
}

type JoinedAttach struct {
	Attach
	Export Export
}

// Exists determines if the Disk exists in the database.
func (d *Attach) Exists() bool {
	return d._exists
}

// Deleted provides information if the Disk has been deleted from the database.
func (d *Attach) Deleted() bool {
	return d._deleted
}

// Save saves the Attach to the database.
func (d *Attach) Save(db XODB) error {
	if d.Exists() {
		return d.Update(db)
	}
	return d.Insert(db)
}

// Insert inserts the Attach to the database.
func (d *Attach) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if d._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query
	const sqlStr = `INSERT INTO attach (` +
		`disk_id, instance_id, cvk_name, cvk_name_orig, attach_status, created_at, updated_at, deleted_at, is_deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlStr, d.DiskId, d.InstanceId, d.CvkName, d.CvkNameOrig, d.AttachStatus, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted)
	_, err = db.Exec(sqlStr, d.DiskId, d.InstanceId, d.CvkName, d.CvkNameOrig, d.AttachStatus, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted)
	if err != nil {
		return err
	}

	// set existence
	d._exists = true

	return nil
}

// Update updates the Attach in the database.
func (d *Attach) Update(db XODB) error {
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
	const sqlStr = `UPDATE attach SET ` +
		`disk_id = ?, instance_id = ?, cvk_name = ?, cvk_name_orig = ?, attach_status = ?, created_at = ?, updated_at = ?, deleted_at = ?, is_deleted = ? ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlStr, d.DiskId, d.InstanceId, d.CvkName, d.CvkNameOrig, d.AttachStatus, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted, d.Id)
	_, err = db.Exec(sqlStr, d.DiskId, d.InstanceId, d.CvkName, d.CvkNameOrig, d.AttachStatus, d.CreateAt, d.UpdateAt, d.DeleteAt, d.IsDeleted, d.Id)
	return err
}

// Delete deletes the Attach from the database.
func (d *Attach) Delete(db XODB) error {
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
	const sqlStr = `UPDATE FROM attach SET delete_at = ?, is_deleted = ? WHERE id = ?`

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

func GetAttach(db XODB, DiskId string) (*Attach, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, instance_id, cvk_name, cvk_name_orig, attach_status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM attach ` +
		`WHERE is_deleted = 0 AND disk_id = ?`

	// run query
	XOLog(sqlStr, DiskId)
	d := Attach{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, DiskId).
		Scan(&d.Id, &d.DiskId, &d.InstanceId, &d.CvkName, &d.CvkNameOrig, &d.AttachStatus, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetAttach2(db XODB, DiskId, InstanceId string) (*Attach, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, instance_id, cvk_name, cvk_name_orig, attach_status, created_at, updated_at, deleted_at, is_deleted ` +
		`FROM attach ` +
		`WHERE is_deleted = 0 AND disk_id = ? AND instance_id = ?`

	// run query
	XOLog(sqlStr, DiskId, InstanceId)
	d := Attach{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, DiskId, InstanceId).
		Scan(&d.Id, &d.DiskId, &d.InstanceId, &d.CvkName, &d.CvkNameOrig, &d.AttachStatus, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetJoinedAttach(db XODB, DiskId string) (*JoinedAttach, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`a.id, a.disk_id, a.instance_id, a.cvk_name, a.cvk_name_orig, a.attach_status, a.created_at, a.updated_at, a.deleted_at, a.is_deleted, e.cvk_lun ` +
		`FROM attach a ` +
		`LEFT JOIN export e ON a.disk_id=e.disk_id AND a.cvk_name=e.cvk_name AND e.cvk_lun != 254 AND e.is_deleted = 0 ` +
		`WHERE a.is_deleted = 0 AND a.disk_id = ?`

	// run query
	XOLog(sqlStr, DiskId)
	var d JoinedAttach

	err = db.QueryRow(sqlStr, DiskId).
		Scan(&d.Id, &d.DiskId, &d.InstanceId, &d.CvkName, &d.CvkNameOrig, &d.AttachStatus, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted, &d.Export.CvkLun)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetJoinAttach2(db XODB, DiskId, InstanceId string) (*JoinedAttach, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`a.id, a.disk_id, a.instance_id, a.cvk_name, a.cvk_name_orig, a.attach_status, a.created_at, a.updated_at, a.deleted_at, a.is_deleted, e.cvk_lun ` +
		`FROM attach a ` +
		`LEFT JOIN export e ON a.disk_id=e.disk_id AND a.cvk_name=e.cvk_name AND e.cvk_lun != 254 AND e.is_deleted = 0 ` +
		`WHERE a.is_deleted = 0 AND a.disk_id = ? AND a.instance_id = ?`

	// run query
	XOLog(sqlStr, DiskId, InstanceId)
	var d JoinedAttach

	err = db.QueryRow(sqlStr, DiskId, InstanceId).
		Scan(&d.Id, &d.DiskId, &d.InstanceId, &d.CvkName, &d.CvkNameOrig, &d.AttachStatus, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.IsDeleted, &d.Export.CvkLun)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
