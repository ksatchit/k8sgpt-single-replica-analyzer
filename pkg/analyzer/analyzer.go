package analyzer

import (
	"context"
	"fmt"

	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer/shared"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SingleReplicaAnalyzer struct {
	client kubernetes.Interface
}

// New creates a new instance of the SingleReplicaAnalyzer
func New() analyzer.Analyzer {
	return &SingleReplicaAnalyzer{}
}

// SetClient sets the kubernetes client
func (a *SingleReplicaAnalyzer) SetClient(client kubernetes.Interface) {
	a.client = client
}

// Name returns the name of the analyzer
func (a *SingleReplicaAnalyzer) Name() string {
	return "single-replica"
}

// StartAnalysis performs the analysis of deployments
func (a *SingleReplicaAnalyzer) StartAnalysis(ctx context.Context, _ chan types.AnalyzerResult) error {
	deployments, err := a.client.AppsV1().Deployments("").List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %v", err)
	}

	results := make([]types.Result, 0)

	for _, deployment := range deployments.Items {
		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 1 {
			results = append(results, types.Result{
				Name:        deployment.Name,
				Namespace:   deployment.Namespace,
				Kind:        "Deployment",
				APIVersion: "apps/v1",
				Error:      fmt.Sprintf("Deployment %s in namespace %s has only one replica", deployment.Name, deployment.Namespace),
				Level:      types.Warning,
				Action: []string{
					fmt.Sprintf("Consider increasing the number of replicas for deployment '%s' to ensure high availability", deployment.Name),
					"You can use the following command to scale the deployment:",
					fmt.Sprintf("kubectl scale deployment %s --replicas=3 -n %s", deployment.Name, deployment.Namespace),
				},
			})
		}
	}

	if len(results) > 0 {
		resultsChan <- types.AnalyzerResult{
			Name:    a.Name(),
			Results: results,
		}
	}

	return nil
}

// Configure sets up the analyzer
func (a *SingleReplicaAnalyzer) Configure(ctx context.Context, _ shared.AnalyzerConfig) error {
	return nil
}
