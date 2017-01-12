package cmd

import (
	"github.com/urfave/cli"
	"6174/cliapp/modules/settings"
	"6174/cliapp/modules/log"
	"6174/cliapp/modules/util"
	"os"
	"path"
	"6174/cliapp/modules/scriptPackage"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/satori/go.uuid"
	"6174/cliapp/modules/compresser"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"encoding/json"
	"net/http"
)

var PublishCommand = cli.Command{
	Name: "publish",
	Usage: "Publish package",
	Description: "Publish package to remote repositry",
	Action: publishPackage,
	Flags: []cli.Flag {
		cli.StringFlag {
			Name: "registry, r",
			Value: "http://spm.idcos.com",
			Usage: "Remote registry",
		},
	},
}

func publishPackage(ctx *cli.Context) error {

	log.Info("args: %v", ctx.Args())
	settings.NewContext()
	log.Info("publish package in the current directory")
	// find the working directory
	pwd := util.PWD()
	var packageDir string
	if ctx.NArg() > 0 {
		packagePath := ctx.Args().Get(0)
		if packagePath[0] == "/"[0] {
			packageDir = packagePath
		} else {
			packageDir = path.Join(pwd, packagePath)
		}
	}
	log.Info("package dir: %s", packageDir)

	// parse pkg info
	pkgInfo, err := parsePackageYaml(packageDir)
	if err != nil {
		log.Fatal(4, "Parse pacakge error: %v", err)
		return nil
	}
	log.Info("package info: %v", pkgInfo)

	// compress pkg into temp folder
	destFilePath, err := compressToTmpFolder(packageDir, pkgInfo)

	if err != nil {
		log.Fatal(4, "Compress error: %v", err)
		return nil
	}
	log.Info("Compress pkg to: %v", destFilePath)

	// upload tar.gz file to remote registry
	registryUri := settings.RemoteRegistryUri

	customRegistry := ctx.String("registry")

	if customRegistry != "" {
		registryUri = customRegistry
	}

	err = uploadPackage(destFilePath, pkgInfo, registryUri)

	if err != nil {
		log.Fatal(4, "Upload error: %v", err)
		return nil
	}
	// remove temp file
	log.Info("end")
	return nil
}

func parsePackageYaml(filePath string) (*scriptPackage.ScriptPackage, error) {
	yamlPath := path.Join(filePath, "/package.yaml")
	var err error
	if _, err = os.Stat(yamlPath); os.IsNotExist(err) {
		return nil, err
	}

	b, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var packageInfo scriptPackage.ScriptPackage
	e := yaml.Unmarshal(b, &packageInfo)
	if e != nil {
		return nil, e
	}

	return &packageInfo, err
}

func compressToTmpFolder(packagePath string, pkgInfo *scriptPackage.ScriptPackage) (string, error){
	tempFolder := settings.AppTempDataPath
	pkgUUID := uuid.NewV1()
	destFilePath := path.Join(tempFolder, pkgUUID.String() + ".tar.gz")
	err := compresser.TarGz(packagePath, destFilePath)
	if err != nil {
		return "", err
	}
	return destFilePath, nil
}

func uploadPackage(packageTarPath string, pkgInfo *scriptPackage.ScriptPackage, registryUri string) error {
	file, err := os.Open(packageTarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(packageTarPath))

	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return err
	}

	extraParams, err := json.Marshal(pkgInfo)

	if err != nil {
		return err
	}

	writer.WriteField("pacakgeInfo", string(extraParams))
	err = writer.Close()

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", registryUri, body)

	req.Header.Set("content-Type", writer.FormDataContentType())

	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	responseBody := &bytes.Buffer{}
	_, rerr := responseBody.ReadFrom(resp.Body)

	if rerr != nil {
		return rerr
	}

	resp.Body.Close()

	log.Info("Response Status: %v", resp.StatusCode)
	log.Info("Response Header: %v", resp.Header)
	log.Info("Response Body: %v", responseBody)

	return nil
}