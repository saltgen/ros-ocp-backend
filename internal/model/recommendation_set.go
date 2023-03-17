package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	database "github.com/redhatinsights/ros-ocp-backend/internal/db"
)

type RecommendationSet struct {
	ID                     uint           `gorm:"primaryKey;not null;autoIncrement"`
	WorkloadID             uint           `gorm:"type:timestamp"`
	Workload               Workload       `gorm:"foreignKey:WorkloadID"`
	MonitoringStartTime    time.Time      `gorm:"type:timestamp"`
	MonitoringEndTime      time.Time      `gorm:"type:timestamp"`
	Recommendations        datatypes.JSON `json:"recommendations"`
	CreatedAt              time.Time      `gorm:"type:timestamp"`
	MonitoringStartTimeStr string         `gorm:"-"`
	MonitoringEndTimeStr   string         `gorm:"-"`
	CreatedAtStr           string         `gorm:"-"`
}

func (r *RecommendationSet) AfterFind(tx *gorm.DB) error {
	r.MonitoringStartTimeStr = r.MonitoringStartTime.Format(time.RFC3339)
	r.MonitoringEndTimeStr = r.MonitoringEndTime.Format(time.RFC3339)
	r.CreatedAtStr = r.CreatedAt.Format(time.RFC3339)
	return nil
}

func (r *RecommendationSet) GetRecommendationSets(orgID string, limit int, offset int, queryParams map[string]interface{}) ([]RecommendationSet, error) {
	
	var recommendationSets []RecommendationSet

	query := database.DB.Joins("JOIN workloads ON recommendation_sets.workload_id = workloads.id").
		Joins("JOIN clusters ON workloads.cluster_id = clusters.id").
		Joins("JOIN rh_accounts ON clusters.tenant_id = rh_accounts.id").
		Preload("Workload.Cluster").
		Order("recommendation_sets.monitoring_start_time").
		Where("rh_accounts.org_id = ?", orgID)

	for key, value := range queryParams{
		query.Where(key, value)
	}

	err := query.Offset(offset).Limit(limit).Find(&recommendationSets).Error

	if err != nil{
		return nil, err
	}
	
	return recommendationSets, nil
}

func (r *RecommendationSet) GetRecommendationSetByID(orgID string, recommendationID int) (RecommendationSet, error) {
	
	var recommendationSet RecommendationSet

	database.DB.Joins("JOIN workloads ON recommendation_sets.workload_id = workloads.id").
		Joins("JOIN clusters ON workloads.cluster_id = clusters.id").
		Joins("JOIN rh_accounts ON clusters.tenant_id = rh_accounts.id").
		Preload("Workload.Cluster").
		Where("rh_accounts.org_id = ?", orgID).
		Where("recommendation_sets.id = ?", recommendationID).
		First(&recommendationSet)

	return recommendationSet, nil
}

func (r *RecommendationSet) CreateRecommendationSet() error {
	db := database.GetDB()
	result := db.Create(r)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
