package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// LoadJsonConfig is used to load config files as json format to config.
// config should be a pointer to structure, if not, panic
func LoadJsonConfig(config interface{}, fullPath string) {
	var err error
	var decoder *json.Decoder

	file := OpenFile(fullPath)
	defer file.Close()

	decoder = json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		log.WithFields(log.Fields{
			"file": fullPath,
			"err":  err,
		}).Error("Decode json fail for config file")
		return
	}

	json.Marshal(config)
}

func LoadJsonFile(fullPath string) (cfg string) {
	file := OpenFile(fullPath)
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file": fullPath,
			"err":  err,
		}).Error("Read config to string error")
		return
	}

	cfg = string(content)

	return cfg
}

func OpenFile(fullPath string) *os.File {
	var file *os.File
	var err error

	if file, err = os.Open(fullPath); err != nil {
		log.WithFields(log.Fields{
			"file": fullPath,
			"err":  err,
		}).Error("Read config to string error")
		return nil
	}

	return file
}
