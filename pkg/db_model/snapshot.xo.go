// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"database/sql"
	"errors"
	"time"
)

// Snapshot represents a row from 'immortality.snapshot'.
type Snapshot struct {
	Id           int          `json:"id"`
	SnapshotID   string       `json:"snapshot_id"`   // snapshot_id
	Status       int8         `json:"status"`        // status
	DiskID       string       `json:"disk_id"`       // disk_id
	Name         string       `json:"name"`          // name
	Region       string       `json:"region"`        // region
	Zone         string       `json:"zone"`          // zone
	Category     int8         `json:"category"`      // category
	Size         int64        `json:"size"`          // size
	Description  string       `json:"description"`   // description
	UserID       string       `json:"user_id"`       // user_id
	ClusterID    string       `json:"cluster_id"`    // cluster_id
	StorageType  string       `json:"storage_type"`  // storage_type
	SnapshotType int          `json:"snapshot_type"` // snapshot_type
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	Deleted      int8         `json:"deleted"`
	ExpireAt     sql.NullTime `json:"expire_at"` // expire_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Snapshot exists in the database.
func (s *Snapshot) Exists() bool {
	return s._exists
}

// Save saves the Snapshot to the database.
func (s *Snapshot) Save(db XODB) error {
	if s.Exists() {
		return s.Update(db)
	}

	return s.Insert(db)
}

// Insert inserts the Snapshot to the database.
func (s *Snapshot) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if s._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided(18)
	const sqlStr = `INSERT INTO snapshot ( ` +
		`snapshot_id, status, disk_id, name, region, zone, category, size, description, user_id, cluster_id, storage_type, snapshot_type, ` +
		`created_at, updated_at, deleted_at, deleted, expire_at) ` +
		`VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// run query
	XOLog(sqlStr, s.SnapshotID, s.Status, s.DiskID, s.Name, s.Region, s.Zone, s.Category, s.Size, s.Description, s.UserID, s.ClusterID, s.StorageType, s.SnapshotType,
		s.CreatedAt, s.UpdatedAt, s.DeletedAt, s.Deleted, s.ExpireAt)
	_, err = db.Exec(sqlStr, s.SnapshotID, s.Status, s.DiskID, s.Name, s.Region, s.Zone, s.Category, s.Size, s.Description, s.UserID, s.ClusterID, s.StorageType, s.SnapshotType,
		s.CreatedAt, s.UpdatedAt, s.DeletedAt, s.Deleted, s.ExpireAt)
	if err != nil {
		return err
	}

	// set existence
	s._exists = true

	return nil
}

// Update updates the Snapshot in the database.
func (s *Snapshot) Update(db XODB) error {
	var err error

	// sql query
	const sqlStr = `UPDATE snapshot SET ` +
		`snapshot_id = ?, status = ?, disk_id = ?, name = ?, region = ?, zone = ?, category = ?, size = ?, description = ?, user_id = ?, ` +
		`cluster_id = ?, storage_type = ?, snapshot_type = ?, created_at = ?, updated_at = ?, deleted_at = ?, deleted = ?, expire_at = ? ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlStr, s.SnapshotID, s.Status, s.DiskID, s.Name, s.Region, s.Zone, s.Category, s.Size, s.Description, s.UserID, s.ClusterID, s.StorageType, s.SnapshotType,
		s.CreatedAt, s.UpdatedAt, s.DeletedAt, s.Deleted, s.ExpireAt, s.Id)
	_, err = db.Exec(sqlStr, s.SnapshotID, s.Status, s.DiskID, s.Name, s.Region, s.Zone, s.Category, s.Size, s.Description, s.UserID, s.ClusterID, s.StorageType, s.SnapshotType,
		s.CreatedAt, s.UpdatedAt, s.DeletedAt, s.Deleted, s.ExpireAt, s.Id)
	return err
}

// SnapshotsByDiskID retrieves a row from 'immortality.snapshot' as a Snapshot.
func SnapshotsByDiskID(db XODB, diskID string) ([]*Snapshot, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, snapshot_id, status, disk_id, name, region, zone, category, size, description, user_id, cluster_id, storage_type, snapshot_type, ` +
		`created_at, updated_at, deleted_at, deleted, expire_at ` +
		`FROM snapshot WHERE disk_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, diskID)
	q, err := db.Query(sqlStr, diskID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Snapshot{}
	for q.Next() {
		s := Snapshot{
			_exists: true,
		}

		// scan
		err = q.Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
			&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &s)
	}

	return res, nil
}

