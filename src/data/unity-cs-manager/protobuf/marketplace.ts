/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "";

export interface MarketplaceMetadata {
  Name: string;
  DisplayName: string;
  Version: string;
  Channel: string;
  Owner: string;
  Description: string;
  Repository: string;
  Tags: string[];
  Category: string;
  IamRoles: MarketplaceMetadata_Iamroles | undefined;
  Package: string;
  ManagedDependencies: MarketplaceMetadata_Manageddependencies[];
  Backend: string;
  Entrypoint: string;
  WorkDirectory: string;
  DefaultDeployment: MarketplaceMetadata_Defaultdeployment | undefined;
}

export interface MarketplaceMetadata_Statement {
  Effect: string;
  Action: string[];
  Resource: string[];
}

export interface MarketplaceMetadata_Iamroles {
  Statement: MarketplaceMetadata_Statement[];
}

export interface MarketplaceMetadata_Eks {
  MinimumVersion: string;
}

export interface MarketplaceMetadata_Manageddependencies {
  Eks: MarketplaceMetadata_Eks | undefined;
}

export interface MarketplaceMetadata_Variables {
  someTerraformVariable: string;
}

export interface MarketplaceMetadata_Nodegroup1 {
  MinNodes: number;
  MaxNodes: number;
  DesiredNodes: number;
  InstanceType: string;
}

export interface MarketplaceMetadata_Nodegroups {
  NodeGroup1: MarketplaceMetadata_Nodegroup1 | undefined;
}

export interface MarketplaceMetadata_Eksspec {
  NodeGroups: MarketplaceMetadata_Nodegroups[];
}

export interface MarketplaceMetadata_Defaultdeployment {
  Variables: MarketplaceMetadata_Variables | undefined;
  EksSpec: MarketplaceMetadata_Eksspec | undefined;
}

function createBaseMarketplaceMetadata(): MarketplaceMetadata {
  return {
    Name: "",
    DisplayName: "",
    Version: "",
    Channel: "",
    Owner: "",
    Description: "",
    Repository: "",
    Tags: [],
    Category: "",
    IamRoles: undefined,
    Package: "",
    ManagedDependencies: [],
    Backend: "",
    Entrypoint: "",
    WorkDirectory: "",
    DefaultDeployment: undefined,
  };
}

