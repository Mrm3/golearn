package db_model

import "time"

type DiskWithUser struct {
	Disk
	GlobalID string `json:"global_id"`
	UserName string `json:"user_name"`
}

func UpdateDiskThreeParStatus(db XODB, status string, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, three_par_status = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, diskId)
	return err
}

func UpdateDiskStatusAvailable(db XODB, wwn string, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = 2, three_par_wwn = ? WHERE disk_id = ? AND deleted = 0", time.Now(), wwn, diskId)
	return err
}

func MarkDiskThreeParStatus(db XODB, status string, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, three_par_status = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, diskId)
	return err
}

func MarkDiskStatus(db XODB, status int8, diskId string) (err error) {
	XOLog("UPDATE disk SET updated_at = ?, status = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, diskId)
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, diskId)
	return err
}

func MarkDiskResized(db XODB, status int8, size int64, diskId string) (err error) {
	XOLog("UPDATE disk SET updated_at = ?, status = ?, size = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, size, diskId)
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = ?, size = ? WHERE disk_id = ? AND deleted = 0", time.Now(), status, size, diskId)
	return err
}

func MarkDiskDeleted(db XODB, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET deleted = 1 , deleted_at = ? WHERE disk_id = ? ", time.Now(), diskId)
	return err
}

func MarkDiskAvailable(db XODB, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = 2 WHERE disk_id = ? AND deleted = 0", time.Now(), diskId)
	return err
}

func MarkDiskInUse(db XODB, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = 7 WHERE disk_id = ? AND deleted = 0", time.Now(), diskId)
	return err
}

func MarkDiskResetting(db XODB, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, status = 8 WHERE disk_id = ? AND deleted = 0", time.Now(), diskId)
	return err
}

func DiskCreatingSnapshot(db XODB, diskId string) (count int64, err error) {
	err = db.QueryRow("SELECT Count(*) FROM snapshot WHERE disk_id = ? AND status = 1 AND deleted = 0", diskId).
		Scan(&count)

	return
}

func MarkSnapshotAvailable(db XODB, snapshotId string) (err error) {
	_, err = db.Exec("UPDATE snapshot SET updated_at = ?, status = 2 WHERE snapshot_id = ? AND deleted = 0", time.Now(), snapshotId)
	return err
}

func MarkImageAvailable(db XODB, imageId string) (err error) {
	_, err = db.Exec("UPDATE image SET updated_at = ?, status = 2 WHERE image_id = ? AND deleted = 0", time.Now(), imageId)
	return err
}

func ResizeDisk(db XODB, size uint64, diskId string) (err error) {
	_, err = db.Exec("UPDATE disk SET updated_at = ?, size = ? WHERE disk_id = ? AND deleted = 0", time.Now(), size, diskId)
	return err
}

func UpdateExportIsDeleted(db XODB, diskId, cvkName string, updateTime, deleteTime time.Time) (err error) {
	XOLog("UPDATE export SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ? AND cvk_name = ?", 1, updateTime, deleteTime, diskId, cvkName)
	_, err = db.Exec("UPDATE export SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ? AND cvk_name = ?", 1, updateTime, deleteTime, diskId, cvkName)
	return err
}

func UpdateAttachIsDeleted(db XODB, diskId, instanceId string, updateTime, deleteTime time.Time) (err error) {
	XOLog("UPDATE attach SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ? AND instance_id = ?", 1, updateTime, deleteTime, diskId, instanceId)
	_, err = db.Exec("UPDATE attach SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ? AND instance_id = ?", 1, updateTime, deleteTime, diskId, instanceId)
	return err
}

func UpdateAttachIsDeleted2(db XODB, diskId string, updateTime, deleteTime time.Time) (err error) {
	XOLog("UPDATE attach SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ?", 1, updateTime, deleteTime, diskId)
	_, err = db.Exec("UPDATE attach SET is_deleted = ?, updated_at = ?, deleted_at = ?"+
		" WHERE is_deleted = 0 AND disk_id = ?", 1, updateTime, deleteTime, diskId)
	return err
}
