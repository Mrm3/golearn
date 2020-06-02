package ceph

import (
	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"immortality-demo/config"
	"immortality-demo/pkg/logger"
	"immortality-demo/pkg/util"
)

type Pool string

const (
	DiskPool  Pool = "disk"
	ImagePool Pool = "image"
)

type ceph struct {
	ID       string
	conn     *rados.Conn
	diskCtx  *rados.IOContext
	imageCtx *rados.IOContext
}

func newCeph(conf config.CephConfig) (c ceph, err error) {
	c = ceph{}

	c.conn, err = rados.NewConnWithUser(conf.User)
	if err != nil {
		logger.Log.Error("ceph failed:", err)
		return
	}

	err = c.conn.ReadConfigFile(conf.ConfPath)
	if err != nil {
		logger.Log.Error("ceph failed(", conf.ConfPath, "):", err)
		return
	}

	logger.Log.Info("ceph config", conf.ConfPath, err)

	err = c.conn.Connect()
	if err != nil {
		logger.Log.Error("ceph failed:", err)
		return
	}

	c.ID, err = c.conn.GetFSID()
	if err != nil {
		logger.Log.Error("ceph failed:", err)
		return
	}

	c.diskCtx, err = c.conn.OpenIOContext(string(DiskPool))
	if err != nil {
		logger.Log.Error("ceph failed(rbd):", err)
		return
	}

	c.imageCtx, err = c.conn.OpenIOContext(string(ImagePool))
	if err != nil {
		logger.Log.Error("ceph failed(image):", err)
		return
	}

	return
}

func (c ceph) close() {
	c.diskCtx.Destroy()
	c.imageCtx.Destroy()
	c.conn.Shutdown()
}

func (c ceph) contextByPool(pool Pool) *rados.IOContext {
	if pool == DiskPool {
		return c.diskCtx
	}
	if pool == ImagePool {
		return c.imageCtx
	}
	return nil
}

func (c ceph) get(pool Pool, name string) (image *rbd.Image, err error) {
	ioContext := c.contextByPool(pool)
	image = rbd.GetImage(ioContext, name)
	err = image.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = image.Close()
	if err != nil {
		logger.Log.Error(err)
	}
	return
}

func (c ceph) create(pool Pool, name string, size uint64) (image *rbd.Image, err error) {
	ioContext := c.contextByPool(pool)
	image, err = rbd.Create(ioContext, name, size, 22, 3)
	if err != nil {
		logger.Log.Error(err)
	}
	return
}

func (c ceph) resizeImage(pool Pool, name string, size uint64) (err error) {
	image, err := c.get(pool, name)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	err = image.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}

	defer func() {
		_ = image.Close()
	}()

	err = image.Resize(size)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	return
}

func (c ceph) removeImage(pool Pool, name string) (err error) {
	image, err := c.get(pool, name)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = image.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer image.Close()
	snapshots, err := image.GetSnapshotNames()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	for _, snapshot := range snapshots {
		snapshot := image.GetSnapshot(snapshot.Name)
		_ = snapshot.Unprotect()
		err = snapshot.Remove()
		if err != nil {
			logger.Log.Error(err)
			return
		}
	}
	err = image.Close()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = image.Remove()
	if err != nil {
		logger.Log.Error(err)
	}
	return
}

func (c ceph) createDiskFromSnapshot(diskName, snapName, destName string) (err error) {
	var disk, dest *rbd.Image
	disk, err = c.get(DiskPool, diskName)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	dest, err = disk.Clone(snapName, c.diskCtx, destName, 3, 22)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = dest.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer func() {
		_ = dest.Close()
		if err != nil {
			_ = c.removeImage(DiskPool, destName)
		}
	}()
	err = dest.Flatten()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	return nil
}

func (c ceph) createSnapshot(diskName, destName string) (snap *rbd.Snapshot, err error) {
	var disk *rbd.Image

	disk, err = c.get(DiskPool, diskName)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	err = disk.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}

	defer disk.Close()

	snap, err = disk.CreateSnapshot(destName)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	err = snap.Protect()
	if err != nil {
		logger.Log.Error(err)
	}

	return
}

func (c ceph) rollbackSnapshot(diskName, snapName string) (err error) {
	var disk *rbd.Image
	disk, err = c.get(DiskPool, diskName)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = disk.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer disk.Close()
	snap := disk.GetSnapshot(snapName)

	err = snap.Rollback()
	if err != nil {
		logger.Log.Error(err)
	}
	return
}

func (c ceph) removeSnapshot(diskName, snapName string) (err error) {
	var disk *rbd.Image
	disk, err = c.get(DiskPool, diskName)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = disk.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer disk.Close()
	var snap *rbd.Snapshot
	snap = disk.GetSnapshot(snapName)
	var protected bool
	protected, err = snap.IsProtected()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	if protected {
		err = snap.Unprotect()
		if err != nil {
			logger.Log.Error(err)
			return
		}
	}
	err = snap.Remove()
	if err != nil {
		logger.Log.Error(err)
	}
	return
}

func (c ceph) createImageSameCluster(diskName, imageName string) (err error) {
	sourceDisk, err := c.get(DiskPool, diskName)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = sourceDisk.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer sourceDisk.Close()
	snapshotId := util.GenerateRandomId()
	sourceSnapshot, err := sourceDisk.CreateSnapshot(snapshotId)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = sourceSnapshot.Protect()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer func() {
		_ = sourceSnapshot.Unprotect()
		_ = sourceSnapshot.Remove()
	}()
	clonedImage, err := sourceDisk.Clone(snapshotId, c.imageCtx,
		imageName, 3, 22)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = clonedImage.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer func() {
		if err != nil {
			_ = clonedImage.Remove()
		}
		_ = clonedImage.Close()
	}()
	return clonedImage.Flatten()
}

func (c ceph) createDiskFromImage(imageType, imageName, diskName string) (err error) {
	sourceImage := rbd.GetImage(c.imageCtx, imageName)
	err = sourceImage.Open()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	defer func() {
		_ = sourceImage.Close()
	}()
	snap, err := sourceImage.CreateSnapshot(imageName + "-snapshot")
	if err != nil {
		logger.Log.Warn("CreateSnapshot error:", err,
			"(it's safe to ignore this warn)")
	} else {
		_ = snap.Protect()
	}
	snap = sourceImage.GetSnapshot(imageName + "-snapshot")
	dest, err := sourceImage.Clone(imageName+"-snapshot",
		c.diskCtx, diskName, 3, 22)
	if err != nil {
		logger.Log.Error(err)
	}

	if imageType == "custom" {
		err = dest.Open()
		if err != nil {
			logger.Log.Error(err)
			return
		}
		defer dest.Close()
		err = dest.Flatten()
		if err != nil {
			logger.Log.Error(err)
		}
	}
	return
}
