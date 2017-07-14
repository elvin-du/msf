package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type CountLables struct {
	Caller   string
	Method   string
	RemoteIp string
	ErrCode  int64
	ErrMsg   string
}

type DurationLables struct {
	Caller   string
	Method   string
	RemoteIp string
	ErrCode  int64
	ErrMsg   string
	CostTime float64
}

var (
	rpcRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_request_count",
			Help: "rpc request count",
		},
		[]string{"caller", "method", "remoete_ip", "errcode", "errmsg"},
	)

	rpcRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "rpc_request_duration",
			Help: "rpc request duration",
		},
		[]string{"caller", "method", "remoete_ip", "errcode", "errmsg"},
	)
)

func Init() {
	prometheus.Register(rpcRequestCount)
	prometheus.Register(rpcRequestDuration)
}

func LogPrometheusCount(cl *CountLables) {
	//传入的数据顺序和数量需要和定义labels的一模一样
	counter, err := rpcRequestCount.GetMetricWithLabelValues(
		cl.Caller,
		cl.Method,
		cl.RemoteIp,
		fmt.Sprintf("%d", cl.ErrCode),
		cl.ErrMsg,
	)
	if nil != err {
		return
	}

	counter.Inc()
}

func LogPrometheusDuration(dl *DurationLables) {
	//传入的数据顺序和数量需要和定义labels的一模一样
	summary, err := rpcRequestDuration.GetMetricWithLabelValues(
		dl.Caller,
		dl.Method,
		dl.RemoteIp,
		fmt.Sprintf("%d", dl.ErrCode),
		dl.ErrMsg,
	)
	if nil != err {
		return
	}
	summary.Observe(dl.CostTime)
}
