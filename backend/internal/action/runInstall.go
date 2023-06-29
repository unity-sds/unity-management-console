package action

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/aws"
	"github.com/unity-sds/unity-control-plane/backend/internal/processes"
	"github.com/unity-sds/unity-cs-manager/lib"
	"github.com/unity-sds/unity-cs-manager/marketplace"

	"os"
)

func RunInstall(install *marketplace.Install, conn *websocket.Conn, appConfig config.AppConfig, r processes.ActRunnerImpl) error {

	if install.Extensions != nil {
		err := spinUpExtensions(conn, appConfig, install.Extensions, r)
		if err != nil {
			return err
		}
	}

	if install.Applications != nil {
		spinUpProjects(install.Applications)
	}
	return nil
}

func spinUpProjects(applications *marketplace.Install_Applications) {

}

func spinUpExtensions(conn *websocket.Conn, appConfig config.AppConfig, extensions *marketplace.Install_Extensions, r processes.ActRunnerImpl) error {
	if extensions.Eks != nil {
		ami, err := aws.ReadSSMParameter("/unity/account/ami/eksClusterAmi")
		secgrp, err := aws.ReadSSMParameter("/unity/account/securityGroups/eksSecurityGroup")
		nodesecgrp, err := aws.ReadSSMParameter("/unity/account/securityGroups/eksSharedNodeSecurityGroup")
		pubsuba, err := aws.ReadSSMParameter("/unity/account/network/subnets/eks/publicA")
		pubsubb, err := aws.ReadSSMParameter("/unity/account/network/subnets/eks/publicB")
		privsuba, err := aws.ReadSSMParameter("/unity/account/network/subnets/eks/privateA")
		privsubb, err := aws.ReadSSMParameter("/unity/account/network/subnets/eks/privateB")
		rolearn, err := aws.ReadSSMParameter("/unity/account/roles/eksInstanceRoleArn")
		svcarn, err := aws.ReadSSMParameter("/unity/account/roles/eksServiceRoleArn")
		//userarn, err := aws.ReadSSMParameter("/unity/account/roles/mcpRoleArns")

		extensions.Eks.EKSClusterAMI = *ami.Parameter.Value
		extensions.Eks.EKSSecurityGroup = *secgrp.Parameter.Value
		extensions.Eks.EKSSharedNodeSecurityGroup = *nodesecgrp.Parameter.Value
		extensions.Eks.EKSPublicSubnetA = *pubsuba.Parameter.Value
		extensions.Eks.EKSPublicSubnetB = *pubsubb.Parameter.Value
		extensions.Eks.EKSPrivateSubnetA = *privsuba.Parameter.Value
		extensions.Eks.EKSPrivateSubnetB = *privsubb.Parameter.Value
		extensions.Eks.EKSInstanceRoleArn = *rolearn.Parameter.Value
		extensions.Eks.EKSServiceArn = *svcarn.Parameter.Value
		extensions.Eks.EKSClusterRegion = "us-west-2"
		out, err := lib.GenerateEKSTemplate(extensions.Eks)
		if err != nil {
			log.WithError(err).Error("Problem generating application metadata")
			return err
		}

		// Install package
		inputs := map[string]string{
			"DEPLOYMENTSOURCE": "act",
			"AWSCONNECTION":    "keys",
		}

		//TODO Figure out how to use packaged workflows from within act runner
		env := map[string]string{
			"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
			"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
			"AWS_REGION":            "us-west-2",
			"EKSTEMPLATE":           out,
			"CLUSTERNAME":           extensions.Eks.Clustername,
			"EKSClusterVersion":     "1.24",
		}

		secrets := map[string]string{
			"token": appConfig.GithubToken,
		}
		action := appConfig.WorkflowBasePath + "/deploy_eks_model.yml"

		return r.RunAct(action, inputs, env, secrets, conn, appConfig)

	}

	if extensions.Apigateway != nil {

		for _, n := range extensions.Apigateway.Name {
			inputs := map[string]string{}

			env := map[string]string{
				"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
				"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
				"AWS_REGION":            "us-west-2",
				"TARGET_ENV":            "mcp",
				"TARGET_STAGE":          "DEV",
				"COMMIT_HASH":           "main",
				"TARGET_PROJECT":        "UNITY",
				"TARGET_OWNER":          "nightly",
				"TARGET_API":            n,
				"TF_DIRECTORY":          "terraform-project-api-gateway_module",
			}

			secrets := map[string]string{
				"token": appConfig.GithubToken,
			}

			action := appConfig.WorkflowBasePath + "/deploy_project_apigateway.yml"
			return r.RunAct(action, inputs, env, secrets, conn, appConfig)
		}

	}

	return nil
}
