package config

import (
	"crypto/sha256"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

var config Config

func loadConfiguration(filepath string) (Config, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
func startlistening(updates <-chan string, errors <-chan error) {
	for {
		select {
		case filepath := <-updates:
			c, err := loadConfiguration(filepath)
			if err != nil {
				log.Println("error loading config:", err)
				continue
			}
			config = c
		case err := <-errors:
			log.Println("error watching config:", err)
		}

	}
}
func init() {
	updates, errors, err := watchConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	go startlistening(updates, errors)
}
func calculateFileHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum, err
}

func watchConfig(filepath string) (<-chan string, <-chan error, error) {
	errs := make(chan error)
	changes := make(chan string)
	hash := ""
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			newhash, err := calculateFileHash(filepath)
			if err != nil {
				errs <- err
				continue
			}
			if hash != newhash {
				hash = newhash
				changes <- filepath
			}

		}
	}()
	return changes, errs, nil
}
func watchFileSystemNotify(filepath string) (<-chan string, <-chan error, error) {
	changes := make(chan string)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	err = watcher.Add(filepath)
	if err != nil {
		return nil, nil, err
	}
	go func() {
		changes <- filepath
		for event := range watcher.Events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				changes <- event.Name
			}
		}
	}()
	return changes, watcher.Errors, nil
}
