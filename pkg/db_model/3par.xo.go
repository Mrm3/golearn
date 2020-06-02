// Package db_model contains the types for schema 'immortality'.
package db_model

import (
	"database/sql"
	"time"
)

type ThreePar struct {
	ThreeParId         string       `json:"three_par_id"`
	Ipv4AddrManagement string       `json:"ipv4_addr_management"`
	Ipv4AddrSSH        string       `json:"ipv4_addr_ssh"`
	Ipv4AddrController string       `json:"ipv4_addr_controller"`
	Username           string       `json:"username"`
	Password           string       `json:"password"`
	LaunchTime         time.Time    `json:"launch_time"`
	InCharge           string       `json:"in_charge"`
	Contact            string       `json:"contact"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	DeletedAt          sql.NullTime `json:"deleted_at"`
	Deleted            int          `json:"deleted"`
}

func (t *ThreePar) Insert(db XODB) error {
	var err error

	const sqlStr = `INSERT INTO 3par (` +
		`three_par_id, ipv4_addr_management, ipv4_addr_ssh, ipv4_addr_controller, username, password, launch_time, in_charge, contact, created_at, updated_at, deleted_at, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlStr, t.ThreeParId, t.Ipv4AddrManagement, t.Ipv4AddrSSH, t.Ipv4AddrController, t.Username, t.Password, t.LaunchTime, t.InCharge, t.Contact, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Deleted)
	_, err = db.Exec(sqlStr, t.ThreeParId, t.Ipv4AddrManagement, t.Ipv4AddrSSH, t.Ipv4AddrController, t.Username, t.Password, t.LaunchTime, t.InCharge, t.Contact, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Deleted)
	if err != nil {
		return err
	}

	return nil
}

func (t *ThreePar) Update(db XODB) error {
	var err error

	const sqlStr = `UPDATE 3par SET ipv4_addr_management = ?, ipv4_addr_ssh = ?, ipv4_addr_controller = ?, username = ?, password = ?, ` +
		`launch_time = ?, in_charge = ?, contact = ?, created_at = ?, updated_at = ?, deleted_at = ?, deleted = ? WHERE three_par_id = ?`

	// run query
	XOLog(sqlStr, t.Ipv4AddrManagement, t.Ipv4AddrSSH, t.Ipv4AddrController, t.Username, t.Password,
		t.LaunchTime, t.InCharge, t.Contact, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Deleted, t.ThreeParId)
	_, err = db.Exec(sqlStr, t.Ipv4AddrManagement, t.Ipv4AddrSSH, t.Ipv4AddrController, t.Username, t.Password,
		t.LaunchTime, t.InCharge, t.Contact, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Deleted, t.ThreeParId)
	if err != nil {
		return err
	}

	return nil
}

func (t *ThreePar) Delete(db XODB) error {
	var err error
	var timeNow = time.Now()

	const sqlStr = `UPDATE 3par SET deleted_at = ?, deleted = 1 WHERE three_par_id = ?`

	// run query
	XOLog(sqlStr, timeNow, t.ThreeParId)
	_, err = db.Exec(sqlStr, timeNow, t.ThreeParId)
	if err != nil {
		return err
	}

	return nil
}

func ThreePars(db XODB) ([]*ThreePar, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`three_par_id, ipv4_addr_management, ipv4_addr_ssh, ipv4_addr_controller, username, password, launch_time, in_charge, contact, created_at, updated_at, deleted_at, deleted ` +
		`FROM 3par WHERE deleted = 0 `

	// run query
	XOLog(sqlStr)
	q, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*ThreePar
	for q.Next() {
		d := ThreePar{}

		// scan
		err = q.Scan(&d.ThreeParId, &d.Ipv4AddrManagement, &d.Ipv4AddrSSH, &d.Ipv4AddrController,
			&d.Username, &d.Password, &d.LaunchTime, &d.InCharge, &d.Contact, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &d)
	}

	return res, nil
}

func ThreeParByThreeParID(db XODB, threeParId string) (*ThreePar, error) {
	var err error

	// sql query
	const sqlStr = `SELECT ` +
		`three_par_id, ipv4_addr_management, ipv4_addr_ssh, ipv4_addr_controller, username, password, launch_time, in_charge, contact, created_at, updated_at, deleted_at, deleted ` +
		`FROM 3par ` +
		`WHERE three_par_id = ? AND deleted = 0`

	// run query
	XOLog(sqlStr, threeParId)
	d := ThreePar{}

	err = db.QueryRow(sqlStr, threeParId).Scan(&d.ThreeParId, &d.Ipv4AddrManagement, &d.Ipv4AddrSSH, &d.Ipv4AddrController,
		&d.Username, &d.Password, &d.LaunchTime, &d.InCharge, &d.Contact, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, &d.Deleted)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func ThreeParExistByThreeParId(db XODB, threeParId string) (bool, error) {
	// sql query
	const sqlStr = `SELECT * FROM 3par WHERE three_par_id = ?`

	// run query
	XOLog(sqlStr, threeParId)
	rows, err := db.Query(sqlStr, threeParId)
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
