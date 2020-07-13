package tp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// fmt.Sprintf("Ошибка при открытии файла %s: %v",filePath, err)
func CreateOpenFile(filePath string)(file *os.File,err error)  {
	file, err = os.Create(filePath)
	return
}

func OpenWriteFile(filePath string)(file *os.File,err error)   {
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	return
}

func OpenReadFile(filePath string)(dat []byte,err error)  {
	dat, err = ioutil.ReadFile(filePath)
	return
}

func OpenStandartConfig(workDir string)(map[string]string,error) {
	// Открываем конфиг, должен находиться в заданной папке
	config := make(map[string]string)
	configDir := filepath.Join(workDir, "Config.json")
	fi, err := OpenReadFile(configDir)
	if err != nil {
		return config,err
		//Ошибка при открытии файла конфигурации Config.json
	}
	err = json.Unmarshal(fi, &config)
	if err != nil {
		return config,err
	}
	return config,err
}