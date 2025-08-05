package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsConstants(t *testing.T) {
	// Verificar que las constantes de métricas estén definidas
	assert.NotEmpty(t, MetricMessageCreated)
	assert.NotEmpty(t, MetricMessageSuccess)
	assert.NotEmpty(t, MetricMessageError)
	assert.NotEmpty(t, MetricMessageDuration)

	assert.NotEmpty(t, MetricTimelineSuccess)
	assert.NotEmpty(t, MetricTimelineError)
	assert.NotEmpty(t, MetricTimelineDuration)

	assert.NotEmpty(t, MetricUserMessagesSuccess)
	assert.NotEmpty(t, MetricUserMessagesError)
	assert.NotEmpty(t, MetricUserMessagesDuration)

	assert.NotEmpty(t, MetricValidationError)
}

func TestPutCountMetric(t *testing.T) {
	// Test que PutCountMetric no falle
	assert.NotPanics(t, func() {
		PutCountMetric(MetricMessageSuccess, 1)
	})
}

func TestPutDurationMetric(t *testing.T) {
	// Test que PutDurationMetric no falle
	assert.NotPanics(t, func() {
		PutDurationMetric(MetricMessageDuration, 100.5)
	})
}

func TestInit(t *testing.T) {
	// Test que Init no falle
	assert.NotPanics(t, func() {
		Init("us-east-1")
	})
}
