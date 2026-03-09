// Package types defines data structures for K8s Cost Prophet.
package types

import "time"

// ResourceType represents a type of Kubernetes resource.
type ResourceType string

const (
	ResourceCPU     ResourceType = "cpu"
	ResourceMemory  ResourceType = "memory"
	ResourceGPU     ResourceType = "gpu"
	ResourceStorage ResourceType = "storage"
)

// CostModel represents a pricing model.
type CostModel string

const (
	CostModelAWS    CostModel = "aws"
	CostModelGCP    CostModel = "gcp"
	CostModelAzure  CostModel = "azure"
	CostModelCustom CostModel = "custom"
)

// PricingConfig holds cost rates.
type PricingConfig struct {
	Model             CostModel          `yaml:"model" json:"model"`
	CPUHour           float64            `yaml:"cpu_hour" json:"cpu_hour"`
	MemoryGBHour      float64            `yaml:"memory_gb_hour" json:"memory_gb_hour"`
	StorageGBMonth    float64            `yaml:"storage_gb_month" json:"storage_gb_month"`
	GPUPricing        map[string]float64 `yaml:"gpus" json:"gpus"` // GPU type -> hourly rate
	ElectricityKWH    float64            `yaml:"electricity_kwh" json:"electricity_kwh"`
	DepreciationYears int                `yaml:"depreciation_years" json:"depreciation_years"`
}

// DefaultPricing returns default pricing config.
func DefaultPricing() PricingConfig {
	return PricingConfig{
		Model:          CostModelCustom,
		CPUHour:        0.05,
		MemoryGBHour:   0.01,
		StorageGBMonth: 0.10,
		GPUPricing: map[string]float64{
			"nvidia-a100":     3.00,
			"nvidia-v100":     1.50,
			"nvidia-t4":       0.50,
			"amd-mi210":       2.00,
			"amd-7900xtx":     0.50,
			"nvidia-gtx980ti": 0.10,
		},
		ElectricityKWH:    0.12,
		DepreciationYears: 3,
	}
}

// WorkloadResources represents resources requested by a workload.
type WorkloadResources struct {
	CPUCores  float64 `json:"cpu_cores"`
	MemoryGB  float64 `json:"memory_gb"`
	GPUCount  float64 `json:"gpu_count"`
	GPUType   string  `json:"gpu_type,omitempty"`
	StorageGB float64 `json:"storage_gb"`
	Replicas  int     `json:"replicas"`
}

// WorkloadCost represents the cost breakdown for a workload.
type WorkloadCost struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Kind        string            `json:"kind"` // Deployment, StatefulSet, etc.
	Resources   WorkloadResources `json:"resources"`
	HourlyCost  float64           `json:"hourly_cost"`
	MonthlyCost float64           `json:"monthly_cost"`
	Breakdown   CostBreakdown     `json:"breakdown"`
}

// CostBreakdown shows cost by resource type.
type CostBreakdown struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	GPU     float64 `json:"gpu"`
	Storage float64 `json:"storage"`
}

// NamespaceCost aggregates costs by namespace.
type NamespaceCost struct {
	Namespace string         `json:"namespace"`
	Workloads []WorkloadCost `json:"workloads"`
	TotalCost float64        `json:"total_cost"`
	Breakdown CostBreakdown  `json:"breakdown"`
}

// CostReport is a complete cost analysis.
type CostReport struct {
	GeneratedAt time.Time       `json:"generated_at"`
	Period      string          `json:"period"` // hour, day, month
	TotalCost   float64         `json:"total_cost"`
	Namespaces  []NamespaceCost `json:"namespaces"`
	Summary     CostSummary     `json:"summary"`
}

// CostSummary provides high-level cost insights.
type CostSummary struct {
	TotalWorkloads int            `json:"total_workloads"`
	TotalCPUs      float64        `json:"total_cpus"`
	TotalMemoryGB  float64        `json:"total_memory_gb"`
	TotalGPUs      float64        `json:"total_gpus"`
	TotalStorageGB float64        `json:"total_storage_gb"`
	Breakdown      CostBreakdown  `json:"breakdown"`
	TopWorkloads   []WorkloadCost `json:"top_workloads"` // Top 5 by cost
}

// WhatIfScenario represents a hypothetical change.
type WhatIfScenario struct {
	Name        string           `json:"name"`
	Changes     []ScenarioChange `json:"changes"`
	CurrentCost float64          `json:"current_cost"`
	NewCost     float64          `json:"new_cost"`
	Difference  float64          `json:"difference"`
	Percentage  float64          `json:"percentage"`
}

// ScenarioChange represents a single change in a what-if scenario.
type ScenarioChange struct {
	Workload  string `json:"workload"`
	Namespace string `json:"namespace"`
	Field     string `json:"field"` // replicas, cpu, memory, gpu_type
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
}

// Config holds the complete configuration.
type Config struct {
	Pricing    PricingConfig    `yaml:"pricing" json:"pricing"`
	Prometheus PrometheusConfig `yaml:"prometheus" json:"prometheus"`
	Thresholds ThresholdConfig  `yaml:"thresholds" json:"thresholds"`
	Namespaces NamespaceConfig  `yaml:"namespaces" json:"namespaces"`
}

// PrometheusConfig for historical data.
type PrometheusConfig struct {
	URL     string            `yaml:"url" json:"url"`
	Metrics map[string]string `yaml:"metrics" json:"metrics"`
}

// ThresholdConfig for alerts.
type ThresholdConfig struct {
	WarnMonthly  float64 `yaml:"warn_monthly" json:"warn_monthly"`
	WarnIncrease float64 `yaml:"warn_increase" json:"warn_increase"` // Percentage
}

// NamespaceConfig for filtering.
type NamespaceConfig struct {
	Exclude []string `yaml:"exclude" json:"exclude"`
	Include []string `yaml:"include" json:"include"`
}
