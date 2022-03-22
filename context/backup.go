package context

import (
	"fmt"
	"sync"
	"encoding/json"
	"github.com/free5gc/openapi/models"
	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type BackupIDSet struct {
	EeSubscriptionIDGenerator               int
	SdmSubscriptionIDGenerator              int
	PolicyDataSubscriptionIDGenerator       int
	SubscriptionDataSubscriptionIDGenerator int
}

type BackupDataMap struct {
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
}

type BackupCollection struct {
	UESubsCollection                        sync.Map // map[ueId]*UESubsData
	UEGroupCollection                       sync.Map // map[ueGroupId]*UEGroupSubsData
}

/* TODO: update and get BackupCollection */
// ref: https://stackoverflow.com/questions/46390409/how-to-decode-json-strings-to-sync-map-instead-of-normal-map-in-go1-9

func (context *UDRContext) UpdateBackupCollection() {

}

func (context *UDRContext) GetBackupCollection() {

}

func (context *UDRContext) UpdateSubscriptionOrPolicyData() {
	var backupDataMap = &BackupDataMap{}
	backupDataMap.SubscriptionDataSubscriptions = context.SubscriptionDataSubscriptions
	backupDataMap.PolicyDataSubscriptions = context.PolicyDataSubscriptions

	jsonData, _ := json.Marshal(backupDataMap)
	len := len(jsonData)

	cpsv.NonFixedStore("UDR_SubscriptionData", jsonData, int(len))
}

func (context *UDRContext) UpdateSubscriptionID() {
	var backupIDSet = &BackupIDSet{}
	backupIDSet.EeSubscriptionIDGenerator = context.EeSubscriptionIDGenerator
	backupIDSet.SdmSubscriptionIDGenerator = context.SdmSubscriptionIDGenerator
	backupIDSet.PolicyDataSubscriptionIDGenerator = context.PolicyDataSubscriptionIDGenerator
	backupIDSet.SubscriptionDataSubscriptionIDGenerator = context.SubscriptionDataSubscriptionIDGenerator

	jsonData, _ := json.Marshal(backupIDSet)
	len := len(jsonData)

	cpsv.Store("UDR_SubscriptionID", jsonData, int(len), 0)
}

func (context *UDRContext) GetSubscriptionOrPolicyData() error {
	readData, err := cpsv.NonFixedLoad("UDR_SubscriptionData")

	if err == nil {
		fmt.Println(readData)
		var backupDataMap = BackupDataMap{}
		json.Unmarshal(readData, &backupDataMap)
		context.SubscriptionDataSubscriptions = backupDataMap.SubscriptionDataSubscriptions
		context.PolicyDataSubscriptions = backupDataMap.PolicyDataSubscriptions
		return nil
	} else {
		return err
	}
}

func (context *UDRContext) GetSubscriptionID() error {
	len := cpsv.GetSize(BackupIDSet{})
	readData, err := cpsv.Load("UDR_SubscriptionID", 0, len)

	if err == nil {
		fmt.Println(readData)
		var backupIDSet = BackupIDSet{}
		json.Unmarshal(readData, &backupIDSet)
		context.EeSubscriptionIDGenerator = backupIDSet.EeSubscriptionIDGenerator
		context.SdmSubscriptionIDGenerator = backupIDSet.SdmSubscriptionIDGenerator
		context.PolicyDataSubscriptionIDGenerator = backupIDSet.PolicyDataSubscriptionIDGenerator
		context.SubscriptionDataSubscriptionIDGenerator = backupIDSet.SubscriptionDataSubscriptionIDGenerator
		return nil
	} else {
		return err
	}
}