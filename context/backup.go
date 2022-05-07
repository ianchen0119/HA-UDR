package context

import (
	"fmt"
	"sync"
	"encoding/json"
	"github.com/free5gc/openapi/models"
	"github.com/ianchen0119/GO-CPSV/cpsv"
	"github.com/free5gc/udr/logger"
	"syscall"
	"os"
	"os/signal"
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

type Location struct {
	X int32
	Y int32
	Z int32
}

func (context *UDRContext) StartCkptSvc() {
	cpsv.Start("safCkpt=HAUDR,safApp=safCkptService")
	for {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
		sig := <-sigs
		fmt.Println("Signal:")
		fmt.Println(sig)
		if sig == syscall.SIGUSR1 {
			fmt.Println("Swtiching to Active mode...")
			context.GetUEGroupColl()
			context.GetUESubsColl()
			// context.GetSubscriptionData()
			// context.GetPolicyData()
			// context.GetSubscriptionID()

			readData, _ := cpsv.NonFixedLoad("udrSubId")
			logger.DataRepoLog.Infoln("Get Subscription ID")
			logger.DataRepoLog.Infoln(readData)
			logger.DataRepoLog.Infoln(string(readData))
			var backupIDSet BackupIDSet
			json.Unmarshal(readData, &backupIDSet)
			context.EeSubscriptionIDGenerator = backupIDSet.EeSubscriptionIDGenerator
			context.SdmSubscriptionIDGenerator = backupIDSet.SdmSubscriptionIDGenerator
			context.PolicyDataSubscriptionIDGenerator = backupIDSet.PolicyDataSubscriptionIDGenerator
			context.SubscriptionDataSubscriptionIDGenerator = backupIDSet.SubscriptionDataSubscriptionIDGenerator
		} else if sig == syscall.SIGUSR2 {
			fmt.Println("Swtiching to Standby mode...")
			// context.UpdateSubscriptionID()
			context.UpdateUEGroupColl()
			context.UpdateUESubsColl()
			// context.UpdateSubscriptionData()
			// context.UpdatePolicyData()


			var backupIDSet BackupIDSet
			backupIDSet.EeSubscriptionIDGenerator = context.EeSubscriptionIDGenerator
			backupIDSet.SdmSubscriptionIDGenerator = context.SdmSubscriptionIDGenerator
			backupIDSet.PolicyDataSubscriptionIDGenerator = context.PolicyDataSubscriptionIDGenerator
			backupIDSet.SubscriptionDataSubscriptionIDGenerator = context.SubscriptionDataSubscriptionIDGenerator

			jsonData, _ := json.Marshal(&backupIDSet)
			length := len(jsonData)

			logger.DataRepoLog.Infoln("Update Subscription ID")
			logger.DataRepoLog.Infoln(jsonData)
			cpsv.NonFixedStore("udrSubId", jsonData, int(length))
		}
	}
}

func (context *UDRContext) UpdateUEGroupColl() {
	ueGroupColl :=  make(map[interface{}]interface{})
	context.UEGroupCollection.Range(func(k, v interface{}) bool {
		ueGroupColl[k] = v
		return true
	})
	jsonData, _ := json.Marshal(ueGroupColl)
	len := len(jsonData)

	logger.DataRepoLog.Infoln("Update UEGroupColl")

	if len != 0 {
		logger.DataRepoLog.Infoln(string(jsonData))
		cpsv.NonFixedStore("UDR_UEGroupColl", jsonData, int(len))
	}
}

func (context *UDRContext) UpdateUESubsColl() {
	ueSubsColl :=  make(map[interface{}]interface{})
	context.UESubsCollection.Range(func(k, v interface{}) bool {
		ueSubsColl[k] = v
		return true
	})
	jsonData, _ := json.Marshal(ueSubsColl)
	len := len(jsonData)

	logger.DataRepoLog.Infoln("Update UESubsColl")
	
	if len != 0 {
		logger.DataRepoLog.Infoln(string(jsonData))
		cpsv.NonFixedStore("UDR_UESubsColl", jsonData, int(len))
	}
}

func (context *UDRContext) GetUEGroupColl() error {
	readData, err := cpsv.NonFixedLoad("UDR_UEGroupColl")
	if err == nil {
		logger.DataRepoLog.Infoln("Get UEGroupColl")
		tmpMap :=  make(map[interface{}]interface{})
		logger.DataRepoLog.Infoln(string(readData))
		err = json.Unmarshal(readData, &tmpMap)
		if err == nil {
			ueGroupColl := &sync.Map{}
			for k, v := range tmpMap {
				ueGroupColl.Store(k, v)
			}
			context.UEGroupCollection = *ueGroupColl
		}
		return err
	} else {
		logger.DataRepoLog.Infoln("ERR: Get UEGroupColl")
		return err
	}
}

func (context *UDRContext) GetUESubsColl() error {
	readData, err := cpsv.NonFixedLoad("UDR_UESubsColl")
	if err == nil {
		logger.DataRepoLog.Infoln("Get UESubsColl")
		tmpMap :=  make(map[interface{}]interface{})
		logger.DataRepoLog.Infoln(string(readData))
		err = json.Unmarshal(readData, &tmpMap)
		if err == nil {
			ueSubsColl := &sync.Map{}
			for k, v := range tmpMap {
				ueSubsColl.Store(k, v)
			}
			context.UESubsCollection = *ueSubsColl
		}
		return err
	} else {
		logger.DataRepoLog.Infoln("ERR: Get UESubsColl")
		return err
	}
}

func (context *UDRContext) UpdateSubscriptionData() error {

	/* TODO: convert the data */
	// ref: https://github.dev/free5gc/openapi/tree/main/Nudr_DataRepository
	// UeId string
	// CallbackReference string
	// OriginalCallbackReference string
	// MonitoredResourceUri      []string
	// Expiry                    *time.Time 
	// SupportedFeatures         string

	var subscriptionData = &SubscriptionData{}
	subscriptionData.SubscriptionDataSubscriptions = context.SubscriptionDataSubscriptions

	jsonData, _ := json.Marshal(subscriptionData)
	len := len(jsonData)

	logger.DataRepoLog.Infoln("Update Subscription Data")

	if len != 0 {
		logger.DataRepoLog.Infoln(string(jsonData))
		cpsv.NonFixedStore("UDR_SubscriptionData", jsonData, int(len))
	}

	return nil
}

func (context *UDRContext) UpdatePolicyData() error {
	/* TODO: convert the data */
	// ref: https://github.dev/free5gc/openapi/tree/main/Nudr_DataRepository
	// NotificationUri       string   `json:"notificationUri" bson:"notificationUri"`
	// MonitoredResourceUris []string `json:"monitoredResourceUris" bson:"monitoredResourceUris"`
	// SupportedFeatures     string   `json:"supportedFeatures,omitempty" bson:"supportedFeatures"`

	var policyData = &PolicyData{}
	policyData.PolicyDataSubscriptions = context.PolicyDataSubscriptions

	jsonData, _ := json.Marshal(policyData)
	len := len(jsonData)

	logger.DataRepoLog.Infoln("Update Policy Data")

	if len != 0 {
		logger.DataRepoLog.Infoln(string(jsonData))
		cpsv.NonFixedStore("UDR_PolicyData", jsonData, int(len))
	}

	return nil
}

func (context *UDRContext) UpdateSubscriptionID() error {
	var backupIDSet BackupIDSet
	backupIDSet.EeSubscriptionIDGenerator = context.EeSubscriptionIDGenerator
	backupIDSet.SdmSubscriptionIDGenerator = context.SdmSubscriptionIDGenerator
	backupIDSet.PolicyDataSubscriptionIDGenerator = context.PolicyDataSubscriptionIDGenerator
	backupIDSet.SubscriptionDataSubscriptionIDGenerator = context.SubscriptionDataSubscriptionIDGenerator

	jsonData, _ := json.Marshal(&backupIDSet)
	len := len(jsonData)

	logger.DataRepoLog.Infoln("Update Subscription ID")

	if len != 0 {
		logger.DataRepoLog.Infoln(string(jsonData))
		cpsv.NonFixedStore("udrSubId", jsonData, int(len))
	}

	return nil
}

func (context *UDRContext) GetSubscriptionData() error {
	readData, err := cpsv.NonFixedLoad("UDR_SubscriptionData")

	if err == nil {
		logger.DataRepoLog.Infoln("Get Sub Data")
		logger.DataRepoLog.Infoln(string(readData))
		var subscriptionData = SubscriptionData{}
		json.Unmarshal(readData, &subscriptionData)
		context.SubscriptionDataSubscriptions = subscriptionData.SubscriptionDataSubscriptions
		return nil
	} else {
		logger.DataRepoLog.Infoln("ERR: Get Sub Data")
		return err
	}
}

func (context *UDRContext) GetPolicyData() error {
	readData, err := cpsv.NonFixedLoad("UDR_PolicyData")

	if err == nil {
		logger.DataRepoLog.Infoln("Get Policy Data")
		logger.DataRepoLog.Infoln(string(readData))
		var policyData = PolicyData{}
		json.Unmarshal(readData, &policyData)
		context.PolicyDataSubscriptions = policyData.PolicyDataSubscriptions
		return nil
	} else {
		logger.DataRepoLog.Infoln("ERR: Get Policy Data")
		return err
	}
}

func (context *UDRContext) GetSubscriptionID() error {
	readData, err := cpsv.NonFixedLoad("udrSubId")

	if err == nil {
		logger.DataRepoLog.Infoln("Get Subscription ID")
		logger.DataRepoLog.Infoln(string(readData))
		var backupIDSet BackupIDSet
		json.Unmarshal(readData, &backupIDSet)
		context.EeSubscriptionIDGenerator = backupIDSet.EeSubscriptionIDGenerator
		context.SdmSubscriptionIDGenerator = backupIDSet.SdmSubscriptionIDGenerator
		context.PolicyDataSubscriptionIDGenerator = backupIDSet.PolicyDataSubscriptionIDGenerator
		context.SubscriptionDataSubscriptionIDGenerator = backupIDSet.SubscriptionDataSubscriptionIDGenerator
		return nil
	} else {
		logger.DataRepoLog.Infoln("ERR: Get Subscription ID")
		return err
	}
}
