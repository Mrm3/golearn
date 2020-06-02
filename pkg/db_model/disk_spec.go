package db_model

import (
	"database/sql"
	"errors"
	"time"
)

// DiskSpec represents a row from 'immortality.disk_spec'
type DiskSpec struct {
	ID              string       `json:"id"`
	CapacityMin     int64        `json:"capacity_min"` // 单盘最小容量
	CapacityMax     int64        `json:"capacity_max"` // 单盘最大容量
	IOPSBase        uint16       `json:"iops_base"`    // IOPS 基础值
	IOPSFactor      float64      `json:"iops_factor"`  // IOPS 系数
	IOPSMax         uint16       `json:"iops_max"`     // IOPS 最大值
	BandwidthBase   uint16       `json:"bw_base"`      // 吞吐量基础值
	BandwidthFactor float64      `json:"bw_factor"`    // 吞吐量系数
	BandwidthMax    uint16       `json:"bw_max"`       // 吞吐量最大值
	Family          string       `json:"family"`       // 规格族
	Code            string       `json:"code"`         // 规格编码
	Category        int8         `json:"category"`     // category
	CreatedAt       time.Time    `json:"created_at"`   // 创建时间
	UpdatedAt       time.Time    `json:"updated_at"`   // 最后更新时间
	DeletedAt       sql.NullTime `json:"deleted_at"`   // 删除时间

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the DiskSpec exists in the database.
func (d *DiskSpec) Exists() bool {
	return d._exists
}

// Save saves the Disk to the database.
func (d *DiskSpec) Save(db XODB) error {
	if d.Exists() {
		return d.Update(db)
	}

	return d.Insert(db)
}

// Insert inserts the Disk to the database.
func (d *DiskSpec) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if d._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided(24)
	const sqlStr = `INSERT INTO disk_spec (` +
		`id, capacity_min, capacity_max, iops_base, iops_factor, iops_max, bw_base, bw_factor, bw_max, family, code, category, created_at, updated_at, deleted_at` +
		`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// run query
	XOLog(sqlStr, d.ID, d.CapacityMin, d.CapacityMax,
		d.IOPSBase, d.IOPSFactor, d.IOPSMax,
		d.BandwidthBase, d.BandwidthFactor, d.BandwidthMax,
		d.Family, d.Code, d.Category,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt)
	_, err = db.Exec(sqlStr, d.ID, d.CapacityMin, d.CapacityMax,
		d.IOPSBase, d.IOPSFactor, d.IOPSMax,
		d.BandwidthBase, d.BandwidthFactor, d.BandwidthMax,
		d.Family, d.Code, d.Category,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt)
	if err != nil {
		return err
	}

	// set existence
	d._exists = true

	return nil
}

// Update updates the Disk in the database.
func (d *DiskSpec) Update(db XODB) error {
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
	const sqlStr = `UPDATE disk_spec SET` +
		` capacity_min = ?, capacity_max = ?,` +
		` iops_base = ?, iops_factor = ?, iops_max = ?, bw_base = ?, bw_factor = ?, bw_max = ?,` +
		` family = ?, code = ?, category = ?, created_at = ?, updated_at = ?, deleted_at = ? WHERE id = ?`

	// run query
	XOLog(sqlStr, d.CapacityMin, d.CapacityMax,
		d.IOPSBase, d.IOPSFactor, d.IOPSMax,
		d.BandwidthBase, d.BandwidthFactor, d.BandwidthMax,
		d.Family, d.Code, d.Category,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt,
		d.ID)
	_, err = db.Exec(sqlStr, d.CapacityMin, d.CapacityMax,
		d.IOPSBase, d.IOPSFactor, d.IOPSMax,
		d.BandwidthBase, d.BandwidthFactor, d.BandwidthMax,
		d.Family, d.Code, d.Category,
		d.CreatedAt, d.UpdatedAt, d.DeletedAt,
		d.ID)
	return err
}

func DiskSpecs(db XODB) (specs []DiskSpec, err error) {
	// sql query
	const sqlStr = `SELECT` +
		` id, capacity_min, capacity_max, iops_base, iops_factor, iops_max, bw_base, bw_factor, bw_max, family, code, category, created_at, updated_at, deleted_at` +
		` FROM disk_spec WHERE deleted_at IS NULL ORDER BY created_at`

	// run query
	XOLog(sqlStr)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return
	}

	for rows.Next() {
		spec := DiskSpec{
			_exists: true,
		}

		err = rows.Scan(&spec.ID, &spec.CapacityMin, &spec.CapacityMax,
			&spec.IOPSBase, &spec.IOPSFactor, &spec.IOPSMax,
			&spec.BandwidthBase, &spec.BandwidthFactor, &spec.BandwidthMax,
			&spec.Family, &spec.Code, &spec.Category,
			&spec.CreatedAt, &spec.UpdatedAt, &spec.DeletedAt)
		if err != nil {
			return
		}

		specs = append(specs, spec)
	}

	return
}

func DiskSpecByID(db XODB, id string) (DiskSpec, error) {
	// sql query
	const sqlStr = `SELECT` +
		` id, capacity_min, capacity_max, iops_base, iops_factor, iops_max, bw_base, bw_factor, bw_max, family, code, category, created_at, updated_at, deleted_at` +
		` FROM disk_spec WHERE id = ?`

	// run query
	XOLog(sqlStr, id)
	spec := DiskSpec{
		_exists: true,
	}

	err := db.QueryRow(sqlStr, id).
		Scan(&spec.ID, &spec.CapacityMin, &spec.CapacityMax,
			&spec.IOPSBase, &spec.IOPSFactor, &spec.IOPSMax,
			&spec.BandwidthBase, &spec.BandwidthFactor, &spec.BandwidthMax,
			&spec.Family, &spec.Code, &spec.Category,
			&spec.CreatedAt, &spec.UpdatedAt, &spec.DeletedAt)
	if err != nil {
		return spec, err
	}

	return spec, nil
}
