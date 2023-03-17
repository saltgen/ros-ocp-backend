package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	database "github.com/redhatinsights/ros-ocp-backend/internal/db"
)

type Cluster struct {
	ID                uint `gorm:"primaryKey;not null;autoIncrement"`
	TenantID          uint
	RHAccount         RHAccount `gorm:"foreignKey:TenantID" json:"-"`
	ClusterID		  string
	ClusterName       string    `gorm:"type:text;unique"`
	LastReportedAt    time.Time
	LastReportedAtStr string    `gorm:"-"`
}

func (c *Cluster) AfterFind(tx *gorm.DB) error {
	c.LastReportedAtStr = c.LastReportedAt.Format(time.RFC3339)
	return nil
}

func (c *Cluster) CreateCluster() error {
	db := database.GetDB()
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "cluster_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_reported_at"}),
	}).Create(c)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
