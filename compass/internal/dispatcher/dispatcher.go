/*
 *
 *  Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package dispatcher

import (
	"compass/internal/configuration"
	"compass/internal/metric"
	"compass/internal/metricsgroup"
	"compass/internal/util"
	"compass/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"sync"
	"time"
)

type UseCases interface {
	Start(stopChan chan bool) error
}

type Dispatcher struct {
	metric metric.UseCases
	mux    sync.Mutex
}

var (
	metricsReachedOpts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "metrics_reached_total",
		Help: "The total of metrics reached",
	})
)

func NewDispatcher(metric metric.UseCases) UseCases {
	return &Dispatcher{metric, sync.Mutex{}}
}

func (dispatcher *Dispatcher) dispatch() {
	logger.Info("Starting metrics dispatcher", time.Now())
	metricExecutions, err := dispatcher.metric.FindAllMetricExecutions()
	if err != nil {
		logger.Panic("Cannot find active metric executions", "Dispatch", err, nil)
	}

	for _, execution := range metricExecutions {
		go dispatcher.getMetricResult(execution)
	}

	logger.Info("Finishing metrics dispatcher", time.Now())
}

func compareResultWithMetricThreshold(result float64, threshold float64, condition string) bool {
	switch condition {
	case metricsgroup.EQUAL.String():
		return result == threshold
	case metricsgroup.GREATER_THAN.String():
		return result > threshold
	case metricsgroup.LOWER_THAN.String():
		return result < threshold
	default:
		return false
	}
}

func (dispatcher *Dispatcher) getNewStatusForExecution(metricResult float64, currentMetric metric.Metric) string {
	if compareResultWithMetricThreshold(metricResult, currentMetric.Threshold, currentMetric.Condition) {
		return metric.MetricReached
	}

	return metric.MetricActive
}

func (dispatcher *Dispatcher) getMetricResult(execution metric.MetricExecution) {
	defer dispatcher.mux.Unlock()
	dispatcher.mux.Lock()

	currentMetric, err := dispatcher.metric.FindMetricById(execution.MetricID.String())
	if err != nil {
		return
	}

	metricResult, err := dispatcher.metric.ResultQuery(currentMetric)
	if err != nil {
		logger.Error(util.ResultByGroupMetricError, "getMetricResult", err, currentMetric)
		execution.Status = metric.MetricError
		dispatcher.metric.UpdateMetricExecution(execution)
		return
	}

	if metricResult != execution.LastValue || execution.Status == metric.MetricUpdated {
		dispatcher.metric.UpdateMetricExecution(metric.MetricExecution{
			BaseModel: execution.BaseModel,
			MetricID:  execution.MetricID,
			LastValue: metricResult,
			Status:    dispatcher.getNewStatusForExecution(metricResult, currentMetric),
		})

		metricsReachedOpts.Inc()
	}
}

func (dispatcher *Dispatcher) getInterval() (time.Duration, error) {
	return time.ParseDuration(configuration.GetConfiguration("DISPATCHER_INTERVAL"))
}

func (dispatcher *Dispatcher) Start(stopChan chan bool) error {
	interval, err := dispatcher.getInterval()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			dispatcher.dispatch()
		case <-stopChan:
			return nil
		}
	}
}
