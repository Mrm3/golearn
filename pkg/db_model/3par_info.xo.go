// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"time"
)

type ThreeParInfo struct {
	ThreeParId         string    `json:"three_par_id"`
	HDDUsage           float64   `json:"hdd_usage"`
	HDDTotal           float64   `json:"hdd_total"`
	HDDUsagePercent    float64   `json:"hdd_usage_percent"`
	SSDUsage           float64   `json:"ssd_usage"`
	SSDTotal           float64   `json:"ssd_total"`
	SSDUsagePercent    float64   `json:"ssd_usage_percent"`
	HybridUsagePercent float64   `json:"hybrid_usage_percent"`
	UpdateAt           time.Time `json:"update_at"`
	PodId              string    `json:"pod_id"`
	Reserve            int       `json:"reserve"`
}

// Update updates the Disk in the database.
func (c *ThreeParInfo) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE 3par_info SET ` +
		`hdd_usage = ?, hdd_total = ?, hdd_usage_percent = ?, ssd_usage = ?, ssd_total = ?, ssd_usage_percent = ?, hybrid_usage_percent = ?, update_at = ?` +
		` WHERE three_par_id = ?`

	// run query
	XOLog(sqlstr, c.HDDUsage, c.HDDTotal, c.HDDUsagePercent, c.SSDUsage, c.SSDTotal, c.SSDUsagePercent, c.HybridUsagePercent, c.UpdateAt, c.ThreeParId)
	_, err = db.Exec(sqlstr, c.HDDUsage, c.HDDTotal, c.HDDUsagePercent, c.SSDUsage, c.SSDTotal, c.SSDUsagePercent, c.HybridUsagePercent, c.UpdateAt, c.ThreeParId)
	return err
}

func GetThreeParInfos(db XODB) ([]*ThreeParInfo, error) {
	var err error

	// sql query
	const sqlstr = `SELECT three_par_id, pod_id FROM 3par_info `

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*ThreeParInfo
	for q.Next() {
		info := ThreeParInfo{}
		// scan
		err = q.Scan(&info.ThreeParId, &info.PodId)
		if err != nil {
			return nil, err
		}

		res = append(res, &info)
	}

	return res, nil
}
