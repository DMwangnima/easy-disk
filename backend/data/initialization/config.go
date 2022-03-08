package initialization

import (
	"fmt"
	"github.com/DMwangnima/easy-disk/data/util"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	LocalConfig DataConfig
)

const (
	// 暂不更改
	DefaultBlockSize = uint64(4 * (1 << 10))
)

type DataConfig struct {
	Storage StorageConfig `yaml:"storage"`
	Server  ServerConfig  `yaml:"server"`
}

type StorageConfig struct {
	BasePath  string `yaml:"base_path"`
	BlockNum  uint64 `yaml:"block_num,omitempty"`
	ChunkSize uint64 `yaml:"chunk_size,omitempty" default:"1048576"`
}

type ServerConfig struct {
	ListenIp string `yaml:"listen_ip,omitempty" default:"0.0.0.0"`
	Port     int    `yaml:"port,omitempty" default:"9000"`
	LogPath  string `yaml:"log_path,omitempty"`
}

func (conf *StorageConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	type rawStorageConfig StorageConfig
	def := new(rawStorageConfig)
	if err = defaults.Set(def); err != nil {
		return err
	}
	if err = unmarshal(def); err != nil {
		return err
	}
	*conf = (StorageConfig)(*def)
	return nil
}

func (conf *ServerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	type rawServerConfig ServerConfig
	def := new(rawServerConfig)
	if err = defaults.Set(def); err != nil {
		return err
	}
	if err = unmarshal(def); err != nil {
		return err
	}
	*conf = (ServerConfig)(*def)
	return nil
}

func initDataConfig(configPath string) error {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config file failed, read path: %s, read err: %s", configPath, err)
	}
	if err = yaml.Unmarshal(buf, &LocalConfig); err != nil {
		return fmt.Errorf("unmarshal config file failed, err: %s", err)
	}
	if err = setDefaultBlockNum(&LocalConfig); err != nil {
		return err
	}
	if err = dumpDataConfig(configPath, &LocalConfig); err != nil {
		return err
	}
	return nil
}

func setDefaultBlockNum(config *DataConfig) error {
	if config.Storage.BlockNum != 0 {
		return nil
	}
	freeSpace, err := util.DiskFree(config.Storage.BasePath)
	if err != nil {
		return fmt.Errorf("get disk free space failed, err: %s", err)
	}
	if freeSpace < DefaultBlockSize {
		return fmt.Errorf("disk space is not enough, expect %dB, but %dB", DefaultBlockSize, freeSpace)
	}
	rawBlockNum := freeSpace / DefaultBlockSize
	// todo: adding more block allocation strategies
	blockNum := util.FloorPowerOf2(rawBlockNum)
	config.Storage.BlockNum = blockNum
	return nil
}

func dumpDataConfig(configPath string, config *DataConfig) error {
	buf, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal DataConfig failed, err: %s", err)
	}
	file, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open config file in %s failed, err: %s", configPath, err)
	}
	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("write config file in %s failed, err: %s", configPath, err)
	}
	return nil
}
