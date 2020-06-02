// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Disk represents a row from 'immortality.disk'.
type Disk struct {
	Id             int          `json:"id"`
	DiskID         string       `json:"disk_id"`       // disk_id
	StatusOrig     int8         `json:"status_orig"`   // status_orig
	Status         int8         `json:"status"`        // status
	DiskType       int8         `json:"disk_type"`     // disk_type
	Name           string       `json:"name"`          // name
	Region         string       `json:"region"`        // region
	Zone           string       `json:"zone"`          // zone
	Category       int8         `json:"category"`      // category
	Size           int64        `json:"size"`          // size
	Description    string       `json:"description"`   // description
	FromSnapshot   string       `json:"from_snapshot"` // from_snapshot
	UserID         string       `json:"user_id"`       // user_id
	ClusterID      string       `json:"cluster_id"`    // cluster_id
	FromImage      string       `json:"from_image"`    // from_image
	StorageType    string       `json:"storage_type"`
	IsShare        int8         `json:"is_share"`
	Qos            string       `json:"qos"`
	ThreeParWWN    string       `json:"three_par_wwn"`
	ThreeParStatus string       `json:"three_par_status"`
	Deleted        int8         `json:"deleted"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Disk exists in the database.
func (d *Disk) Exists() bool {
	return d._exists
}

// Deleted provides information if the Disk has been deleted from the database.
func (d *Disk) IsDeleted() bool {
	return d._deleted
}

// Save saves the Disk to the database.
func (d *Disk) Save(db XODB) error {
	if d.Exists() {
		return d.Update(db)
	}

	return d.Insert(db)
}

// Insert inserts the Disk to the database.
func (d *Disk) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if d._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided(24)
	const sqlStr = `INSERT INTO disk (` +
		`disk_id, status_orig, status, disk_type, name, region, zone, category, size, description, from_snapshot, user_id, cluster_id,` +
		`from_image, storage_type, is_share, qos, three_par_wwn, three_par_status,` +
		`created_at, updated_at, deleted_at, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlStr, d.DiskID, d.StatusOrig, d.Status, d.DiskType, d.Name, d.Region, d.Zone, d.Category, d.Size, d.Description, d.FromSnapshot, d.UserID, d.ClusterID,
		d.FromImage, d.StorageType, d.IsShare, d.Qos, d.ThreeParWWN, d.ThreeParStatus,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt, d.Deleted)
	_, err = db.Exec(sqlStr, d.DiskID, d.StatusOrig, d.Status, d.DiskType, d.Name, d.Region, d.Zone, d.Category, d.Size, d.Description, d.FromSnapshot, d.UserID, d.ClusterID,
		d.FromImage, d.StorageType, d.IsShare, d.Qos, d.ThreeParWWN, d.ThreeParStatus,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt, d.Deleted)
	if err != nil {
		return err
	}

	// set existence
	d._exists = true

	return nil
}

// Update updates the Disk in the database.
func (d *Disk) Update(db XODB) error {
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
	const sqlStr = `UPDATE disk SET disk_id = ?, status_orig = ?, status = ?, disk_type = ?, name = ?, region = ?, zone = ?, category = ?, size = ?, description = ?, from_snapshot = ?,` +
		`user_id = ?, cluster_id = ?, from_image = ?, storage_type = ?, is_share = ?, qos = ?, three_par_wwn = ?, three_par_status = ?,` +
		`created_at = ?, updated_at = ? ,deleted_at = ?, deleted = ? WHERE id = ?`

	// run query
	XOLog(sqlStr, d.DiskID, d.StatusOrig, d.Status, d.DiskType, d.Name, d.Region, d.Zone, d.Category, d.Size,
		d.Description, d.FromSnapshot, d.UserID, d.ClusterID, d.FromImage,
		d.StorageType, d.IsShare, d.Qos, d.ThreeParWWN, d.ThreeParStatus,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt, d.Deleted, d.Id)
	_, err = db.Exec(sqlStr, d.DiskID, d.StatusOrig, d.Status, d.DiskType, d.Name, d.Region, d.Zone, d.Category, d.Size,
		d.Description, d.FromSnapshot, d.UserID, d.ClusterID, d.FromImage,
		d.StorageType, d.IsShare, d.Qos, d.ThreeParWWN, d.ThreeParStatus,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt, d.Deleted, d.Id)
	return err
}

