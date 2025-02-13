package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/redhatinsights/ros-ocp-backend/internal/redis"
	"github.com/sirupsen/logrus"
)

var (
	log         = logrus.New()
	redisClient = redis.Client()
)

type MigrationStatus struct {
	Version uint  `json:"version"`
	Dirty   bool  `json:"dirty"`
	Err     error `json:"error"`
}

func getMigrationVersion() (uint, error) {
	/*
		Migration version is the number of migration migrationFiles applied
		Since there are _up and _down migrationFiles ~ file_count / 2
	*/
	migrationFiles, err := filepath.Glob("./migrations/*.sql")
	if err != nil {
		return 0, err
	}
	return uint(len(migrationFiles) / 2), nil
}

func getMigrationStatus(migrationInstance *migrate.Migrate, ctx context.Context) bool {
	var version uint
	var dirty bool
	var migrationErr error

	migrationStatus, err := redisClient.Get(ctx, "migration_status").Result()
	if err != nil {
		version, dirty, migrationErr = migrationInstance.Version()
	} else {
		var migrationStatusJSON MigrationStatus
		if err := json.Unmarshal([]byte(migrationStatus), &migrationStatusJSON); err != nil {
			log.Printf("failed to parse migration status data: %v", err)
			return false
		}
		log.Infof("fetched migration status from cache: %+v", migrationStatusJSON)
		version = migrationStatusJSON.Version
		dirty = migrationStatusJSON.Dirty
		migrationErr = migrationStatusJSON.Err
	}

	migrationVersion, _ := getMigrationVersion()
	switch {
	case migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNilVersion):
		log.Errorf("unable to check migration status: %v", err)
	case dirty:
		log.Error("unable to apply migrations, database is dirty")
		return false
	case version == migrationVersion:
		log.Info("migrations already applied, skipping...")
		return false
	case version < migrationVersion:
		log.Info("forward database migration")
		return true
	}
	return false
}