// SnapshotsByUserID retrieves a row from 'immortality.snapshot' as a Snapshot.
func SnapshotsByUserID(db XODB, userID string) ([]*Snapshot, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, snapshot_id, status, disk_id, name, region, zone, category, size, description, user_id, cluster_id, storage_type, snapshot_type, 
		created_at, updated_at, deleted_at, deleted, expire_at ` +
		`FROM snapshot WHERE user_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, userID)
	q, err := db.Query(sqlStr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Snapshot{}
	for q.Next() {
		s := Snapshot{
			_exists: true,
		}

		// scan
		err = q.Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
			&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &s)
	}

	return res, nil
}

// SnapshotBySnapshotID retrieves a row from 'immortality.snapshot' as a Snapshot.
func SnapshotBySnapshotID(db XODB, snapshotID string) (*Snapshot, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, snapshot_id, status, disk_id, name, region, zone, category, size, description, ` +
		`user_id, cluster_id, storage_type, snapshot_type, created_at, updated_at, deleted_at, deleted, expire_at ` +
		`FROM snapshot WHERE snapshot_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, snapshotID)
	s := Snapshot{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, snapshotID).Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
		&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func SnapshotBySnapshotIDForUpdate(db XODB, snapshotID string) (*Snapshot, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, snapshot_id, status, disk_id, name, region, zone, category, size, description, ` +
		`user_id, cluster_id, storage_type, snapshot_type, created_at, updated_at, deleted_at, deleted, expire_at ` +
		`FROM snapshot WHERE snapshot_id = ? AND deleted = 0 FOR UPDATE`

	// run query
	XOLog(sqlStr, snapshotID)
	s := Snapshot{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, snapshotID).Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
		&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func DeleteSnapshotByDiskId(db XODB, diskId string) (err error) {
	_, err = db.Exec("UPDATE snapshot SET deleted_at = ?, deleted = 1 WHERE disk_id = ? AND deleted = 0", time.Now(), diskId)
	return err
}

func AutoSnapshotsByDiskID(db XODB, diskID string) ([]*Snapshot, error) {
	// sql query
	const sqlStr = `SELECT id, snapshot_id, status, disk_id, name, region, zone, category, ` +
		`size, description, user_id, cluster_id, storage_type, snapshot_type, ` +
		`created_at, updated_at, deleted_at, deleted, expire_at FROM snapshot ` +
		`WHERE disk_id = ? AND snapshot_type = 1 AND deleted = 0 ORDER BY created_at`

	// run query
	result, err := db.Query(sqlStr, diskID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	// load results
	var res []*Snapshot
	for result.Next() {
		s := Snapshot{}
		// scan
		err = result.Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
			&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &s)
	}

	return res, nil
}

func SnapshotExistBySnapshotId(db XODB, snapshotId string) (bool, error) {
	// sql query
	const sqlStr = `SELECT * FROM snapshot WHERE snapshot_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, snapshotId)
	rows, err := db.Query(sqlStr, snapshotId)
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

func LastAutoSnapshotByDiskId(db XODB, diskId string) (*Snapshot, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`id, snapshot_id, status, disk_id, name, region, zone, category, size, description, user_id, cluster_id, ` +
		`storage_type, snapshot_type, created_at, updated_at, deleted_at, deleted, expire_at ` +
		`FROM snapshot WHERE disk_id = ? AND snapshot_type = 1 AND deleted = 0 ORDER BY created_at DESC LIMIT 1`

	// run query
	XOLog(sqlStr, diskId)
	s := Snapshot{
		_exists: true,
	}

	err = db.QueryRow(sqlStr, diskId).Scan(&s.Id, &s.SnapshotID, &s.Status, &s.DiskID, &s.Name, &s.Region, &s.Zone, &s.Category,
		&s.Size, &s.Description, &s.UserID, &s.ClusterID, &s.StorageType, &s.SnapshotType,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.Deleted, &s.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
