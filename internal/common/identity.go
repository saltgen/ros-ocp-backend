package common

import (
	"encoding/json"
	"encoding/base64"

	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/ros-ocp-backend/internal/logging"
)



func GetOrgIDFromRequest(c echo.Context) (string, error) {
	log := logging.GetLogger()
	encodedIdentity := c.Request().Header.Get("X-Rh-Identity")
	decodedIdentity, err := base64.StdEncoding.DecodeString(encodedIdentity)
	if err != nil {
		log.Error("unable to ascertain identity")
		return "", err
	}

	type System struct {
		CN       string `json:"cn"`
		CertType string `json:"cert_type"`
	}

	type Internal struct {
		OrgID     string `json:"org_id"`
		AuthTime  int    `json:"auth_time"`
	}

	type Identity struct {
		OrgID      string   `json:"org_id"`
		Type       string   `json:"type"`
		AuthType   string   `json:"auth_type"`
		System     System   `json:"system"`
		Internal   Internal `json:"internal"`
	}
	
	type IdentityData struct {
		Identity Identity `json:"identity"`
	}

	var identityData IdentityData

	if err := json.Unmarshal(decodedIdentity, &identityData); err != nil {
		log.Error("unable to unmarshall identity data")
		return "", err
	}

	return identityData.Identity.OrgID, nil
}