package db_model

import (
	"database/sql"
	"time"
)

type DiskQoS struct {
	ID        string       `json:"id"`         // 主键
	DiskID    string       `json:"disk_id"`    // 磁盘ID
	LevelID   string       `json:"level_id"`   // QoSLevel ID
	GroupID   int          `json:"group_id"`   // QoSLevel NextGroupID
	CreatedAt time.Time    `json:"created_at"` // 创建时间
	UpdatedAt time.Time    `json:"updated_at"` // 更新时间
	DeletedAt sql.NullTime `json:"deleted_at"` // 删除时间
}

func (q *DiskQoS) Insert(db XODB) (err error) {
	const sqlStr = `INSERT INTO disk_qos (` +
		`id, disk_id, level_id, group_id, created_at, updated_at, deleted_at` +
		`) VALUES (?, ?, ?, ?, ?, ?, ?)`

	XOLog(sqlStr, q.ID, q.DiskID, q.LevelID, q.GroupID, q.CreatedAt, q.UpdatedAt, q.DeletedAt)
	_, err = db.Exec(sqlStr, q.ID, q.DiskID, q.LevelID, q.GroupID, q.CreatedAt, q.UpdatedAt, q.DeletedAt)

	return
}

func (q *DiskQoS) Update(db XODB) (err error) {
	const sqlStr = `UPDATE disk_qos SET ` +
		`disk_id = ?, level_id = ?, group_id = ?, created_at = ?, updated_at = ?, deleted_at = ? WHERE id = ?`

	XOLog(sqlStr, q.DiskID, q.LevelID, q.GroupID, q.CreatedAt, q.UpdatedAt, q.DeletedAt, q.ID)
	_, err = db.Exec(sqlStr, q.DiskID, q.LevelID, q.GroupID, q.CreatedAt, q.UpdatedAt, q.DeletedAt, q.ID)

	return
}

func (q *DiskQoS) Delete(db XODB) (err error) {
	const sqlStr = `UPDATE disk_qos SET deleted_at = ? WHERE id = ?`

	at := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	XOLog(sqlStr, at, q.ID)
	_, err = db.Exec(sqlStr, at, q.ID)

	return
}

// DiskQoSIDByDiskID 根据DiskID查询所对应的 Level ID
func DiskQoSIDByDiskID(db XODB, diskID string) (string, error) {
	const sqlStr = `SELECT level_id FROM disk_qos WHERE disk_id = ? AND deleted_at IS NULL`

	XOLog(sqlStr, diskID)

	var id string
	err := db.QueryRow(sqlStr, diskID).Scan(&id)

	return id, err
}

// DiskQoSByDiskIDForUpdate 根据DiskID查询所对应的 Level 并在事务下加锁
func DiskQoSByDiskIDForUpdate(db XODB, diskID string) (*DiskQoS, error) {
	const sqlStr = `SELECT id, disk_id, level_id, group_id, created_at, updated_at, deleted_at ` +
		`FROM disk_qos WHERE disk_id = ? AND deleted_at IS NULL FOR UPDATE`

	XOLog(sqlStr, diskID)

	var rule DiskQoS
	err := db.QueryRow(sqlStr, diskID).
		Scan(&rule.ID, &rule.DiskID, &rule.LevelID, &rule.GroupID, &rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt)

	return &rule, err
}
