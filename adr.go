package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// AdrConfig ADR configuration, loaded and used by each sub-command
type AdrConfig struct {
	BaseDir    string `json:"base_directory"`
	CurrentAdr int    `json:"current_id"`
	ReadmeName string `json:"readme_name"`
}

// AdrEntry basic structure
type AdrEntry struct {
	Number int
	Title  string
	Date   string
	Status AdrStatus
}

// AdrStatus type
type AdrStatus string

// ADR status enums
const (
	PROPOSED   AdrStatus = "Proposed"
	ACCEPTED   AdrStatus = "Accepted"
	DEPRECATED AdrStatus = "Deprecated"
	SUPERSEDED AdrStatus = "Superseded"
)

const (
	adrConfigFileName     = "config.json"
	adrConfigTemplateName = "template.md"
	adrDefaultBaseDirName = "architecture-decision-records"
	adrDefaultReadmeName  = "Readme.md"
)

type AdrHelper struct {
	baseDir    string
	readmeName string
}

func NewAdrHelper(baseDir, readmeName string) *AdrHelper {
	helper := &AdrHelper{readmeName: readmeName}
	helper.SetBaseDir(baseDir)
	return helper
}

func (a *AdrHelper) getAdrTemplateFilePath() string {
	return filepath.Join(a.baseDir, adrConfigTemplateName)
}

func (a *AdrHelper) getAdrConfigFilePath() string {
	return filepath.Join(a.baseDir, adrConfigFileName)
}

func (a *AdrHelper) getBaseDir(initDir string) string {
	if initDir == "" {
		path, err := os.Getwd()
		if err != nil {
			color.Red("ops failed to get base dir got err: " + err.Error())
			os.Exit(1)
		}
		return fmt.Sprintf("%s/%s", path, adrDefaultBaseDirName)
	}
	return initDir

}

func (a *AdrHelper) SetBaseDir(dir string) {
	basedir := a.getBaseDir(dir)
	a.baseDir = basedir
}

func (a *AdrHelper) InitBaseDir(initDir string) error {
	a.baseDir = a.getBaseDir(initDir)
	if _, err := os.Stat(a.baseDir); os.IsNotExist(err) {
		os.Mkdir(a.baseDir, 0744)
	} else {
		color.Red(a.baseDir + " already exists, skipping folder creation")
	}

	a.SetAdrBlockInReadme(a.readmeName)
	return nil
}

func (a *AdrHelper) InitConfig() error {
	if _, err := os.Stat(a.baseDir); os.IsNotExist(err) {
		color.Red(a.baseDir + " did not exists, please call init")
	}
	config := AdrConfig{BaseDir: a.baseDir, ReadmeName: a.readmeName, CurrentAdr: 0}
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(a.getAdrConfigFilePath(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *AdrHelper) InitTemplate() error {
	body := []byte(`
# {{.Title}}

## Context

## Decision

## Status
{{.Status}}

## Consequences

`)

	err := ioutil.WriteFile(a.getAdrTemplateFilePath(), body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdrHelper) UpdateConfig(config AdrConfig) error {
	bytes, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(a.getAdrConfigFilePath(), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdrHelper) GetConfig() AdrConfig {
	var currentConfig AdrConfig

	configPath := a.getAdrConfigFilePath()
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		color.Red("No ADR configuration is found!")
		color.HiGreen("Start by initializing ADR configuration, check 'adr init --help' for more help")
		os.Exit(1)
	}

	json.Unmarshal(bytes, &currentConfig)
	return currentConfig
}

func (a *AdrHelper) NewAdr(config AdrConfig, adrName string) {
	adr := AdrEntry{
		Title:  adrName,
		Date:   time.Now().Format("02-01-2006 15:04:05"),
		Number: config.CurrentAdr,
		Status: PROPOSED,
	}
	templateFilePath := a.getAdrTemplateFilePath()
	template, err := template.ParseFiles(templateFilePath)

	if err != nil {
		panic(err)
	}
	adrFileName := strconv.Itoa(adr.Number) + "-" + strings.Join(strings.Split(strings.Trim(adr.Title, "\n \t"), " "), "-") + ".md"
	adrFullPath := filepath.Join(config.BaseDir, adrFileName)
	f, err := os.Create(adrFullPath)
	if err != nil {
		panic(err)
	}
	template.Execute(f, adr)
	f.Close()

	adrPathForReadme := fmt.Sprintf("%s/%s", a.baseDir, adrFileName)
	a.AppendRecordIndexToReadme(config.ReadmeName, adrPathForReadme, adr)
	color.Green("ADR number " + strconv.Itoa(adr.Number) + " was successfully written to : " + adrPathForReadme)

}

func (a *AdrHelper) SetAdrBlockInReadme(filename string) {
	text := fmt.Sprintf("\n##ADR index\n")
	a.AppendTextToEndOfFile(filename, text)

}

func (a *AdrHelper) AppendRecordIndexToReadme(ReadmeFilename, RecordFilename string, recordEntry AdrEntry) {
	text := fmt.Sprintf("* [%s](%s)\n", recordEntry.Title, RecordFilename)
	a.AppendTextToEndOfFile(ReadmeFilename, text)
}

func (a *AdrHelper) AppendTextToEndOfFile(filename, text string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		color.Red(filename + " did not exists")
	}

	//Append second line
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(text); err != nil {
		log.Fatal(err)
	}

}
