package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/udr/logger"
	udr_context "github.com/free5gc/udr/context"
	udr_service "github.com/free5gc/udr/service"
	"github.com/free5gc/version"
	"github.com/ianchen0119/GO-CPSV/cpsv"
)

var UDR = &udr_service.UDR{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "udr"
	appLog.Infoln(app.Name)
	appLog.Infoln("UDR version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -udrcfg udr configuration file"
	app.Action = action
	app.Flags = UDR.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("UDR Run error: %v", err)
	} else {
		cpsv.Start("safCkpt=HAUDR,safApp=safCkptService")
	}
}

func stateManage() {
	udrSelf := udr_context.UDR_Self()
	for {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
		sig := <-sigs
		fmt.Println("Signal:")
		fmt.Println(sig)
		if sig == syscall.SIGUSR1 {
			fmt.Println("Swtiching to Active mode...")
			udrSelf.GetUEGroupColl()
			udrSelf.GetUESubsColl()
			udrSelf.GetSubscriptionData()
			udrSelf.GetPolicyData()
			udrSelf.GetSubscriptionID()
		} else if sig == syscall.SIGUSR2 {
			fmt.Println("Swtiching to Standby mode...")
			go udrSelf.UpdateUEGroupColl()
			go udrSelf.UpdateUESubsColl()
			go udrSelf.UpdateSubscriptionData()
			go udrSelf.UpdatePolicyData()
			go udrSelf.UpdateSubscriptionID()
		}
	}
}

func action(c *cli.Context) error {
	if err := UDR.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	go stateManage()

	UDR.Start()

	return nil
}
