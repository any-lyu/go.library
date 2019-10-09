package config

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	// "github.com/spf13/viper"
)

// Viper config struct
type Viper struct {
	path           string
	onConfigChange func(fsnotify.Event)
}

// LoadConfig load config
func LoadConfig(path string, v interface{}) (err error) {
	var j []byte
	if j, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = json.Unmarshal(j, v); err != nil {
		return
	}
	return
}

// WatchConfigs watch config
func WatchConfigs(confs map[string]interface{}, watch bool) (err error) {
	for path, v := range confs {
		if err = WatchConfig(path, v, watch); err != nil {
			return
		}
	}
	return
}

// WatchConfig watch config
func WatchConfig(path string, v interface{}, watch bool) (err error) {
	viper := &Viper{path: path}
	if err = LoadConfig(path, v); err != nil {
		fmt.Printf("WatchConfig |err=%v", err)
	}
	if watch {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			// model.Log.Debugf("Config file changed: %v", e)
			if err = LoadConfig(viper.path, v); err != nil {
				fmt.Printf("OnConfigChange |LoadConfig|err=%v", e)
			}
			confJSON, err := json.Marshal(v)
			fmt.Printf("OnConfigChange |new config=%v ok=%v", string(confJSON), err)
		})
	}
	return
}

// OnConfigChange  config change listener
func (v *Viper) OnConfigChange(run func(in fsnotify.Event)) {
	v.onConfigChange = run
}

// LocalPath local path
func LocalPath(in string) (out string) {
	file, _ := exec.LookPath(in)
	out = filepath.Dir(file)
	return
}

// WatchConfig viper watch config
func (v *Viper) WatchConfig() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		configFile := filepath.Clean(v.path)
		configDir, _ := filepath.Split(configFile)

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					// we only care about the config file
					if filepath.Clean(event.Name) == configFile {
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
							v.onConfigChange(event)
						}
					}
				case err := <-watcher.Errors:
					fmt.Println("error:", err)
				}
			}
		}()

		_ = watcher.Add(configDir)
		<-done
	}()
}
