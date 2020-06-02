package db_model

import (
	"time"
)

type Strategy struct {
	StrategyId   string    `json:"strategy_id"`
	StrategyName string    `json:"strategy_name"`
	UserId       string    `json:"user_id"`
	Hours        string    `json:"hours"`
	Weeks        string    `json:"weeks"`
	Status       int32     `json:"status"`
	Duration     int32     `json:"duration"`
	DiskQuota    int32     `json:"disk_quota"`
	CreateTime   time.Time `json:"create_at"`
	UpdateTime   time.Time `json:"update_at"`
}

func (s *Strategy) Insert(db XODB) error {
	var err error

	const sqlstr = `INSERT INTO snapshot_strategy (` +
		`strategy_id, strategy_name, user_id, hours, weeks, duration, status, disk_quota, create_at, update_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, s.StrategyId, s.StrategyName, s.UserId, s.Hours, s.Weeks, s.Duration, s.Status, s.DiskQuota, s.CreateTime, s.UpdateTime)
	_, err = db.Exec(sqlstr, s.StrategyId, s.StrategyName, s.UserId, s.Hours, s.Weeks, s.Duration, s.Status, s.DiskQuota, s.CreateTime, s.UpdateTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *Strategy) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE snapshot_strategy SET strategy_name = ?, user_id = ?, hours = ?, weeks = ?, duration = ?, status = ?, disk_quota = ?, create_at = ?, update_at = ?` +
		` WHERE strategy_id = ?`

	// run query
	XOLog(sqlstr, s.StrategyName, s.UserId, s.Hours, s.Weeks, s.Duration, s.Status, s.DiskQuota, s.CreateTime, s.UpdateTime, s.StrategyId)
	_, err = db.Exec(sqlstr, s.StrategyName, s.UserId, s.Hours, s.Weeks, s.Duration, s.Status, s.DiskQuota, s.CreateTime, s.UpdateTime, s.StrategyId)
	return err
}

func (s *Strategy) Delete(db XODB) error {
	var err error

	const sqlStr = `DELETE FROM snapshot_strategy WHERE strategy_id = ?`
	XOLog(sqlStr, s.StrategyId)
	_, err = db.Exec(sqlStr, s.StrategyId)
	return err
}

func StrategyByStrategyId(db XODB, strategyId string) (*Strategy, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`strategy_id, strategy_name, user_id, hours, weeks, duration, status, disk_quota, create_at, update_at ` +
		`FROM snapshot_strategy ` +
		`WHERE strategy_id = ?`

	// run query
	XOLog(sqlstr, strategyId)
	s := Strategy{}

	err = db.QueryRow(sqlstr, strategyId).Scan(&s.StrategyId, &s.StrategyName, &s.UserId, &s.Hours, &s.Weeks, &s.Duration, &s.Status, &s.DiskQuota,
		&s.CreateTime, &s.UpdateTime)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func StrategyByIDForUpdate(db XODB, strategyId string) (*Strategy, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`strategy_id, strategy_name, user_id, hours, weeks, duration, status, disk_quota, create_at, update_at ` +
		`FROM snapshot_strategy ` +
		`WHERE strategy_id = ? FOR UPDATE`

	// run query
	XOLog(sqlstr, strategyId)
	s := Strategy{}

	err = db.QueryRow(sqlstr, strategyId).Scan(&s.StrategyId, &s.StrategyName, &s.UserId,
		&s.Hours, &s.Weeks, &s.Duration, &s.Status, &s.DiskQuota, &s.CreateTime, &s.UpdateTime)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func StrategyExistByStrategyId(db XODB, strategyId string) (bool, error) {
	const sqlstr = `SELECT ` +
		`strategy_id, strategy_name, user_id, hours, weeks, duration, status, disk_quota, create_at, update_at ` +
		`FROM snapshot_strategy ` +
		`WHERE strategy_id = ?`

	// run query
	XOLog(sqlstr, strategyId)
	rows, err := db.Query(sqlstr, strategyId)
	if err != nil {
		return true, err
	}
	if rows.Next() {
		return true, err
	}

	return false, err
}
