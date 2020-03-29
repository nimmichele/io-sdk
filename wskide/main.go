package wskide

import (
	"fmt"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/pagopa/io-sdk/wskide/icon"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// VerboseFlag is flag for verbose

var (
	// TestModeFlag enable behaviours useful for testing

	// global flags
	debugFlag = kingpin.Flag("debug", "Enable Debug log").Hidden().Default("false").Bool()
	traceFlag = kingpin.Flag("trace", "Enable Trace log").Hidden().Default("false").Bool()
	testFlag  = kingpin.Flag("test", "Enable Test Mode").Hidden().Default("false").Bool()

	// hidden global flags
	skipDockerVersion = kingpin.Flag("skip-docker-version", "Skip check of docker version").Hidden().Default("false").Bool()

	// hidden debug commands
	debugCmd        = kingpin.Command("debug", "debug").Hidden()
	ideDeployCmd    = debugCmd.Command("ide-deploy", "Create IDE deployment").Hidden()
	ideDestroyCmd   = debugCmd.Command("ide-destroy", "Destroy IDE deployment").Hidden()
	whiskDeployCmd  = debugCmd.Command("whisk-deploy", "Create Whisk deployment").Hidden()
	whiskDestroyCmd = debugCmd.Command("whisk-destroy", "Destroy Whisk deployment").Hidden()
	redisDeployCmd  = debugCmd.Command("redis-deploy", "Create Redis deployment").Hidden()
	redisDestroyCmd = debugCmd.Command("redius-destroy", "Destroy Redis deployment").Hidden()

	// start, stop, init and status
	startCmd    = kingpin.Command("start", "Start Development Enviroment")
	startDirArg = startCmd.Arg("dir", "Project dir").Required().String()
	// init
	initCmd      = kingpin.Command("init", "Initialise SDK Repository")
	initLangFlag = initCmd.Flag("language", "SDK language").Default("javascript").String()
	initDirArg   = initCmd.Arg("directory", "work directory").Required().String()
	// stop
	stopCmd = kingpin.Command("stop", "Stop Development Environment")
	// status
	statusCmd = kingpin.Command("status", "Check Containers Status")
	// config
	configCmd = kingpin.Command("config", "Edit config file $HOME/.iosdk")
	// prova
	provaCmd = kingpin.Command("prova", "Prova").Hidden()
)

func parse(cmd string) {
	switch cmd {
	// Debug
	case ideDeployCmd.FullCommand():
		IdeDeploy("")
	case ideDestroyCmd.FullCommand():
		IdeDestroy()
	case whiskDeployCmd.FullCommand():
		WhiskDeploy()
	case whiskDestroyCmd.FullCommand():
		WhiskDestroy()
	case redisDeployCmd.FullCommand():
		RedisDeploy()
	case redisDestroyCmd.FullCommand():
		RedisDestroy()
	// Start
	case startCmd.FullCommand():
		err := Start(*startDirArg)
		ShowError(err)
		if err == nil {
			time.Sleep(2 * time.Second)
			browser.OpenURL(BrowserURL)
		}
		systray.Run(onReady, onExit)

		// Stop
	case stopCmd.FullCommand():
		Stop()
	// Init
	case initCmd.FullCommand():
		ShowError(Init(*initDirArg, *initLangFlag, os.Stderr))
	// Status
	case statusCmd.FullCommand():
		dockerStatus("openwhisk")
		dockerStatus("redis")
		dockerStatus("ide-js")
	case configCmd.FullCommand():
		config()
	case provaCmd.FullCommand():
		prova()
	default:
		kingpin.Usage()
	}
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("IO SDK")
	//systray.SetTitle("IO SDK")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		Stop()
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// Sets the icon of a menu item. Only available on Mac.
	mQuitOrig.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
	os.Exit(1)
}

// Main entrypoint for wskide
func Main() {

	cmd := kingpin.Parse()
	if _, err := LoadConfig(); err != nil && cmd != configCmd.FullCommand() {
		fmt.Println("You need to run 'iosdk config', first.")
		os.Exit(1)
	}
	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	}
	if *traceFlag {
		log.SetLevel(log.TraceLevel)
	}
	parse(cmd)
}