// Generated from index 'disk_disk_id_pkey'.
func DiskByDiskID(db XODB, diskID string) (*Disk, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, disk_id, status_orig, status, disk_type, name, region, zone, category, size, description, from_snapshot, user_id, cluster_id, ` +
		`from_image, storage_type, is_share, qos, three_par_wwn, three_par_status, created_at, updated_at, deleted_at, deleted ` +
		`FROM disk WHERE disk_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, diskID)
	d := Disk{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, diskID).
		Scan(&d.Id, &d.DiskID, &d.StatusOrig, &d.Status, &d.DiskType, &d.Name, &d.Region, &d.Zone, &d.Category,
			&d.Size, &d.Description, &d.FromSnapshot, &d.UserID, &d.ClusterID, &d.FromImage,
			&d.StorageType, &d.IsShare, &d.Qos, &d.ThreeParWWN, &d.ThreeParStatus,
			&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.Deleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func DisksByDiskIDs(db XODB, diskIDs []string) ([]*Disk, error) {
	var err error

	n := len(diskIDs)
	ps := make([]string, n)
	for i := 0; i < n; i++ {
		ps[i] = "?"
	}

	var sqlStr = fmt.Sprintf(`SELECT * FROM disk WHERE disk_id in (%s) AND deleted = 0`, strings.Join(ps, ","))

	args := make([]interface{}, len(diskIDs))
	for i, id := range diskIDs {
		args[i] = id
	}

	XOLog(sqlStr, diskIDs)

	rows, err := db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disks []*Disk

	for rows.Next() {
		d := &Disk{
			_exists: true,
		}

		err = rows.Scan(&d.Id, &d.DiskID, &d.StatusOrig, &d.Status, &d.DiskType, &d.Name, &d.Region, &d.Zone, &d.Category,
			&d.Size, &d.Description, &d.FromSnapshot, &d.UserID, &d.ClusterID, &d.FromImage,
			&d.StorageType, &d.StorageType, &d.Qos, &d.ThreeParWWN, &d.ThreeParStatus,
			&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.Deleted)
		if err != nil {
			return nil, err
		}

		disks = append(disks, d)
	}

	return disks, nil
}

// DiskByDiskIDForUpdate retrieves a row from 'immortality.disk' as a Disk, and locks the row.
func DiskByDiskIDForUpdate(db XODB, diskID string) (*Disk, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, disk_id, status_orig, status, disk_type, name, region, zone, category, size, description, from_snapshot, user_id, cluster_id, ` +
		`from_image, storage_type, is_share, qos, three_par_wwn, three_par_status, created_at, updated_at, deleted_at, deleted ` +
		`FROM disk WHERE disk_id = ? AND deleted = 0 ` +
		`FOR UPDATE `

	// run query
	XOLog(sqlstr, diskID)
	d := Disk{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, diskID).
		Scan(&d.Id, &d.DiskID, &d.StatusOrig, &d.Status, &d.DiskType, &d.Name, &d.Region, &d.Zone, &d.Category,
			&d.Size, &d.Description, &d.FromSnapshot, &d.UserID, &d.ClusterID, &d.FromImage,
			&d.StorageType, &d.IsShare, &d.Qos, &d.ThreeParWWN, &d.ThreeParStatus,
			&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.Deleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func DiskExistByDiskId(db XODB, diskId string) (bool, error) {
	// sql query
	const sqlStr = `SELECT * FROM disk WHERE disk_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, diskId)
	rows, err := db.Query(sqlStr, diskId)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	}

	return false, err
}
