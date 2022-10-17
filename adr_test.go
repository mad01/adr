package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testBaseDir = "just-a-test-dir"

func testHelperBaseDir() (string, string) {
	randString := func() string {
		rand.Seed(time.Now().Unix())
		charset := "abcdefghijklmnopqrstuvwxyz"

		shuff := []rune(charset)

		// Shuffling the string
		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})

		// Displaying the random string
		return string(shuff)
	}

	tempDir := fmt.Sprintf("%s/%s", testBaseDir, randString())
	os.MkdirAll(tempDir, 0744)

	readmeName := "readme.md"
	createReadme := func(testDir, readmeName string) {
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", tempDir, readmeName), []byte("#Just a readme test\n"), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	createReadme(tempDir, readmeName)

	return tempDir, readmeName
}

func TestAdrHelper_GetConfig(t *testing.T) {
	tempDir, readmeName := testHelperBaseDir()
	helper := NewAdrHelper(tempDir, readmeName)
	assert.Nil(t, helper.InitBaseDir(tempDir))
	assert.Nil(t, helper.InitConfig())
	currentConfig := helper.GetConfig()
	assert.NotNil(t, currentConfig)

}

func TestAdrHelper_InitBaseDir(t *testing.T) {
	/*
		tc := struct {
			file     string
			expected bool
		}{
			"just/a/random/dir",
			false,
		}
	*/

	tempDir, readmeName := testHelperBaseDir()
	helper := NewAdrHelper(tempDir, readmeName)
	err := helper.InitBaseDir("")
	assert.Nil(t, err)

}

func TestAdrHelper_NewAdr(t *testing.T) {
	tempDir, readme := testHelperBaseDir()

	helper := NewAdrHelper(tempDir, readme)
	assert.Nil(t, helper.InitBaseDir(tempDir))
	assert.Nil(t, helper.InitConfig())
	assert.Nil(t, helper.InitTemplate())
	currentConfig := helper.GetConfig()
	assert.NotNil(t, currentConfig)

	currentConfig.CurrentAdr++
	helper.NewAdr(currentConfig, "some name")

}
