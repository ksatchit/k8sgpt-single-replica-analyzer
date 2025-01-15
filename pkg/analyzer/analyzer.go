package analyzer

import (
	"context"
	"fmt"

	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	//appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SingleReplicaAnalyzer struct {
	client kubernetes.Interface
}

func New() common.IAnalyzer {
	return &SingleReplicaAnalyzer{}
}

func (a *SingleReplicaAnalyzer) SetClient(client kubernetes.Interface) {
	a.client = client
}

func (a *SingleReplicaAnalyzer) Name() string {
	return "single-replica"
}

func (a *SingleReplicaAnalyzer) Analyze(analyzerCtx common.Analyzer) ([]common.Result, error) {
	ctx := context.Background()
	deployments, err := a.client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}

	var results []common.Result

	for _, deployment := range deployments.Items {
		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 1 {
			results = append(results, common.Result{
				Kind: "Deployment",
				Name: deployment.Name,
				Error: []common.Failure{
					{
						Text: fmt.Sprintf("Deployment %s in namespace %s has only one replica", deployment.Name, deployment.Namespace),
					},
					{
						Text: fmt.Sprintf("Consider increasing the number of replicas for deployment '%s' to ensure high availability", deployment.Name),
					},
					{
						Text: "You can use the following command to scale the deployment:",
					},
					{
						Text: fmt.Sprintf("kubectl scale deployment %s --replicas=3 -n %s", deployment.Name, deployment.Namespace),
					},
				},
			})
		}
	}

	return results, nil
}

func (a *SingleReplicaAnalyzer) Configure(ctx context.Context, config map[string]string) error {
	return nil
}
