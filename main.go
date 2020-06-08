package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/miiniper/loges"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("server starting...")
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		loges.Loges.Info("Config file changed: ", zap.Any("", e.Name))
	})
	Init()
	ProGetA()

}

type K8sConfig struct {
	ClusterName string `json:"clustername"`
	ConfigFile  string `json:"configfile"`
}

type K8sConfigs []K8sConfig

type HttpStatus struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var ClusterCfgs K8sConfigs

func Init() {
	ClusterCfgs = GetConfig()
}

func GetConfig() K8sConfigs {
	session, err := mgo.Dial(viper.GetString("db.addr"))
	if err != nil {
		loges.Loges.Error("conn mgo is err:", zap.Error(err))
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	err = session.DB("admin").Login(viper.GetString("db.dbuser"), viper.GetString("db.dbpass"))
	if err != nil {
		loges.Loges.Error("auth mgo is err:", zap.Error(err))
	}
	aa := K8sConfigs{}
	c := session.DB("check").C("k8sconfig")
	err = c.Find(nil).All(&aa)
	if err != nil {
		loges.Loges.Error("select db is err:", zap.Error(err))
	}

	return aa
}

func K8sCli(k8sCfg string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(k8sCfg))
	if err != nil {
		loges.Loges.Error("REST Config From KubeConfig is err:", zap.Error(err))
		return nil, err
	}

	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		loges.Loges.Error("new KubeConfig is err:", zap.Error(err))
		return nil, err
	}
	return cli, nil

}

func ProGetA() {
	for _, ClusterCfg := range ClusterCfgs {
		cli, _ := K8sCli(ClusterCfg.ConfigFile)
		s1, err := cli.ServerResourcesForGroupVersion("monitoring.coreos.com/v1")
		if err != nil {
			loges.Loges.Error(" err:", zap.Error(err))
		}

		fmt.Println(s1)

	}
}
