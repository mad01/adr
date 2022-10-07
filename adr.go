package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
	adrConfigFolderName   = ".adr"
	adrConfigFileName     = "config.json"
	adrConfigTemplateName = "template.md"
	adrDefaultBaseDirName = "architecture-decision-records"
)

type AdrHelper struct {
	baseDir string
}

func NewAdrHelper(baseDir string) *AdrHelper {
	helper := &AdrHelper{}
	helper.SetBaseDir(baseDir)
	return helper
}

func (a *AdrHelper) getAdrTemplateFilePath() string {
	return filepath.Join(a.baseDir, adrConfigTemplateName)
}

func (a *AdrHelper) getAdrConfigFilePath() string {
	return filepath.Join(a.baseDir, adrConfigFileName)
}

func (a *AdrHelper) SetBaseDir(dir string) {
	basedir := a.getBaseDir(dir)
	a.baseDir = basedir
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

func (a *AdrHelper) InitBaseDir(initDir string) error {
	a.baseDir = a.getBaseDir(initDir)
	if _, err := os.Stat(a.baseDir); os.IsNotExist(err) {
		os.Mkdir(a.baseDir, 0744)
	} else {
		color.Red(a.baseDir + " already exists, skipping folder creation")
	}
	return nil
}

func (a *AdrHelper) InitConfig() error {
	if _, err := os.Stat(a.baseDir); os.IsNotExist(err) {
		color.Red(a.baseDir + " did not exists, please call init")
	}
	config := AdrConfig{a.baseDir, 0}
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
	template, err := template.ParseFiles(a.getAdrTemplateFilePath())
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
	color.Green("ADR number " + strconv.Itoa(adr.Number) + " was successfully written to : " + adrFullPath)
}
