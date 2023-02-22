package metrics

import (
	"fmt"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
	_ "time"

	"k8s.io/klog/v2"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/cloud/pkg/cloudhub/handler"
	"github.com/kubeedge/kubeedge/cloud/pkg/common/modules"
	"github.com/kubeedge/kubeedge/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// cloudierMetric 用于注册一个beehive的通讯模块
type cloudcoreMetric struct {
	enable bool
	port   uint32
}

// CloudCoreCollector 启动一个metrics模块
type CloudCoreCollector struct {
	promMetrics []*prometheus.Desc
}

// 定义metric输出的label->value
var metricValue []string

// IPLabel 添加IP的label->value
var IPLabel prometheus.Labels = make(map[string]string)

// collectDate 用于存储上一次数据
var collectDate CollectResources

// currentDate 用于输出数据
var currentDate CollectResources

func NewCollectMetric(cm *v1alpha1.Metrics) *cloudcoreMetric {
	if cm == nil {
		return &cloudcoreMetric{
			enable: false,
			port:   uint32(10082),
		}
	}
	if cm.Port == 0 {
		cm.Port = uint32(10082)
	}
	return &cloudcoreMetric{
		enable: cm.Enable,
		port:   cm.Port,
	}
}

func Register(cm *v1alpha1.Metrics) {
	mc := NewCollectMetric(cm)
	core.Register(mc)
}

func (cm *cloudcoreMetric) Name() string {
	return modules.CloudCoreMetricModuleName
}

func (cm *cloudcoreMetric) Group() string {
	return modules.CloudCoreMetricModuleName
}

func (cm *cloudcoreMetric) Enable() bool {
	return cm.enable
}

func (cm *cloudcoreMetric) Start() {
	klog.Info("Starting cloudcoreMetric dispatch message")
	go cm.updateMetric()
	if cm.Enable() {
		ip, _ := util.GetLocalIP("localhost")
		IPLabel["cloudcoreIP"] = ip
		ccm := newCloudCoreCollector()
		prometheus.MustRegister(ccm)
		http.Handle("/metrics", promhttp.Handler())
		klog.Infof("begin to server on port %d", cm.port)
		go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cm.port), nil))
	}

}

func (collector *CloudCoreCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- collector.promMetrics
}

// Collect 对收集的数据进行输出显示处理
func (collector *CloudCoreCollector) Collect(ch chan<- prometheus.Metric) {
	//get the connected count of nodes
	cncount := handler.CloudhubHandler.GetlocalConnection()
	countValue := fmt.Sprint(cncount)
	metricValue = metricValue[0:0]
	metricValue = append(metricValue, countValue)

	t := reflect.TypeOf(currentDate)
	v := reflect.ValueOf(&currentDate)
	for k := 0; k < t.NumField(); k++ {
		if v.Elem().Field(k).Kind() == reflect.Struct {
			t2 := v.Elem().Field(k).Type()
			v2 := v.Elem().Field(k)
			for i := 0; i < v2.NumField(); i++ {
				f := v2.Field(i)
				if t2.Field(i).Name == "ConsumeTime" {
					consumetime := f.Int() / 1000000
					vtostr := fmt.Sprint(consumetime)
					metricValue = append(metricValue, vtostr)
				} else {
					vtostr := fmt.Sprint(f)
					metricValue = append(metricValue, vtostr)
				}
			}
		}
	}
	for k, _ := range metricValue {
		metricV, _ := strconv.ParseFloat(metricValue[k], 64)
		ch <- prometheus.MustNewConstMetric(collector.promMetrics[k], prometheus.CounterValue, metricV, metricValue[k:k+1]...)
	}
}

func newCloudCoreCollector() *CloudCoreCollector {
	rm := getCollectType()
	var metricDesc []*prometheus.Desc
	for _, v := range rm {
		metricLabel := make([]string, 1)
		metricLabel[0] = v
		IPLabel["type"] = v
		metricDesc = append(metricDesc, prometheus.NewDesc("cloudcoreMetric", "Show cloudcoreMetric", metricLabel, IPLabel))
	}
	return &CloudCoreCollector{
		promMetrics: metricDesc,
	}
}

