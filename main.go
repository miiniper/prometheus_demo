package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/fsnotify/fsnotify"

	"github.com/miiniper/loges"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/clientcmd"

	monitorv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"gopkg.in/mgo.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ProGetA(ClusterName string) *monitorv1.PrometheusRule {

	clientSet := SetClient(ClusterName)

	proRuleInt := clientSet.PrometheusRules("checkelk-sre-k8s-loda")
	plist, err := proRuleInt.Get("tencentc-cpu", metav1.GetOptions{})
	if err != nil {
		loges.Loges.Error("get PrometheusRules list  err:", zap.Error(err))
	}

	//fmt.Println(plist)
	//	fmt.Println("=============================================================")
	//fmt.Println("ClusterName: ", plist.ClusterName)
	//fmt.Println("Name: ", plist.Name)
	//fmt.Println("Annotations: ", plist.Annotations)
	//fmt.Println("APIVersion: ", plist.APIVersion)
	fmt.Println("TypeMeta: ", plist.Kind)
	//	fmt.Println("ObjectMeta: ", plist.ObjectMeta)
	//fmt.Println("plist.Spec.Groups[0].Name    : ", plist.Spec.Groups[0].Name)
	//fmt.Println("plist.Spec.Groups[0].Rules   : ", plist.Spec.Groups[0].Rules)
	//fmt.Println("plist.Spec.Groups[0].Interval: ", plist.Spec.Groups[0].Interval)

	fmt.Println("=============================================================")
	fmt.Println("pplist.Spec.Groups[0].Rules[0].Annotations: ", plist.Spec.Groups[0].Rules[0].Annotations)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Alert       : ", plist.Spec.Groups[0].Rules[0].Alert)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Expr        : ", plist.Spec.Groups[0].Rules[0].Expr)
	fmt.Println("plist.Spec.Groups[0].Rules[0].For         : ", plist.Spec.Groups[0].Rules[0].For)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Labels      : ", plist.Spec.Groups[0].Rules[0].Labels)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Record      : ", plist.Spec.Groups[0].Rules[0].Record)
	fmt.Println("=============================================================")
	fmt.Println("plist.Spec.Groups[0].Rules[0].Expr.IntVal: ", plist.Spec.Groups[0].Rules[0].Expr.IntVal)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Expr.StrVal: ", plist.Spec.Groups[0].Rules[0].Expr.StrVal)
	fmt.Println("plist.Spec.Groups[0].Rules[0].Expr.Type  : ", plist.Spec.Groups[0].Rules[0].Expr.Type)

	return plist

}

func ProSetA(ClusterName string) {
	clientSet := SetClient(ClusterName)
	proRuleInt := clientSet.PrometheusRules("checkelk-sre-k8s-loda")

	//s1 := monitorv1.PrometheusRule{}
	//s1.Kind = "PrometheusRule"
	//s1.APIVersion = "monitoring.coreos.com/v1"
	//s1.Labels = map[string]string{"app": "prometheus-operator", "release": "prometheus-operator"}
	//s1.Name = "tencentc-mem"
	//s1.Namespace = "checkelk-sre-k8s-loda"
	//
	////s2 := monitorv1.RuleGroup{Name: "tencentc-mem",Interval: "60s",Rules:s3}
	//s1.Spec.Groups[0].Name = "tencentc-mem"
	//
	//s1.Spec.Groups[0].Interval = "60s"
	//s1.Spec.Groups[0].Rules[0].Alert = "tencentC-mem-alert"
	//s1.Spec.Groups[0].Rules[0].Annotations = map[string]string{"message": "'{{ $labels.node }}节点异常请检查 : {{ $value }}'", "runbook_url": ""}
	//s1.Spec.Groups[0].Rules[0].Expr.StrVal = "sum(kube_pod_container_resource_requests_memory_bytes{node!~\"k8snode[1-7]-.*|k8s-m.*|k8snode5[4-6]-.*|k8snode1[0-2]-.*\"}) / sum(kube_node_status_allocatable_memory_bytes{node!~\"k8snode[1-7]-.*|k8s-m.*|k8snode5[4-6]-.*|k8snode1[0-2]-.*\"}) > 0.7"
	//s1.Spec.Groups[0].Rules[0].For = "5s"
	//s1.Spec.Groups[0].Rules[0].Labels = map[string]string{"cluster": "tencent-c", "monitoring": "tencentc-mem", "namespace": "checkelk-sre-k8s-loda"}

	s1 := monitorv1.PrometheusRule{}

	plist, err := proRuleInt.Create(&s1)
	if err != nil {
		loges.Loges.Error("Set PrometheusRules  err:", zap.Error(err))
	}
	fmt.Println(plist)

}

func MakeBasicRule(ns, name string, groups []monitorv1.RuleGroup) *monitorv1.PrometheusRule {
	return &monitorv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"role": "rulefile",
			},
		},
		Spec: monitorv1.PrometheusRuleSpec{
			Groups: groups,
		},
	}
}

func generateHugePrometheusRule(ns, identifier string) *monitorv1.PrometheusRule {
	alertName := "my-alert"
	groups := []monitorv1.RuleGroup{
		{
			Name:  alertName,
			Rules: []monitorv1.Rule{},
		},
	}
	// One rule marshaled as yaml is ~34 bytes long, the max is ~524288 bytes.
	for i := 0; i < 12000; i++ {
		groups[0].Rules = append(groups[0].Rules, monitorv1.Rule{
			Alert: alertName + "-" + identifier,
			Expr:  intstr.FromString("vector(1)"),
		})
	}
	rule := MakeBasicRule(ns, "prometheus-rule-"+identifier, groups)

	return rule
}

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

	ProGetA("tencent-c")
	ProSetA("tencent-c")

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

func SetClient(ClusterName string) *monitoringv1.MonitoringV1Client {
	cfg := GetConfig(ClusterName)
	if cfg.ClusterName == "" {
		loges.Loges.Error("get cluster config error")
		return nil
	}
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.ConfigFile))
	if err != nil {
		loges.Loges.Error("REST Config From KubeConfig is err:", zap.Error(err))
		return nil
	}
	clientSet, err := monitoringv1.NewForConfig(config)
	if err != nil {
		loges.Loges.Error("NewForConfig From KubeConfig is err:", zap.Error(err))
		return nil
	}

	return clientSet
}

func GetConfig(ClusterName string) K8sConfig {
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

	for _, j := range aa {
		if j.ClusterName == ClusterName {
			return j
		}
	}

	return K8sConfig{}

}