export const MarketplaceMetadata = {
  encode(message: MarketplaceMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Name !== "") {
      writer.uint32(10).string(message.Name);
    }
    if (message.DisplayName !== "") {
      writer.uint32(122).string(message.DisplayName);
    }
    if (message.Version !== "") {
      writer.uint32(18).string(message.Version);
    }
    if (message.Channel !== "") {
      writer.uint32(26).string(message.Channel);
    }
    if (message.Owner !== "") {
      writer.uint32(34).string(message.Owner);
    }
    if (message.Description !== "") {
      writer.uint32(42).string(message.Description);
    }
    if (message.Repository !== "") {
      writer.uint32(50).string(message.Repository);
    }
    for (const v of message.Tags) {
      writer.uint32(58).string(v!);
    }
    if (message.Category !== "") {
      writer.uint32(66).string(message.Category);
    }
    if (message.IamRoles !== undefined) {
      MarketplaceMetadata_Iamroles.encode(message.IamRoles, writer.uint32(74).fork()).ldelim();
    }
    if (message.Package !== "") {
      writer.uint32(82).string(message.Package);
    }
    for (const v of message.ManagedDependencies) {
      MarketplaceMetadata_Manageddependencies.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    if (message.Backend !== "") {
      writer.uint32(98).string(message.Backend);
    }
    if (message.Entrypoint !== "") {
      writer.uint32(106).string(message.Entrypoint);
    }
    if (message.WorkDirectory !== "") {
      writer.uint32(130).string(message.WorkDirectory);
    }
    if (message.DefaultDeployment !== undefined) {
      MarketplaceMetadata_Defaultdeployment.encode(message.DefaultDeployment, writer.uint32(114).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Name = reader.string();
          continue;
        case 15:
          if (tag !== 122) {
            break;
          }

          message.DisplayName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Version = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Channel = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.Owner = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.Description = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.Repository = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.Tags.push(reader.string());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.Category = reader.string();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.IamRoles = MarketplaceMetadata_Iamroles.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.Package = reader.string();
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.ManagedDependencies.push(MarketplaceMetadata_Manageddependencies.decode(reader, reader.uint32()));
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.Backend = reader.string();
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.Entrypoint = reader.string();
          continue;
        case 16:
          if (tag !== 130) {
            break;
          }

          message.WorkDirectory = reader.string();
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.DefaultDeployment = MarketplaceMetadata_Defaultdeployment.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata {
    return {
      Name: isSet(object.Name) ? String(object.Name) : "",
      DisplayName: isSet(object.DisplayName) ? String(object.DisplayName) : "",
      Version: isSet(object.Version) ? String(object.Version) : "",
      Channel: isSet(object.Channel) ? String(object.Channel) : "",
      Owner: isSet(object.Owner) ? String(object.Owner) : "",
      Description: isSet(object.Description) ? String(object.Description) : "",
      Repository: isSet(object.Repository) ? String(object.Repository) : "",
      Tags: Array.isArray(object?.Tags) ? object.Tags.map((e: any) => String(e)) : [],
      Category: isSet(object.Category) ? String(object.Category) : "",
      IamRoles: isSet(object.IamRoles) ? MarketplaceMetadata_Iamroles.fromJSON(object.IamRoles) : undefined,
      Package: isSet(object.Package) ? String(object.Package) : "",
      ManagedDependencies: Array.isArray(object?.ManagedDependencies)
        ? object.ManagedDependencies.map((e: any) => MarketplaceMetadata_Manageddependencies.fromJSON(e))
        : [],
      Backend: isSet(object.Backend) ? String(object.Backend) : "",
      Entrypoint: isSet(object.Entrypoint) ? String(object.Entrypoint) : "",
      WorkDirectory: isSet(object.WorkDirectory) ? String(object.WorkDirectory) : "",
      DefaultDeployment: isSet(object.DefaultDeployment)
        ? MarketplaceMetadata_Defaultdeployment.fromJSON(object.DefaultDeployment)
        : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata): unknown {
    const obj: any = {};
    message.Name !== undefined && (obj.Name = message.Name);
    message.DisplayName !== undefined && (obj.DisplayName = message.DisplayName);
    message.Version !== undefined && (obj.Version = message.Version);
    message.Channel !== undefined && (obj.Channel = message.Channel);
    message.Owner !== undefined && (obj.Owner = message.Owner);
    message.Description !== undefined && (obj.Description = message.Description);
    message.Repository !== undefined && (obj.Repository = message.Repository);
    if (message.Tags) {
      obj.Tags = message.Tags.map((e) => e);
    } else {
      obj.Tags = [];
    }
    message.Category !== undefined && (obj.Category = message.Category);
    message.IamRoles !== undefined &&
      (obj.IamRoles = message.IamRoles ? MarketplaceMetadata_Iamroles.toJSON(message.IamRoles) : undefined);
    message.Package !== undefined && (obj.Package = message.Package);
    if (message.ManagedDependencies) {
      obj.ManagedDependencies = message.ManagedDependencies.map((e) =>
        e ? MarketplaceMetadata_Manageddependencies.toJSON(e) : undefined
      );
    } else {
      obj.ManagedDependencies = [];
    }
    message.Backend !== undefined && (obj.Backend = message.Backend);
    message.Entrypoint !== undefined && (obj.Entrypoint = message.Entrypoint);
    message.WorkDirectory !== undefined && (obj.WorkDirectory = message.WorkDirectory);
    message.DefaultDeployment !== undefined && (obj.DefaultDeployment = message.DefaultDeployment
      ? MarketplaceMetadata_Defaultdeployment.toJSON(message.DefaultDeployment)
      : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata>, I>>(base?: I): MarketplaceMetadata {
    return MarketplaceMetadata.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata>, I>>(object: I): MarketplaceMetadata {
    const message = createBaseMarketplaceMetadata();
    message.Name = object.Name ?? "";
    message.DisplayName = object.DisplayName ?? "";
    message.Version = object.Version ?? "";
    message.Channel = object.Channel ?? "";
    message.Owner = object.Owner ?? "";
    message.Description = object.Description ?? "";
    message.Repository = object.Repository ?? "";
    message.Tags = object.Tags?.map((e) => e) || [];
    message.Category = object.Category ?? "";
    message.IamRoles = (object.IamRoles !== undefined && object.IamRoles !== null)
      ? MarketplaceMetadata_Iamroles.fromPartial(object.IamRoles)
      : undefined;
    message.Package = object.Package ?? "";
    message.ManagedDependencies =
      object.ManagedDependencies?.map((e) => MarketplaceMetadata_Manageddependencies.fromPartial(e)) || [];
    message.Backend = object.Backend ?? "";
    message.Entrypoint = object.Entrypoint ?? "";
    message.WorkDirectory = object.WorkDirectory ?? "";
    message.DefaultDeployment = (object.DefaultDeployment !== undefined && object.DefaultDeployment !== null)
      ? MarketplaceMetadata_Defaultdeployment.fromPartial(object.DefaultDeployment)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Statement(): MarketplaceMetadata_Statement {
  return { Effect: "", Action: [], Resource: [] };
}

export const MarketplaceMetadata_Statement = {
  encode(message: MarketplaceMetadata_Statement, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Effect !== "") {
      writer.uint32(10).string(message.Effect);
    }
    for (const v of message.Action) {
      writer.uint32(18).string(v!);
    }
    for (const v of message.Resource) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Statement {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Statement();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Effect = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Action.push(reader.string());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Resource.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Statement {
    return {
      Effect: isSet(object.Effect) ? String(object.Effect) : "",
      Action: Array.isArray(object?.Action) ? object.Action.map((e: any) => String(e)) : [],
      Resource: Array.isArray(object?.Resource) ? object.Resource.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: MarketplaceMetadata_Statement): unknown {
    const obj: any = {};
    message.Effect !== undefined && (obj.Effect = message.Effect);
    if (message.Action) {
      obj.Action = message.Action.map((e) => e);
    } else {
      obj.Action = [];
    }
    if (message.Resource) {
      obj.Resource = message.Resource.map((e) => e);
    } else {
      obj.Resource = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Statement>, I>>(base?: I): MarketplaceMetadata_Statement {
    return MarketplaceMetadata_Statement.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Statement>, I>>(
    object: I,
  ): MarketplaceMetadata_Statement {
    const message = createBaseMarketplaceMetadata_Statement();
    message.Effect = object.Effect ?? "";
    message.Action = object.Action?.map((e) => e) || [];
    message.Resource = object.Resource?.map((e) => e) || [];
    return message;
  },
};

function createBaseMarketplaceMetadata_Iamroles(): MarketplaceMetadata_Iamroles {
  return { Statement: [] };
}

export const MarketplaceMetadata_Iamroles = {
  encode(message: MarketplaceMetadata_Iamroles, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Statement) {
      MarketplaceMetadata_Statement.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Iamroles {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Iamroles();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Statement.push(MarketplaceMetadata_Statement.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Iamroles {
    return {
      Statement: Array.isArray(object?.Statement)
        ? object.Statement.map((e: any) => MarketplaceMetadata_Statement.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MarketplaceMetadata_Iamroles): unknown {
    const obj: any = {};
    if (message.Statement) {
      obj.Statement = message.Statement.map((e) => e ? MarketplaceMetadata_Statement.toJSON(e) : undefined);
    } else {
      obj.Statement = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Iamroles>, I>>(base?: I): MarketplaceMetadata_Iamroles {
    return MarketplaceMetadata_Iamroles.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Iamroles>, I>>(object: I): MarketplaceMetadata_Iamroles {
    const message = createBaseMarketplaceMetadata_Iamroles();
    message.Statement = object.Statement?.map((e) => MarketplaceMetadata_Statement.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMarketplaceMetadata_Eks(): MarketplaceMetadata_Eks {
  return { MinimumVersion: "" };
}

export const MarketplaceMetadata_Eks = {
  encode(message: MarketplaceMetadata_Eks, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.MinimumVersion !== "") {
      writer.uint32(10).string(message.MinimumVersion);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Eks {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Eks();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.MinimumVersion = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Eks {
    return { MinimumVersion: isSet(object.MinimumVersion) ? String(object.MinimumVersion) : "" };
  },

  toJSON(message: MarketplaceMetadata_Eks): unknown {
    const obj: any = {};
    message.MinimumVersion !== undefined && (obj.MinimumVersion = message.MinimumVersion);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Eks>, I>>(base?: I): MarketplaceMetadata_Eks {
    return MarketplaceMetadata_Eks.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Eks>, I>>(object: I): MarketplaceMetadata_Eks {
    const message = createBaseMarketplaceMetadata_Eks();
    message.MinimumVersion = object.MinimumVersion ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_Manageddependencies(): MarketplaceMetadata_Manageddependencies {
  return { Eks: undefined };
}

export const MarketplaceMetadata_Manageddependencies = {
  encode(message: MarketplaceMetadata_Manageddependencies, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Eks !== undefined) {
      MarketplaceMetadata_Eks.encode(message.Eks, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Manageddependencies {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Manageddependencies();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Eks = MarketplaceMetadata_Eks.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Manageddependencies {
    return { Eks: isSet(object.Eks) ? MarketplaceMetadata_Eks.fromJSON(object.Eks) : undefined };
  },

  toJSON(message: MarketplaceMetadata_Manageddependencies): unknown {
    const obj: any = {};
    message.Eks !== undefined && (obj.Eks = message.Eks ? MarketplaceMetadata_Eks.toJSON(message.Eks) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Manageddependencies>, I>>(
    base?: I,
  ): MarketplaceMetadata_Manageddependencies {
    return MarketplaceMetadata_Manageddependencies.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Manageddependencies>, I>>(
    object: I,
  ): MarketplaceMetadata_Manageddependencies {
    const message = createBaseMarketplaceMetadata_Manageddependencies();
    message.Eks = (object.Eks !== undefined && object.Eks !== null)
      ? MarketplaceMetadata_Eks.fromPartial(object.Eks)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Variables(): MarketplaceMetadata_Variables {
  return { someTerraformVariable: "" };
}

export const MarketplaceMetadata_Variables = {
  encode(message: MarketplaceMetadata_Variables, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.someTerraformVariable !== "") {
      writer.uint32(10).string(message.someTerraformVariable);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Variables {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Variables();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.someTerraformVariable = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Variables {
    return { someTerraformVariable: isSet(object.someTerraformVariable) ? String(object.someTerraformVariable) : "" };
  },

  toJSON(message: MarketplaceMetadata_Variables): unknown {
    const obj: any = {};
    message.someTerraformVariable !== undefined && (obj.someTerraformVariable = message.someTerraformVariable);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Variables>, I>>(base?: I): MarketplaceMetadata_Variables {
    return MarketplaceMetadata_Variables.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Variables>, I>>(
    object: I,
  ): MarketplaceMetadata_Variables {
    const message = createBaseMarketplaceMetadata_Variables();
    message.someTerraformVariable = object.someTerraformVariable ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_Nodegroup1(): MarketplaceMetadata_Nodegroup1 {
  return { MinNodes: 0, MaxNodes: 0, DesiredNodes: 0, InstanceType: "" };
}

export const MarketplaceMetadata_Nodegroup1 = {
  encode(message: MarketplaceMetadata_Nodegroup1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.MinNodes !== 0) {
      writer.uint32(8).uint32(message.MinNodes);
    }
    if (message.MaxNodes !== 0) {
      writer.uint32(16).uint32(message.MaxNodes);
    }
    if (message.DesiredNodes !== 0) {
      writer.uint32(24).uint32(message.DesiredNodes);
    }
    if (message.InstanceType !== "") {
      writer.uint32(34).string(message.InstanceType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Nodegroup1 {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Nodegroup1();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.MinNodes = reader.uint32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.MaxNodes = reader.uint32();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.DesiredNodes = reader.uint32();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.InstanceType = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Nodegroup1 {
    return {
      MinNodes: isSet(object.MinNodes) ? Number(object.MinNodes) : 0,
      MaxNodes: isSet(object.MaxNodes) ? Number(object.MaxNodes) : 0,
      DesiredNodes: isSet(object.DesiredNodes) ? Number(object.DesiredNodes) : 0,
      InstanceType: isSet(object.InstanceType) ? String(object.InstanceType) : "",
    };
  },

  toJSON(message: MarketplaceMetadata_Nodegroup1): unknown {
    const obj: any = {};
    message.MinNodes !== undefined && (obj.MinNodes = Math.round(message.MinNodes));
    message.MaxNodes !== undefined && (obj.MaxNodes = Math.round(message.MaxNodes));
    message.DesiredNodes !== undefined && (obj.DesiredNodes = Math.round(message.DesiredNodes));
    message.InstanceType !== undefined && (obj.InstanceType = message.InstanceType);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Nodegroup1>, I>>(base?: I): MarketplaceMetadata_Nodegroup1 {
    return MarketplaceMetadata_Nodegroup1.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Nodegroup1>, I>>(
    object: I,
  ): MarketplaceMetadata_Nodegroup1 {
    const message = createBaseMarketplaceMetadata_Nodegroup1();
    message.MinNodes = object.MinNodes ?? 0;
    message.MaxNodes = object.MaxNodes ?? 0;
    message.DesiredNodes = object.DesiredNodes ?? 0;
    message.InstanceType = object.InstanceType ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_Nodegroups(): MarketplaceMetadata_Nodegroups {
  return { NodeGroup1: undefined };
}

export const MarketplaceMetadata_Nodegroups = {
  encode(message: MarketplaceMetadata_Nodegroups, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.NodeGroup1 !== undefined) {
      MarketplaceMetadata_Nodegroup1.encode(message.NodeGroup1, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Nodegroups {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Nodegroups();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.NodeGroup1 = MarketplaceMetadata_Nodegroup1.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Nodegroups {
    return {
      NodeGroup1: isSet(object.NodeGroup1) ? MarketplaceMetadata_Nodegroup1.fromJSON(object.NodeGroup1) : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata_Nodegroups): unknown {
    const obj: any = {};
    message.NodeGroup1 !== undefined &&
      (obj.NodeGroup1 = message.NodeGroup1 ? MarketplaceMetadata_Nodegroup1.toJSON(message.NodeGroup1) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Nodegroups>, I>>(base?: I): MarketplaceMetadata_Nodegroups {
    return MarketplaceMetadata_Nodegroups.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Nodegroups>, I>>(
    object: I,
  ): MarketplaceMetadata_Nodegroups {
    const message = createBaseMarketplaceMetadata_Nodegroups();
    message.NodeGroup1 = (object.NodeGroup1 !== undefined && object.NodeGroup1 !== null)
      ? MarketplaceMetadata_Nodegroup1.fromPartial(object.NodeGroup1)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Eksspec(): MarketplaceMetadata_Eksspec {
  return { NodeGroups: [] };
}

export const MarketplaceMetadata_Eksspec = {
  encode(message: MarketplaceMetadata_Eksspec, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.NodeGroups) {
      MarketplaceMetadata_Nodegroups.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Eksspec {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Eksspec();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.NodeGroups.push(MarketplaceMetadata_Nodegroups.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Eksspec {
    return {
      NodeGroups: Array.isArray(object?.NodeGroups)
        ? object.NodeGroups.map((e: any) => MarketplaceMetadata_Nodegroups.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MarketplaceMetadata_Eksspec): unknown {
    const obj: any = {};
    if (message.NodeGroups) {
      obj.NodeGroups = message.NodeGroups.map((e) => e ? MarketplaceMetadata_Nodegroups.toJSON(e) : undefined);
    } else {
      obj.NodeGroups = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Eksspec>, I>>(base?: I): MarketplaceMetadata_Eksspec {
    return MarketplaceMetadata_Eksspec.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Eksspec>, I>>(object: I): MarketplaceMetadata_Eksspec {
    const message = createBaseMarketplaceMetadata_Eksspec();
    message.NodeGroups = object.NodeGroups?.map((e) => MarketplaceMetadata_Nodegroups.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMarketplaceMetadata_Defaultdeployment(): MarketplaceMetadata_Defaultdeployment {
  return { Variables: undefined, EksSpec: undefined };
}

export const MarketplaceMetadata_Defaultdeployment = {
  encode(message: MarketplaceMetadata_Defaultdeployment, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Variables !== undefined) {
      MarketplaceMetadata_Variables.encode(message.Variables, writer.uint32(10).fork()).ldelim();
    }
    if (message.EksSpec !== undefined) {
      MarketplaceMetadata_Eksspec.encode(message.EksSpec, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Defaultdeployment {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Defaultdeployment();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Variables = MarketplaceMetadata_Variables.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.EksSpec = MarketplaceMetadata_Eksspec.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Defaultdeployment {
    return {
      Variables: isSet(object.Variables) ? MarketplaceMetadata_Variables.fromJSON(object.Variables) : undefined,
      EksSpec: isSet(object.EksSpec) ? MarketplaceMetadata_Eksspec.fromJSON(object.EksSpec) : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata_Defaultdeployment): unknown {
    const obj: any = {};
    message.Variables !== undefined &&
      (obj.Variables = message.Variables ? MarketplaceMetadata_Variables.toJSON(message.Variables) : undefined);
    message.EksSpec !== undefined &&
      (obj.EksSpec = message.EksSpec ? MarketplaceMetadata_Eksspec.toJSON(message.EksSpec) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Defaultdeployment>, I>>(
    base?: I,
  ): MarketplaceMetadata_Defaultdeployment {
    return MarketplaceMetadata_Defaultdeployment.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Defaultdeployment>, I>>(
    object: I,
  ): MarketplaceMetadata_Defaultdeployment {
    const message = createBaseMarketplaceMetadata_Defaultdeployment();
    message.Variables = (object.Variables !== undefined && object.Variables !== null)
      ? MarketplaceMetadata_Variables.fromPartial(object.Variables)
      : undefined;
    message.EksSpec = (object.EksSpec !== undefined && object.EksSpec !== null)
      ? MarketplaceMetadata_Eksspec.fromPartial(object.EksSpec)
      : undefined;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
