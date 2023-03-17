package api

import (
	"strconv"
	"net/http"
	
	"github.com/labstack/echo/v4"

	"github.com/redhatinsights/ros-ocp-backend/internal/model"
	"github.com/redhatinsights/ros-ocp-backend/internal/logging"
	"github.com/redhatinsights/ros-ocp-backend/internal/common"
)

func GetRecommendationSetList(c echo.Context) error {
	OrgID, err := common.GetOrgIDFromRequest(c)
	if err != nil{
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "org_id not found"})
	}
	
	log := logging.GetLogger()

	limitStr := c.QueryParam("limit")
	limit := 10 // default value
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = limitInt
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0 // default value
	if offsetStr != "" {
		offsetInt, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = offsetInt
		}
	}

	queryParams := MapQueryParameters(c)
	recommendationSet := model.RecommendationSet{}
	recommendationSets, error := recommendationSet.GetRecommendationSets(OrgID, limit, offset, queryParams)

	if error != nil {
		log.Error("unable to fetch records from database", error)
	}

	allRecommendations := []map[string]interface{}{}

	for _, recommendation := range recommendationSets {
		recommendationData := make(map[string]interface{})
		recommendationData["id"] = recommendation.ID
		recommendationData["cluster_alias"] = recommendation.Workload.Cluster.ClusterName
		recommendationData["workload_type"] = recommendation.Workload.WorkloadType
		recommendationData["workload"] = recommendation.Workload.WorkloadName
		recommendationData["containers"] = recommendation.Workload.Containers
		recommendationData["last_report"] = recommendation.Workload.Cluster.LastReportedAtStr
		recommendationData["values"] = recommendation.Recommendations
		allRecommendations = append(allRecommendations, recommendationData)
	}

	interfaceSlice := make([]interface{}, len(allRecommendations))
	for i, v := range allRecommendations {
		interfaceSlice[i] = v
	}

	return c.JSON(http.StatusOK, CollectionResponse(interfaceSlice, c.Request(), len(allRecommendations), limit, offset))

}

func GetRecommendationSet(c echo.Context) error {
	log := logging.GetLogger()

	OrgID, err := common.GetOrgIDFromRequest(c)
	if err != nil{
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "org_id not found"})
	}
	
	// OrgID := "foo_org1"
	RecommendationID := c.Param("recommendation_id")

	ID, err := strconv.Atoi(RecommendationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "bad recommendation_id"})
	}
	recommendationSetVar := model.RecommendationSet{}
	recommendationSet, error := recommendationSetVar.GetRecommendationSetByID(OrgID, ID)

	if error != nil {
		log.Error("unable to fetch records from database", error)
	}

	recommendationSlice := make(map[string]interface{})
	recommendationSlice["id"] = recommendationSet.ID
	recommendationSlice["cluster_alias"] = recommendationSet.Workload.Cluster.ClusterName
	recommendationSlice["workload_type"] = recommendationSet.Workload.WorkloadType
	recommendationSlice["workload"] = recommendationSet.Workload.WorkloadName
	recommendationSlice["containers"] = recommendationSet.Workload.Containers
	recommendationSlice["last_report"] = recommendationSet.Workload.Cluster.LastReportedAtStr
	recommendationSlice["values"] = recommendationSet.Recommendations


	return c.JSON(http.StatusOK, recommendationSlice)
}