// getCollectType 返回metrics需要输出的label
func getCollectType() (m []string) {
	m = []string{
		"cloudcore_nodeConnectionCount",

		"cloudcore_nodeStatus_getMessageErrorCount",
		"cloudcore_nodeStatus_operateErrorCount",
		"cloudcore_nodeStatus_chanLen",
		"cloudcore_nodeStatus_ReceiveCount",
		"cloudcore_nodeStatus_consumeTime",

		"cloudcore_podStatus_getMessageErrorCount",
		"cloudcore_podStatus_operateErrorCount",
		"cloudcore_podStatus_chanLen",
		"cloudcore_podStatus_ReceiveCount",
		"cloudcore_podStatus_consumeTime",

		"cloudcore_configMap_getMessageErrorCount",
		"cloudcore_configMap_operateErrorCount",
		"cloudcore_configMap_chanLen",
		"cloudcore_configMap_ReceiveCount",
		"cloudcore_configMap_consumeTime",

		"cloudcore_secret_getMessageErrorCount",
		"cloudcore_secret_operateErrorCount",
		"cloudcore_secret_chanLen",
		"cloudcore_secret_ReceiveCount",
		"cloudcore_secret_consumeTime",

		"cloudcore_queryNode_getMessageErrorCount",
		"cloudcore_queryNode_operateErrorCount",
		"cloudcore_queryNode_chanLen",
		"cloudcore_queryNode_ReceiveCount",
		"cloudcore_queryNode_consumeTime",

		"cloudcore_updateNode_getMessageErrorCount",
		"cloudcore_updateNode_operateErrorCount",
		"cloudcore_updateNode_chanLen",
		"cloudcore_updateNode_ReceiveCount",
		"cloudcore_updateNode_consumeTime",

		"cloudcore_podDelete_getMessageErrorCount",
		"cloudcore_podDelete_operateErrorCount",
		"cloudcore_podDelete_chanLen",
		"cloudcore_podDelete_ReceiveCount",
		"cloudcore_podDelete_consumeTime",

		"cloudcore_receiveSituation_errorMessageCount",
		"cloudcore_receiveSituation_totalMessageCount",
		"cloudcore_receiveSituation_consumeTime",

		"cloudcore_sendError_SyncPodCount",
		"cloudcore_sendError_SyncConfigMap",
		"cloudcore_sendError_SyncSecret",
		"cloudcore_sendError_SyncEdgeNodes",

		"cloudcore_SendOK_SyncPodCount",
		"cloudcore_SendOK_SyncConfigMap",
		"cloudcore_SendOK_SyncSecret",
		"cloudcore_SendOK_SyncEdgeNodes",
	}
	return m
}

func (cm *cloudcoreMetric) updateMetric() {
	klog.Info("start cloudcoreMetric send message")
	for {
		select {
		case <-beehiveContext.Done():
			klog.Warning("collect metric channel event queue dispatch message loop stopped")
			return
		default:
		}
		msg, err := beehiveContext.Receive(modules.CloudCoreMetricModuleName)
		if err != nil {
			klog.Warning("the received message is not format message")
			continue
		} else {
			if value, ok := msg.Content.(CollectResources); ok {
				//数据赋值  ---  对数据进行求差
				t := reflect.TypeOf(collectDate)
				v := reflect.ValueOf(&collectDate)
				//nt := reflect.TypeOf(n)
				nv := reflect.ValueOf(&value)
				cv := reflect.ValueOf(&currentDate)
				for k := 0; k < t.NumField(); k++ {
					if v.Elem().Field(k).Kind() == reflect.Struct {
						t2 := v.Elem().Field(k).Type()
						v2 := v.Elem().Field(k)
						//nt2 := nv.Elem().Field(k).Type()
						nv2 := nv.Elem().Field(k)
						cv2 := cv.Elem().Field(k)
						for i := 0; i < v2.NumField(); i++ {
							if t2.Field(i).Name == "ChanLen" {
								chanLen := nv2.Field(i).Int()
								cv2.Field(i).SetInt(chanLen)
							} else if t2.Field(i).Name == "ConsumeTime" {
								consumeTime := nv2.Field(i).Int() - v2.Field(i).Int()
								cv2.Field(i).SetInt(consumeTime)
							} else {
								update := nv2.Field(i).Int() - v2.Field(i).Int()
								cv2.Field(i).SetInt(update)
							}
						}
					}
				}
				if currentDate.NodeStatus.ReceiveCount != 0 {
					currentDate.NodeStatus.ConsumeTime /= time.Duration(currentDate.NodeStatus.ReceiveCount)
				}
				if currentDate.PodStatus.ReceiveCount != 0 {
					currentDate.PodStatus.ConsumeTime /= time.Duration(currentDate.PodStatus.ReceiveCount)
				}
				if currentDate.ConfigMap.ReceiveCount != 0 {
					currentDate.ConfigMap.ConsumeTime /= time.Duration(currentDate.ConfigMap.ReceiveCount)
				}
				if currentDate.Secret.ReceiveCount != 0 {
					currentDate.Secret.ConsumeTime /= time.Duration(currentDate.Secret.ReceiveCount)
				}
				if currentDate.QueryNode.ReceiveCount != 0 {
					currentDate.QueryNode.ConsumeTime /= time.Duration(currentDate.QueryNode.ReceiveCount)
				}
				if currentDate.UpdateNode.ReceiveCount != 0 {
					currentDate.UpdateNode.ConsumeTime /= time.Duration(currentDate.UpdateNode.ReceiveCount)
				}
				if currentDate.PodDelete.ReceiveCount != 0 {
					currentDate.PodDelete.ConsumeTime /= time.Duration(currentDate.PodDelete.ReceiveCount)
				}
				collectDate = value
			} else {
				typeOfA := reflect.TypeOf(value)
				klog.Warning("the received message is not format message ", value, "the type of value is ", typeOfA.Name(), "and ", typeOfA.Kind())
			}
		}
	}
}

//func StartCloudMetricServer() {
//	ip, _ := util.GetLocalIP("localhost")
//	IPLabel["cloudcoreIP"] = ip
//	ccm := newCloudCoreCollector()
//	prometheus.MustRegister(ccm)
//	http.Handle("/metrics", promhttp.Handler())
//	klog.Info("begin to server on port 10082")
//	klog.Fatal(http.ListenAndServe(":10082", nil))
//}
