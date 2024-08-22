package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/containerresources"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func helmChart(ctx *pulumi.Context,
	locals *Locals, createdNamespace *kubernetescorev1.Namespace) error {

	// https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml
	var helmValues = pulumi.Map{
		"fullnameOverride": pulumi.String(locals.OpenfgaKubernetes.Metadata.Name),
		"replicaCount":     pulumi.Int(1),
		"datastore": pulumi.Map{
			"engine": pulumi.String("postgres"),
			"uri":    pulumi.String("postgres://postgres:@pgk8s-planton-cloud-prod-t20240802.planton.live:5432/db_openfga?sslmode=require"),
		},
		"resources": containerresources.ConvertToPulumiMap(locals.OpenfgaKubernetes.Spec.Container.Resources),
		//"postgresql": pulumi.Map{
		//	"enabled": pulumi.Bool(true),
		//	"auth": pulumi.Map{
		//		"postgresPassword": pulumi.String("RCKS63Qtt5CvyG7U1x2fXvnF4mxqhixPjm05Xv2YFknI3ULNuqZLWhTezgqbG0tJ"),
		//		"database":         pulumi.String("postgres"),
		//	},
		//},
	}

	//merge extra helm values provided in the spec with base values
	//mergemaps.MergeMapToPulumiMap(helmValues, locals.OpenfgaKubernetes.Spec.HelmValues)

	//install jenkins helm-chart
	_, err := helmv3.NewChart(ctx,
		locals.OpenfgaKubernetes.Metadata.Id,
		helmv3.ChartArgs{
			Chart:     pulumi.String(vars.HelmChartName),
			Version:   pulumi.String(vars.HelmChartVersion),
			Namespace: createdNamespace.Metadata.Name().Elem(),
			Values:    helmValues,
			FetchArgs: helmv3.FetchArgs{
				Repo: pulumi.String(vars.HelmChartRepoUrl),
			},
		}, pulumi.Parent(createdNamespace), pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "3m", Update: "3m", Delete: "3m"}))
	if err != nil {
		return errors.Wrap(err, "failed to create helm chart")
	}

	return nil
}
