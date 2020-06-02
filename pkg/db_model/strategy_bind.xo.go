package db_model

import "time"

type StrategyBind struct {
	DiskId            string    `json:"disk_id"`
	StrategyId        string    `json:"strategy_id"`
	AutoSnapshotQuota int       `json:"auto_snapshot_quota"`
	BindTime          time.Time `json:"bind_at"`
}

func (s *StrategyBind) Insert(db XODB) error {
	var err error

	const sqlStr = `INSERT INTO snapshot_strategy_binding (disk_id, strategy_id, auto_snapshot_quota, bind_at) VALUES (?, ?, ?, ?)`

	// run query
	XOLog(sqlStr, s.DiskId, s.StrategyId, s.AutoSnapshotQuota, s.BindTime)
	_, err = db.Exec(sqlStr, s.DiskId, s.StrategyId, s.AutoSnapshotQuota, s.BindTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *StrategyBind) Delete(db XODB) error {
	var err error

	const sqlStr = `DELETE FROM snapshot_strategy_binding WHERE disk_id = ?`

	// run query
	XOLog(sqlStr, s.DiskId)
	_, err = db.Exec(sqlStr, s.DiskId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBindByStrategyId(db XODB, strategyId string) error {
	var err error

	const sqlStr = `DELETE FROM snapshot_strategy_binding WHERE strategy_id = ?`

	// run query
	XOLog(sqlStr, strategyId)
	_, err = db.Exec(sqlStr, strategyId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBindByDiskId(db XODB, diskId string) error {
	var err error

	const sqlStr = `DELETE FROM snapshot_strategy_binding WHERE disk_id = ?`

	// run query
	XOLog(sqlStr, diskId)
	_, err = db.Exec(sqlStr, diskId)
	if err != nil {
		return err
	}

	return nil
}

func BindByDiskId(db XODB, diskId string) (*StrategyBind, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`disk_id, strategy_id, auto_snapshot_quota, bind_at ` +
		`FROM snapshot_strategy_binding ` +
		`WHERE disk_id = ?`

	// run query
	XOLog(sqlStr, diskId)
	s := StrategyBind{}

	err = db.QueryRow(sqlStr, diskId).Scan(&s.DiskId, &s.StrategyId, &s.AutoSnapshotQuota, &s.BindTime)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// if exist, the disk can't bind new strategy
func BindExistByDiskId(db XODB, diskId string) (bool, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`disk_id, strategy_id, auto_snapshot_quota, bind_at ` +
		`FROM snapshot_strategy_binding ` +
		`WHERE disk_id = ?`

	// run query
	XOLog(sqlStr, diskId)
	rows, err := db.Query(sqlStr, diskId)
	if err != nil {
		return true, err
	}
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

// if exist, can't delete strategy
func BindExistByStrategyId(db XODB, strategyId string) (bool, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`disk_id, strategy_id, auto_snapshot_quota, bind_at ` +
		`FROM snapshot_strategy_binding ` +
		`WHERE strategy_id = ?`

	// run query
	XOLog(sqlStr, strategyId)
	rows, err := db.Query(sqlStr, strategyId)
	if err != nil {
		return true, err
	}
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetAllBinds(db XODB) ([]*StrategyBind, error) {
	var err error
	// sql query
	const sqlStr = `SELECT disk_id, strategy_id, auto_snapshot_quota, bind_at FROM snapshot_strategy_binding `

	// run query
	XOLog(sqlStr)
	result, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	// load results
	var res []*StrategyBind
	for result.Next() {
		bind := StrategyBind{}
		// scan
		err = result.Scan(&bind.DiskId, &bind.StrategyId, &bind.AutoSnapshotQuota, &bind.BindTime)
		if err != nil {
			return nil, err
		}

		res = append(res, &bind)
	}

	return res, nil
}

// GetDiskCountByStrategyId 获取指定策略已绑定磁盘数量
func GetDiskCountByStrategyId(db XODB, strategyId string) (count int, err error) {
	const sqlStr = `SELECT COUNT(DISTINCT disk_id) FROM snapshot_strategy_binding WHERE strategy_id = ?`

	XOLog(sqlStr, strategyId)
	err = db.QueryRow(sqlStr, strategyId).Scan(&count)

	return
}
