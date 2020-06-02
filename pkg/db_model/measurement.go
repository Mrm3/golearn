package db_model

import (
	"database/sql"
	"errors"
	"time"
)

type DiskMeasurement struct {
	ID       string       `json:"id"`
	DiskID   string       `json:"disk_id"`
	StartAt  sql.NullTime `json:"start_at"`
	StopAt   sql.NullTime `json:"stop_at"`
	Category int8         `json:"category"`
	Size     int64        `json:"size"`

	// xo fields
	_exists, _deleted bool
}

func (d *DiskMeasurement) Exists() bool {
	return d._exists
}

func (d *DiskMeasurement) Deleted() bool {
	return d._deleted
}

func (d *DiskMeasurement) Save(db XODB) error {
	if d.Exists() {
		return d.Update(db)
	}
	return d.Insert(db)
}

// Insert inserts the Attach to the database.
func (d *DiskMeasurement) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if d._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query
	const sqlStr = `INSERT INTO disk_measurement (` +
		`id, disk_id, start_at, stop_at, category, size` +
		`) VALUES ( ?, ?, ?, ?, ?, ? )`

	// run query
	XOLog(sqlStr, d.ID, d.DiskID, d.StartAt, d.StopAt, d.Category, d.Size)
	_, err = db.Exec(sqlStr, d.ID, d.DiskID, d.StartAt, d.StopAt, d.Category, d.Size)
	if err != nil {
		return err
	}

	// set existence
	d._exists = true

	return nil
}

// Update updates the Attach in the database.
func (d *DiskMeasurement) Update(db XODB) error {
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
	const sqlStr = `UPDATE disk_measurement SET` +
		` start_at = ?, stop_at = ?, category = ?, size = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlStr, d.StartAt, d.StopAt, d.Category, d.Size, d.ID)
	_, err = db.Exec(sqlStr, d.StartAt, d.StopAt, d.Category, d.Size, d.ID)
	return err
}

func MeasurementsByDiskID(db XODB, diskID string, startAt, stopAt time.Time) ([]DiskMeasurement, error) {
	var (
		measurements []DiskMeasurement
		err          error
	)

	const sqlStr = `SELECT id, disk_id, start_at, stop_at, category, size FROM disk_measurement` +
		` WHERE disk_id = ? AND (` +
		` start_at BETWEEN ? AND ?` +
		` OR stop_at BETWEEN ? AND ?` +
		` OR (start_at >= ? AND stop_at <= ?)` +
		` OR (start_at <= ? AND stop_at is NULL)` +
		`) ORDER BY start_at ASC`

	// run query
	XOLog(sqlStr, diskID, startAt, stopAt, startAt, stopAt, startAt, stopAt, startAt)
	q, err := db.Query(sqlStr, diskID, startAt, stopAt, startAt, stopAt, startAt, stopAt, startAt)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = q.Close()
	}()

	for q.Next() {
		m := DiskMeasurement{
			_exists: true,
		}

		err = q.Scan(&m.ID, &m.DiskID, &m.StartAt, &m.StopAt, &m.Category, &m.Size)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, m)
	}

	return measurements, nil
}

func LastMeasurementByDiskID(db XODB, diskID string) (DiskMeasurement, error) {
	var (
		err error
		m   = DiskMeasurement{
			_exists: true,
		}
	)
	const sqlStr = `SELECT id, disk_id, start_at, stop_at, category, size FROM disk_measurement` +
		` WHERE disk_id = ?` +
		` ORDER BY start_at DESC LIMIT 1 FOR UPDATE`

	// run query
	XOLog(sqlStr, diskID)
	err = db.QueryRow(sqlStr, diskID).
		Scan(&m.ID, &m.DiskID, &m.StartAt, &m.StopAt, &m.Category, &m.Size)

	return m, err
}
