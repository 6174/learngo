package settings

import (
	"runtime"
	"gopkg.in/ini.v1"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
	"6174/cliapp/modules/util"
	"6174/cliapp/modules/log"
	"path"
)
var (

	// App settings
	AppVersion      string
	AppPath         string
	AppName         string = "SPM"
	AppDataPath     string = "." + AppName
	AppTempDataPath string
	AppURL          string
	HomeDir         string

	// Packages settings
	LocalPackageRepoPath          string
	OSSPackageRepoEndpoint        string
 	OSSPackageRepoAccessKeyId     string
	OSSPackageRepoAccessKeySecret string
	OSSPackageRepoBucketId        string
	RemoteRegistryUri             string = "http://repositry.spm.idcos.com"

	// Server settings
	CertFile       string
	KeyFile        string
	HTTPPort       string
	Domain         string
	EnableGzip     bool

	// Security settings
	InstallLock        bool
	SecretKey          string
	LogInRememberDays  int
	MinPasswordLength  int

	// Database settings

	// OSS settings

	// UI settings

	// session settings

	// Cron tasks

	// Log settings
	LogRootPath     string

	// Global settings
	Cfg             *ini.File
	CustomPath      string
	CustomConf      string
	IsWindows       bool
	RunUser         string

)

func init() {
	IsWindows = runtime.GOOS == "windows"
	// set log
	log.NewLogger(0, "console", `{"level": 0}`)

	var err error
	if AppPath, err = getAppRootPath(); err != nil {
		log.Fatal(4, "fail to get app path: %v\n", err)
	}

	HomeDir, err = util.HomeDir()

	if err != nil {
		log.Fatal(4, "Fail to get home dir: %v\n", err)
	}

	AppDataPath = path.Join(HomeDir, "." + AppName)
	AppTempDataPath = path.Join(AppDataPath, "tmp")
	LogRootPath = path.Join(AppDataPath, "log")
	LocalPackageRepoPath = path.Join(AppDataPath, "packages")

	if _, err = os.Stat(AppDataPath); os.IsNotExist(err) {
		os.MkdirAll(AppDataPath, 0777)
		os.MkdirAll(AppTempDataPath, 0777)
		os.MkdirAll(LogRootPath, 0777)
		os.MkdirAll(LocalPackageRepoPath, 0777)
	}

	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

func getAppRootPath()(string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}


func WorkDir() (string, error) {
	wd := os.Getenv("SPM_WORK_DIR")
	if len(wd) > 0 {
		return wd, nil
	}
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

func NewContext() {
	homeDir, err := util.HomeDir()
	if err != nil {
		log.Fatal(4, "Fail to get home directory: %v", err)
	}
	homeDir = strings.Replace(homeDir, "\\", "/", -1)

	Cfg = ini.Empty()

	if err != nil {
		log.Fatal(4, "Fail to parse app.ini: %v", err)
	}

	CustomPath = os.Getenv("SPM_CUSTOM_PATH")
	if len(CustomPath) == 0 {
		CustomPath = AppDataPath;
	}

	CustomConf = CustomPath + "/app.ini"

	if util.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			log.Fatal(4, "Fail to load custom conf '%s': %v", CustomConf, err)
		}
	} else {
		log.Info("Custom config %s not found, ignore this if you are running first time", CustomConf)
	}

	sec := Cfg.Section("security")
	InstallLock = sec.Key("INSTALL_LOCK").MustBool(false)
	SecretKey = sec.Key("SECRECT_KEY").MustString("!#SPMWEMPS%.*#!")
	RunUser = Cfg.Section("").Key("RUN_USER").MustString(util.CurrentUsername())

	LogRootPath = Cfg.Section("log").Key("ROOT_PATH").MustString(path.Join(CustomPath, "log"))
	forcePathSeparator(LogRootPath)

	err = Cfg.SaveTo(CustomConf)
	if err != nil {
		log.Fatal(4, "Fail to save configure file $s : $v", CustomConf, err)
	}
}

func forcePathSeparator(path string) {
	if strings.Contains(path, "\\") {
		log.Fatal(4, "Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}