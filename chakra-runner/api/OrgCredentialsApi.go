package api

import (
	"chakrarunner/pojo"

	"gorm.io/gorm"
)

func GetOrgCredentials(DbConnection *gorm.DB, orgId int) *[]pojo.OrgCredentials {
	orgCredentials := &[]pojo.OrgCredentials{}
	DbConnection.Where("org_id = ?", orgId).Find(&orgCredentials)
	return orgCredentials
}

func SaveOrgCredentials(db *gorm.DB, orgId int, credentialMap map[string]string) error {
	credentialMapList := make([]pojo.OrgCredentials, 0, len(credentialMap))
	for k := range credentialMap {
		credentialMapList = append(credentialMapList, pojo.OrgCredentials{HeaderParam: k, HeaderValue: credentialMap[k], OrgId: orgId})
	}
	return db.CreateInBatches(credentialMapList, len(credentialMapList)).Error
}

func DeleteOrgCredentials(db *gorm.DB, orgId int) error {
	return db.Where("org_id = ?", orgId).Delete(&pojo.OrgCredentials{}).Error
}
