package load

import (
	"os"
	"path/filepath"

	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xdisk"
	"github.com/imohamedsheta/xioc"
)

func InitStorage(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xdisk.Storage, error) {
		visibilityConverter := xdisk.NewDefaultVisibilityConverter()

		// local: points to the "storage/app" folder.
		// This is the base storage location for private files,
		// not directly exposed to the public web.
		localStoragePath := filepath.Join(appPath(), "storage", "app")
		localAdapter := xdisk.NewLocalAdapter(localStoragePath, visibilityConverter)

		// public: points to the "storage/app/public" folder.
		// This folder is special in Laravel convention:
		// it is symlinked to "public/storage" so that files inside
		// it can be accessed directly via the web server as static assets.
		// Used for things like user uploads, images, etc.
		publicStoragePath := filepath.Join(appPath(), "storage", "app", "public")
		publicAdapter := xdisk.NewLocalAdapter(publicStoragePath, visibilityConverter)

		// backups: points to BACKUPS_PATH env or "storage/app/backups" by default.
		// On production, set BACKUPS_PATH=/var/backups/connect so the connect
		// user can access dumps without touching the app tree.
		backupsStoragePath := backupsPath()
		backupsAdapter := xdisk.NewLocalAdapter(backupsStoragePath, visibilityConverter)

		// storage manager with both disks
		manager := xdisk.NewStorageManager("local", localAdapter)
		manager.RegisterDisk("public", publicAdapter)
		manager.RegisterDisk("backups", backupsAdapter)

		return manager, nil
	})

	if err != nil {
		x.Logger().Error("Failed to load storage module as singleton in the ioc container: " + err.Error())
	}
}

func appPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic("Error initializing Storage module: Failed to get app path from os.Getwd: " + err.Error())
	}
	return path
}

func backupsPath() string {
	if p := os.Getenv("BACKUPS_PATH"); p != "" {
		return p
	}
	return filepath.Join(appPath(), "storage", "app")
}
