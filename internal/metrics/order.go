package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(OrderList)
}

var OrderList = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "order_list_counter",
		Help: "List Order Count",
	},
	[]string{"service"},
)
