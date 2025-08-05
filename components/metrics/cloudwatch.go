package metrics

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"mensajesService/components/logger"
)

var cwClient *cloudwatch.Client

func Init(region string) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.LogError("Failed to init metrics", "error", err)
		return
	}
	cwClient = cloudwatch.NewFromConfig(awsCfg)
}

func PutCountMetric(metricName string, value float64) {
	if cwClient == nil {
		return
	}
	_, err := cwClient.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("MensajesService"),
		MetricData: []cloudwatchTypes.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Timestamp:  aws.Time(time.Now()),
				Value:      aws.Float64(value),
				Unit:       cloudwatchTypes.StandardUnitCount,
			},
		},
	})

	if err != nil {
		logger.LogError("Failed to put metric", "metric", metricName, "error", err)
	}
}

func PutDurationMetric(metricName string, durationMs float64) {
	if cwClient == nil {
		return
	}
	_, err := cwClient.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("MensajesService"),
		MetricData: []cloudwatchTypes.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Timestamp:  aws.Time(time.Now()),
				Value:      aws.Float64(durationMs),
				Unit:       cloudwatchTypes.StandardUnitMilliseconds,
			},
		},
	})

	if err != nil {
		logger.LogError("Failed to put duration metric", "metric", metricName, "error", err)
	}
}
