package helper

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

const configPath = "/etc/yig/yig-billing.toml"

var Conf Config

type Config struct {
	PrometheusUrl         string        `toml:"prometheus_url"`
	RegionId              string        `toml:"region_id"`
	UsageDataDir          string        `toml:"usage_data_dir"`
	BucketUsageDataDir    string        `toml:"bucket_usage_data_dir"`
	LogPath               string        `toml:"log_path"`
	EnablePostBillingCron bool          `toml:"enable_post_billing_cron"`
	PostBillingSpec       string        `toml:"post_billing_spec"`
	SparkHome             string        `toml:"spark_home"`
	TisparkShell          string        `toml:"tispark_shell"`
	TisparkShellBucket    string        `toml:"tispark_shell_bucket"`
	Producer              DummyProducer `toml:"producer"`
	RedisUrl              string        `toml:"redis_url"`
	RedisPassword         string        `toml:"redis_password"`
	EnableUsageCache      bool          `toml:"enable_usage_cache"`
	CleanUpDatabaseSpec   string        `toml:"clean_up_database_spec"`
	TidbConnection        string        `toml:"tidb_connection"`
	DbMaxIdleConns        int           `toml:"db_max_idle_conns"`
	DbMaxOpenConns        int           `toml:"db_max_open_conns"`
	DbConnMaxLifeSeconds  int           `toml:"db_conn_max_life_seconds"`
	KafkaServer           string        `toml:"kafka_server"`
	KafkaGroupId          string        `toml:"kafka_group_id"`
	KafkaTopic            string        `toml:"kafka_topic"`
}

func ReadConfig() {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		if err != nil {
			panic("[ERROR] Cannot open /etc/yig/yig-billing.toml")
		}
	}
	_, err = toml.Decode(string(data), &Conf)
	if err != nil {
		panic("[ERROR] Load yig-billing.toml error: " + err.Error())
	}
}
