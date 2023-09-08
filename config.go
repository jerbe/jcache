package jcache

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/1 09:11
  @describe :
*/

type Config struct {
	Redis  *RedisConfig `json:"redis" yaml:"redis"`
	Memory *MemoryConfig
}

type RedisConfig struct {
	// Mode 模式
	// 支持:single,sentinel,cluster
	Mode       string   `yaml:"mode"`
	MasterName string   `yaml:"master_name"`
	Addrs      []string `yaml:"addrs"`
	Database   string   `yaml:"database"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

type MemoryConfig struct {
	// MaxKeySize 最大的键数量
	MaxKeySize int `yaml:"max_key_size"`
}
