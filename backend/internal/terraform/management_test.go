package terraform

import (
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"testing"
)

func TestParseAdvancedVars(t *testing.T) {
	blueNodeGroup := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"create_iam_role":            {Kind: &structpb.Value_BoolValue{BoolValue: false}},
			"iam_role_arn":               {Kind: &structpb.Value_StringValue{StringValue: "data.aws_ssm_parameter.eks_iam_node_role.value"}},
			"min_size":                   {Kind: &structpb.Value_NumberValue{NumberValue: 1}},
			"max_size":                   {Kind: &structpb.Value_NumberValue{NumberValue: 10}},
			"desired_size":               {Kind: &structpb.Value_NumberValue{NumberValue: 1}},
			"ami_id":                     {Kind: &structpb.Value_StringValue{StringValue: "ami-0c0e3c5bfa15ba56b"}},
			"instance_types":             {Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: []*structpb.Value{{Kind: &structpb.Value_StringValue{StringValue: "t3.large"}}}}}},
			"capacity_type":              {Kind: &structpb.Value_StringValue{StringValue: "SPOT"}},
			"enable_bootstrap_user_data": {Kind: &structpb.Value_BoolValue{BoolValue: true}},
		},
	}

	greenNodeGroup := &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}

	nodeGroups := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"blue":  {Kind: &structpb.Value_StructValue{StructValue: blueNodeGroup}},
			"green": {Kind: &structpb.Value_StructValue{StructValue: greenNodeGroup}},
		},
	}

	advancedValues := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"nodegroups": {Kind: &structpb.Value_StructValue{StructValue: nodeGroups}},
		},
	}

	//pbStruct := &structpb.Struct{
	//	Fields: map[string]*structpb.Value{
	//		"foo": {Kind: &structpb.Value_StringValue{StringValue: "Hello"}},
	//		"bar": {
	//			Kind: &structpb.Value_StructValue{
	//				StructValue: &structpb.Struct{
	//					Fields: map[string]*structpb.Value{
	//						"nestedKey": {Kind: &structpb.Value_NumberValue{NumberValue: 123}},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	variables := marketplace.Install_Variables{
		Values:         nil,
		NestedValues:   nil,
		AdvancedValues: advancedValues,
	}
	application := marketplace.Install_Applications{
		Name:        "",
		Version:     "",
		Variables:   &variables,
		Postinstall: "",
	}
	install := marketplace.Install{
		Applications:   &application,
		DeploymentName: "",
	}

	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()

	cloudenv := rootBody.AppendNewBlock("module", []string{"test"})
	parseAdvancedVariables(&install, cloudenv)
}
