// Package main provides the CLI for K8s Cost Prophet.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.flexinfer.ai/services/k8s-cost-prophet/pkg/types"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "prophet",
		Short:   "K8s Cost Prophet - Predictive cost analysis for Kubernetes",
		Version: version,
	}

	rootCmd.AddCommand(analyzeCmd())
	rootCmd.AddCommand(reportCmd())
	rootCmd.AddCommand(whatifCmd())
	rootCmd.AddCommand(serveCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func analyzeCmd() *cobra.Command {
	var (
		kustomize  bool
		helm       bool
		valuesFile string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "analyze [paths...]",
		Short: "Analyze manifest costs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pricing := types.DefaultPricing()

			// Example output
			report := types.CostReport{
				Period:    "month",
				TotalCost: 247.50,
				Namespaces: []types.NamespaceCost{
					{
						Namespace: "ai",
						TotalCost: 247.50,
						Workloads: []types.WorkloadCost{
							{
								Name:      "vllm-deployment",
								Namespace: "ai",
								Kind:      "Deployment",
								Resources: types.WorkloadResources{
									CPUCores:  4,
									MemoryGB:  16,
									GPUCount:  1,
									GPUType:   "amd-7900xtx",
									Replicas:  1,
								},
								HourlyCost:  0.25,
								MonthlyCost: 180.00,
								Breakdown: types.CostBreakdown{
									CPU:    14.40,
									Memory: 11.52,
									GPU:    360.00,
								},
							},
						},
					},
				},
				Summary: types.CostSummary{
					TotalWorkloads: 1,
					TotalCPUs:      4,
					TotalMemoryGB:  16,
					TotalGPUs:      1,
				},
			}

			if jsonOutput {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}

			// Human-readable output
			fmt.Println("Cost Analysis")
			fmt.Println("=============")
			fmt.Printf("Total Monthly Cost: $%.2f\n\n", report.TotalCost)
			fmt.Printf("Using pricing: CPU $%.2f/hr, Memory $%.3f/GB/hr, GPU varied\n\n", pricing.CPUHour, pricing.MemoryGBHour)

			for _, ns := range report.Namespaces {
				fmt.Printf("Namespace: %s ($%.2f/month)\n", ns.Namespace, ns.TotalCost)
				for _, w := range ns.Workloads {
					gpu := ""
					if w.Resources.GPUCount > 0 {
						gpu = fmt.Sprintf(" [GPU: %.0fx %s]", w.Resources.GPUCount, w.Resources.GPUType)
					}
					fmt.Printf("  %-30s $%7.2f/month%s\n", w.Name, w.MonthlyCost, gpu)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&kustomize, "kustomize", false, "Treat path as Kustomize directory")
	cmd.Flags().BoolVar(&helm, "helm", false, "Treat path as Helm chart")
	cmd.Flags().StringVar(&valuesFile, "values", "", "Helm values file")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}

func reportCmd() *cobra.Command {
	var (
		prometheusURL string
		period        string
		groupBy       string
		jsonOutput    bool
	)

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate cost report from live cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Generating %s cost report...\n", period)
			if prometheusURL != "" {
				fmt.Printf("Using Prometheus at %s for historical data\n", prometheusURL)
			}
			fmt.Println("(report generation not yet implemented)")
			return nil
		},
	}

	cmd.Flags().StringVar(&prometheusURL, "prometheus", "", "Prometheus URL")
	cmd.Flags().StringVar(&period, "period", "month", "Report period (hour, day, month)")
	cmd.Flags().StringVar(&groupBy, "by", "namespace", "Group by (namespace, workload)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}

func whatifCmd() *cobra.Command {
	var (
		scale       []string
		gpuType     string
		replicas    int
		spotFraction float64
	)

	cmd := &cobra.Command{
		Use:   "whatif",
		Short: "What-if scenario analysis",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("What-If Analysis")
			fmt.Println("================")

			scenario := types.WhatIfScenario{
				Name:        "Scale vllm to 2 replicas",
				CurrentCost: 247.50,
				NewCost:     495.00,
				Difference:  247.50,
				Percentage:  100.0,
				Changes: []types.ScenarioChange{
					{
						Workload:  "vllm-deployment",
						Namespace: "ai",
						Field:     "replicas",
						OldValue:  "1",
						NewValue:  "2",
					},
				},
			}

			fmt.Printf("Scenario: %s\n", scenario.Name)
			fmt.Printf("Current Cost: $%.2f/month\n", scenario.CurrentCost)
			fmt.Printf("New Cost:     $%.2f/month\n", scenario.NewCost)
			fmt.Printf("Difference:   $%.2f (+%.1f%%)\n", scenario.Difference, scenario.Percentage)

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&scale, "scale", nil, "Scale workload (name=replicas)")
	cmd.Flags().StringVar(&gpuType, "gpu-type", "", "Change GPU type")
	cmd.Flags().IntVar(&replicas, "replicas", 0, "Set replicas")
	cmd.Flags().Float64Var(&spotFraction, "spot-fraction", 0, "Fraction using spot instances")

	return cmd
}

func serveCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "mcp-serve",
		Short: "Start MCP server for cost queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Starting Cost Prophet MCP server on port %d\n", port)
			fmt.Println("(MCP server not yet implemented)")
			return nil
		},
	}

	cmd.Flags().IntVar(&port, "port", 8089, "Server port")

	return cmd
}
