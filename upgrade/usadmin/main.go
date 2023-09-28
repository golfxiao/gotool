package main

import (
	"gotool/upgrade"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/docopt/docopt-go"
)

func main() {

	usage := `Package Deploy:

Usage:
                        usadmin deploy <filename> <xml> [--config=<conf>]
                        usadmin revoke <version> [--config=<conf>]
                        usadmin -h | --help

Arguments:
                        <filename>             Will be upgradge zip package. e.g: (UC_Upgrade_PC20_2.2.428-Test.zip)
						<xml>                  Will be deployed upgrade config. e.g: (config.xml)
                        <version>              Will be upgradge to veriosn. e.g: (2.2.428)

Options:
                        -c, --config=<conf>         config path. [default:/uc/etc/usserver.conf]
                        -h, --help                  show details
`
	arg, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Printf("parser command error:%s", err.Error())
		os.Exit(-1)
	}

	if arg["--config"] == nil {
		log.Printf("--config not specified")
		os.Exit(-1)

	}
	if err := upgrade.Init(arg["--config"].(string)); err != nil {
		log.Printf("init beego framework error: %s", err.Error())
		os.Exit(-1)
	}

	log.Printf("parser argument results:%v", arg)
	rootPath := beego.AppConfig.String("package_path")
	deployCmd := NewUSDeploy(rootPath)

	var deployFile, deployXml, deployVersion string
	if arg["<filename>"] != nil {
		deployFile = arg["<filename>"].(string)
	}
	if arg["<xml>"] != nil {
		deployXml = arg["<xml>"].(string)
	}
	if arg["<version>"] != nil {
		deployVersion = arg["<version>"].(string)
	}

	switch {
	case arg["deploy"].(bool):
		err = deployCmd.Deploy(deployFile, deployXml)
	case arg["revoke"].(bool):
		err = deployCmd.Revoke(deployVersion)
	}
	if err != nil {
		log.Printf("deploy package error:%s", err.Error())
		os.Exit(-1)
	}

	log.Printf("deploy package success")
	os.Exit(0)
}
