// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.22.3
// source: marketplace.proto

package marketplace

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MarketplaceMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                string                                     `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Version             string                                     `protobuf:"bytes,2,opt,name=Version,proto3" json:"Version,omitempty"`
	Channel             string                                     `protobuf:"bytes,3,opt,name=Channel,proto3" json:"Channel,omitempty"`
	Owner               string                                     `protobuf:"bytes,4,opt,name=Owner,proto3" json:"Owner,omitempty"`
	Description         string                                     `protobuf:"bytes,5,opt,name=Description,proto3" json:"Description,omitempty"`
	Repository          string                                     `protobuf:"bytes,6,opt,name=Repository,proto3" json:"Repository,omitempty"`
	Tags                []string                                   `protobuf:"bytes,7,rep,name=Tags,proto3" json:"Tags,omitempty"`
	Category            string                                     `protobuf:"bytes,8,opt,name=Category,proto3" json:"Category,omitempty"`
	IamRoles            *MarketplaceMetadata_Iamroles              `protobuf:"bytes,9,opt,name=IamRoles,proto3" json:"IamRoles,omitempty"`
	Package             string                                     `protobuf:"bytes,10,opt,name=Package,proto3" json:"Package,omitempty"`
	ManagedDependencies []*MarketplaceMetadata_Manageddependencies `protobuf:"bytes,11,rep,name=ManagedDependencies,proto3" json:"ManagedDependencies,omitempty"`
	Backend             string                                     `protobuf:"bytes,12,opt,name=Backend,proto3" json:"Backend,omitempty"`
	DefaultDeployment   *MarketplaceMetadata_Defaultdeployment     `protobuf:"bytes,13,opt,name=DefaultDeployment,proto3" json:"DefaultDeployment,omitempty"`
}

func (x *MarketplaceMetadata) Reset() {
	*x = MarketplaceMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata) ProtoMessage() {}

func (x *MarketplaceMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0}
}

func (x *MarketplaceMetadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MarketplaceMetadata) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *MarketplaceMetadata) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *MarketplaceMetadata) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *MarketplaceMetadata) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *MarketplaceMetadata) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *MarketplaceMetadata) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *MarketplaceMetadata) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *MarketplaceMetadata) GetIamRoles() *MarketplaceMetadata_Iamroles {
	if x != nil {
		return x.IamRoles
	}
	return nil
}

func (x *MarketplaceMetadata) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

func (x *MarketplaceMetadata) GetManagedDependencies() []*MarketplaceMetadata_Manageddependencies {
	if x != nil {
		return x.ManagedDependencies
	}
	return nil
}

func (x *MarketplaceMetadata) GetBackend() string {
	if x != nil {
		return x.Backend
	}
	return ""
}

func (x *MarketplaceMetadata) GetDefaultDeployment() *MarketplaceMetadata_Defaultdeployment {
	if x != nil {
		return x.DefaultDeployment
	}
	return nil
}

type MarketplaceMetadata_Statement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Effect   string   `protobuf:"bytes,1,opt,name=Effect,proto3" json:"Effect,omitempty"`
	Action   []string `protobuf:"bytes,2,rep,name=Action,proto3" json:"Action,omitempty"`
	Resource []string `protobuf:"bytes,3,rep,name=Resource,proto3" json:"Resource,omitempty"`
}

func (x *MarketplaceMetadata_Statement) Reset() {
	*x = MarketplaceMetadata_Statement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Statement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Statement) ProtoMessage() {}

func (x *MarketplaceMetadata_Statement) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Statement.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Statement) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MarketplaceMetadata_Statement) GetEffect() string {
	if x != nil {
		return x.Effect
	}
	return ""
}

func (x *MarketplaceMetadata_Statement) GetAction() []string {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *MarketplaceMetadata_Statement) GetResource() []string {
	if x != nil {
		return x.Resource
	}
	return nil
}

type MarketplaceMetadata_Iamroles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statement []*MarketplaceMetadata_Statement `protobuf:"bytes,1,rep,name=Statement,proto3" json:"Statement,omitempty"`
}

func (x *MarketplaceMetadata_Iamroles) Reset() {
	*x = MarketplaceMetadata_Iamroles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Iamroles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Iamroles) ProtoMessage() {}

func (x *MarketplaceMetadata_Iamroles) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Iamroles.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Iamroles) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 1}
}

func (x *MarketplaceMetadata_Iamroles) GetStatement() []*MarketplaceMetadata_Statement {
	if x != nil {
		return x.Statement
	}
	return nil
}

type MarketplaceMetadata_Eks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinimumVersion string `protobuf:"bytes,1,opt,name=MinimumVersion,proto3" json:"MinimumVersion,omitempty"`
}

func (x *MarketplaceMetadata_Eks) Reset() {
	*x = MarketplaceMetadata_Eks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Eks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Eks) ProtoMessage() {}

func (x *MarketplaceMetadata_Eks) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Eks.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Eks) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 2}
}

func (x *MarketplaceMetadata_Eks) GetMinimumVersion() string {
	if x != nil {
		return x.MinimumVersion
	}
	return ""
}

type MarketplaceMetadata_Manageddependencies struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eks *MarketplaceMetadata_Eks `protobuf:"bytes,1,opt,name=Eks,proto3" json:"Eks,omitempty"`
}

func (x *MarketplaceMetadata_Manageddependencies) Reset() {
	*x = MarketplaceMetadata_Manageddependencies{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Manageddependencies) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Manageddependencies) ProtoMessage() {}

func (x *MarketplaceMetadata_Manageddependencies) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Manageddependencies.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Manageddependencies) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 3}
}

func (x *MarketplaceMetadata_Manageddependencies) GetEks() *MarketplaceMetadata_Eks {
	if x != nil {
		return x.Eks
	}
	return nil
}

type MarketplaceMetadata_Variables struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SomeTerraformVariable string `protobuf:"bytes,1,opt,name=some_terraform_variable,json=someTerraformVariable,proto3" json:"some_terraform_variable,omitempty"`
}

func (x *MarketplaceMetadata_Variables) Reset() {
	*x = MarketplaceMetadata_Variables{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Variables) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Variables) ProtoMessage() {}

func (x *MarketplaceMetadata_Variables) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Variables.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Variables) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 4}
}

func (x *MarketplaceMetadata_Variables) GetSomeTerraformVariable() string {
	if x != nil {
		return x.SomeTerraformVariable
	}
	return ""
}

type MarketplaceMetadata_Nodegroup1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinNodes     uint32 `protobuf:"varint,1,opt,name=MinNodes,proto3" json:"MinNodes,omitempty"`
	MaxNodes     uint32 `protobuf:"varint,2,opt,name=MaxNodes,proto3" json:"MaxNodes,omitempty"`
	DesiredNodes uint32 `protobuf:"varint,3,opt,name=DesiredNodes,proto3" json:"DesiredNodes,omitempty"`
	InstanceType string `protobuf:"bytes,4,opt,name=InstanceType,proto3" json:"InstanceType,omitempty"`
}

func (x *MarketplaceMetadata_Nodegroup1) Reset() {
	*x = MarketplaceMetadata_Nodegroup1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Nodegroup1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Nodegroup1) ProtoMessage() {}

func (x *MarketplaceMetadata_Nodegroup1) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Nodegroup1.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Nodegroup1) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 5}
}

func (x *MarketplaceMetadata_Nodegroup1) GetMinNodes() uint32 {
	if x != nil {
		return x.MinNodes
	}
	return 0
}

func (x *MarketplaceMetadata_Nodegroup1) GetMaxNodes() uint32 {
	if x != nil {
		return x.MaxNodes
	}
	return 0
}

func (x *MarketplaceMetadata_Nodegroup1) GetDesiredNodes() uint32 {
	if x != nil {
		return x.DesiredNodes
	}
	return 0
}

func (x *MarketplaceMetadata_Nodegroup1) GetInstanceType() string {
	if x != nil {
		return x.InstanceType
	}
	return ""
}

type MarketplaceMetadata_Nodegroups struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeGroup1 *MarketplaceMetadata_Nodegroup1 `protobuf:"bytes,1,opt,name=NodeGroup1,proto3" json:"NodeGroup1,omitempty"`
}

func (x *MarketplaceMetadata_Nodegroups) Reset() {
	*x = MarketplaceMetadata_Nodegroups{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Nodegroups) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Nodegroups) ProtoMessage() {}

func (x *MarketplaceMetadata_Nodegroups) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Nodegroups.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Nodegroups) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 6}
}

func (x *MarketplaceMetadata_Nodegroups) GetNodeGroup1() *MarketplaceMetadata_Nodegroup1 {
	if x != nil {
		return x.NodeGroup1
	}
	return nil
}

type MarketplaceMetadata_Eksspec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeGroups []*MarketplaceMetadata_Nodegroups `protobuf:"bytes,1,rep,name=NodeGroups,proto3" json:"NodeGroups,omitempty"`
}

func (x *MarketplaceMetadata_Eksspec) Reset() {
	*x = MarketplaceMetadata_Eksspec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Eksspec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Eksspec) ProtoMessage() {}

func (x *MarketplaceMetadata_Eksspec) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Eksspec.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Eksspec) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 7}
}

func (x *MarketplaceMetadata_Eksspec) GetNodeGroups() []*MarketplaceMetadata_Nodegroups {
	if x != nil {
		return x.NodeGroups
	}
	return nil
}

type MarketplaceMetadata_Defaultdeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Variables *MarketplaceMetadata_Variables `protobuf:"bytes,1,opt,name=Variables,proto3" json:"Variables,omitempty"`
	EksSpec   *MarketplaceMetadata_Eksspec   `protobuf:"bytes,2,opt,name=EksSpec,proto3" json:"EksSpec,omitempty"`
}

func (x *MarketplaceMetadata_Defaultdeployment) Reset() {
	*x = MarketplaceMetadata_Defaultdeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketplace_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketplaceMetadata_Defaultdeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketplaceMetadata_Defaultdeployment) ProtoMessage() {}

func (x *MarketplaceMetadata_Defaultdeployment) ProtoReflect() protoreflect.Message {
	mi := &file_marketplace_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketplaceMetadata_Defaultdeployment.ProtoReflect.Descriptor instead.
func (*MarketplaceMetadata_Defaultdeployment) Descriptor() ([]byte, []int) {
	return file_marketplace_proto_rawDescGZIP(), []int{0, 8}
}

func (x *MarketplaceMetadata_Defaultdeployment) GetVariables() *MarketplaceMetadata_Variables {
	if x != nil {
		return x.Variables
	}
	return nil
}

func (x *MarketplaceMetadata_Defaultdeployment) GetEksSpec() *MarketplaceMetadata_Eksspec {
	if x != nil {
		return x.EksSpec
	}
	return nil
}

var File_marketplace_proto protoreflect.FileDescriptor

var file_marketplace_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x96, 0x0a, 0x0a, 0x13, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x52,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x54, 0x61, 0x67, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x39, 0x0a, 0x08, 0x49,
	0x61, 0x6d, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x49, 0x61, 0x6d, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x08, 0x49, 0x61,
	0x6d, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65,
	0x12, 0x5a, 0x0a, 0x13, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x44, 0x65, 0x70, 0x65, 0x6e,
	0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x64, 0x65, 0x70, 0x65, 0x6e,
	0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x52, 0x13, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64,
	0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x42,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x12, 0x54, 0x0a, 0x11, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x64,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x11, 0x44, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x57, 0x0a, 0x09,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x45, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x45, 0x66, 0x66, 0x65, 0x63,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x48, 0x0a, 0x08, 0x49, 0x61, 0x6d, 0x72, 0x6f, 0x6c, 0x65,
	0x73, 0x12, 0x3c, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61,
	0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a,
	0x2d, 0x0a, 0x03, 0x45, 0x6b, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x4d, 0x69, 0x6e, 0x69, 0x6d, 0x75,
	0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x4d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x41,
	0x0a, 0x13, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65,
	0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x03, 0x45, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x45, 0x6b, 0x73, 0x52, 0x03, 0x45, 0x6b,
	0x73, 0x1a, 0x43, 0x0a, 0x09, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x36,
	0x0a, 0x17, 0x73, 0x6f, 0x6d, 0x65, 0x5f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d,
	0x5f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x15, 0x73, 0x6f, 0x6d, 0x65, 0x54, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x61,
	0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x1a, 0x8c, 0x01, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x31, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x69, 0x6e, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x4d, 0x69, 0x6e, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x61, 0x78, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x4d, 0x61, 0x78, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x22, 0x0a,
	0x0c, 0x44, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0c, 0x44, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x22, 0x0a, 0x0c, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x4d, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x12, 0x3f, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74,
	0x70, 0x6c, 0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x31, 0x52, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x31, 0x1a, 0x4a, 0x0a, 0x07, 0x45, 0x6b, 0x73, 0x73, 0x70, 0x65, 0x63, 0x12,
	0x3f, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63,
	0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x52, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x1a, 0x89, 0x01, 0x0a, 0x11, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x64, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x09, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x4d, 0x61, 0x72, 0x6b,
	0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x52, 0x09, 0x56, 0x61, 0x72, 0x69, 0x61,
	0x62, 0x6c, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x07, 0x45, 0x6b, 0x73, 0x53, 0x70, 0x65, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x45, 0x6b, 0x73, 0x73,
	0x70, 0x65, 0x63, 0x52, 0x07, 0x45, 0x6b, 0x73, 0x53, 0x70, 0x65, 0x63, 0x42, 0x1e, 0x5a, 0x1c,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_marketplace_proto_rawDescOnce sync.Once
	file_marketplace_proto_rawDescData = file_marketplace_proto_rawDesc
)

func file_marketplace_proto_rawDescGZIP() []byte {
	file_marketplace_proto_rawDescOnce.Do(func() {
		file_marketplace_proto_rawDescData = protoimpl.X.CompressGZIP(file_marketplace_proto_rawDescData)
	})
	return file_marketplace_proto_rawDescData
}

var file_marketplace_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_marketplace_proto_goTypes = []interface{}{
	(*MarketplaceMetadata)(nil),                     // 0: MarketplaceMetadata
	(*MarketplaceMetadata_Statement)(nil),           // 1: MarketplaceMetadata.Statement
	(*MarketplaceMetadata_Iamroles)(nil),            // 2: MarketplaceMetadata.Iamroles
	(*MarketplaceMetadata_Eks)(nil),                 // 3: MarketplaceMetadata.Eks
	(*MarketplaceMetadata_Manageddependencies)(nil), // 4: MarketplaceMetadata.Manageddependencies
	(*MarketplaceMetadata_Variables)(nil),           // 5: MarketplaceMetadata.Variables
	(*MarketplaceMetadata_Nodegroup1)(nil),          // 6: MarketplaceMetadata.Nodegroup1
	(*MarketplaceMetadata_Nodegroups)(nil),          // 7: MarketplaceMetadata.Nodegroups
	(*MarketplaceMetadata_Eksspec)(nil),             // 8: MarketplaceMetadata.Eksspec
	(*MarketplaceMetadata_Defaultdeployment)(nil),   // 9: MarketplaceMetadata.Defaultdeployment
}
var file_marketplace_proto_depIdxs = []int32{
	2, // 0: MarketplaceMetadata.IamRoles:type_name -> MarketplaceMetadata.Iamroles
	4, // 1: MarketplaceMetadata.ManagedDependencies:type_name -> MarketplaceMetadata.Manageddependencies
	9, // 2: MarketplaceMetadata.DefaultDeployment:type_name -> MarketplaceMetadata.Defaultdeployment
	1, // 3: MarketplaceMetadata.Iamroles.Statement:type_name -> MarketplaceMetadata.Statement
	3, // 4: MarketplaceMetadata.Manageddependencies.Eks:type_name -> MarketplaceMetadata.Eks
	6, // 5: MarketplaceMetadata.Nodegroups.NodeGroup1:type_name -> MarketplaceMetadata.Nodegroup1
	7, // 6: MarketplaceMetadata.Eksspec.NodeGroups:type_name -> MarketplaceMetadata.Nodegroups
	5, // 7: MarketplaceMetadata.Defaultdeployment.Variables:type_name -> MarketplaceMetadata.Variables
	8, // 8: MarketplaceMetadata.Defaultdeployment.EksSpec:type_name -> MarketplaceMetadata.Eksspec
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_marketplace_proto_init() }
func file_marketplace_proto_init() {
	if File_marketplace_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_marketplace_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Statement); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Iamroles); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Eks); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Manageddependencies); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Variables); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Nodegroup1); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Nodegroups); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Eksspec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketplace_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketplaceMetadata_Defaultdeployment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_marketplace_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_marketplace_proto_goTypes,
		DependencyIndexes: file_marketplace_proto_depIdxs,
		MessageInfos:      file_marketplace_proto_msgTypes,
	}.Build()
	File_marketplace_proto = out.File
	file_marketplace_proto_rawDesc = nil
	file_marketplace_proto_goTypes = nil
	file_marketplace_proto_depIdxs = nil
}
