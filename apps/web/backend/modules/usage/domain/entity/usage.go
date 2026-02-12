package entity

type UsageMetric struct {
	Used int64
}

type UsageStats struct {
	Specview UsageMetric
	Analysis UsageMetric
}
