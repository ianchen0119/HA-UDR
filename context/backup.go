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

type SubscriptionData struct {
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
}

type PolicyData struct {
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
}

type UEGroupColl struct {
	UEGroupCollection                       sync.Map // map[ueGroupId]*UEGroupSubsData
}

type UESubsColl struct {
	UESubsCollection                        sync.Map // map[ueId]*UESubsData
}

/* TODO: update and get BackupCollection */
// ref: https://stackoverflow.com/questions/46390409/how-to-decode-json-strings-to-sync-map-instead-of-normal-map-in-go1-9

func (context *UDRContext) UpdateUEGroupColl() {

}

func (context *UDRContext) UpdateUESubsColl() {

}

func (context *UDRContext) GetUEGroupColl() {

}

func (context *UDRContext) GetUESubsColl() {

}

func (context *UDRContext) UpdateSubscriptionData() {
	var subscriptionData = &SubscriptionData{}
	subscriptionData.SubscriptionDataSubscriptions = context.SubscriptionDataSubscriptions

	jsonData, _ := json.Marshal(subscriptionData)
	len := len(jsonData)

	cpsv.NonFixedStore("UDR_SubscriptionData", jsonData, int(len))
}

func (context *UDRContext) UpdatePolicyData() {
	var policyData = &PolicyData{}
	policyData.PolicyDataSubscriptions = context.PolicyDataSubscriptions

	jsonData, _ := json.Marshal(policyData)
	len := len(jsonData)

	cpsv.NonFixedStore("UDR_PolicyData", jsonData, int(len))
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

func (context *UDRContext) GetSubscriptionData() error {
	readData, err := cpsv.NonFixedLoad("UDR_SubscriptionData")

	if err == nil {
		fmt.Println(readData)
		var subscriptionData = SubscriptionData{}
		json.Unmarshal(readData, &subscriptionData)
		context.SubscriptionDataSubscriptions = subscriptionData.SubscriptionDataSubscriptions
		return nil
	} else {
		return err
	}
}

func (context *UDRContext) GetPolicyData() error {
	readData, err := cpsv.NonFixedLoad("UDR_PolicyData")

	if err == nil {
		fmt.Println(readData)
		var policyData = PolicyData{}
		json.Unmarshal(readData, &policyData)
		context.PolicyDataSubscriptions = policyData.PolicyDataSubscriptions
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