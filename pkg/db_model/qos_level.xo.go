package db_model

type QoSLevel struct {
	ID          string `json:"id"`
	CapMin      int    `json:"cap_min"`       // 该等级最低容量
	CapMax      int    `json:"cap_max"`       // 该等级最高容量
	GroupCount  int    `json:"group_count"`   // 该等级包含的QoS Rule条数
	NextGroupID int    `json:"next_group_id"` // 该等级下次添加卷要使用的 QoS Rule ID
}

func (q *QoSLevel) Insert(db XODB) (err error) {
	const sqlStr = `INSERT INTO qos_level (` +
		`id, cap_min, cap_max, group_count, next_group_id` +
		`) VALUES (?, ?, ?, ?, ?)`

	XOLog(sqlStr, q.ID, q.CapMin, q.CapMax, q.GroupCount, q.NextGroupID)
	_, err = db.Exec(sqlStr, q.ID, q.CapMin, q.CapMax, q.GroupCount, q.NextGroupID)

	return
}

func (q *QoSLevel) Update(db XODB) (err error) {
	const sqlStr = `UPDATE qos_level SET ` +
		`cap_min = ?, cap_max = ?, group_count = ?, next_group_id = ? WHERE id = ?`

	XOLog(sqlStr, q.CapMin, q.CapMax, q.GroupCount, q.NextGroupID, q.ID)
	_, err = db.Exec(sqlStr, q.CapMin, q.CapMax, q.GroupCount, q.NextGroupID, q.ID)

	return
}

// QoSLevelByCap 根据容量查询所对应的 Level ID
func QoSLevelByCap(db XODB, cap int64) (string, error) {
	const sqlStr = `SELECT id FROM qos_level WHERE cap_min <= ? AND cap_max > ?`

	XOLog(sqlStr, cap, cap)

	var id string
	err := db.QueryRow(sqlStr, cap, cap).Scan(&id)

	return id, err
}

// QoSLevelByCapForUpdate 根据容量查询所对应的 Level
func QoSLevelByCapForUpdate(db XODB, cap int64) (*QoSLevel, error) {
	const sqlStr = `SELECT id, cap_min, cap_max, group_count, next_group_id FROM qos_level WHERE cap_min <= ? AND cap_max > ? FOR UPDATE`

	XOLog(sqlStr, cap, cap)

	var rule QoSLevel
	err := db.QueryRow(sqlStr, cap, cap).
		Scan(&rule.ID, &rule.CapMin, &rule.CapMax, &rule.GroupCount, &rule.NextGroupID)

	return &rule, err
}
