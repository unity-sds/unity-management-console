/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Struct } from "../../google/protobuf/struct";

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
  PostInstall: string;
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

export interface MarketplaceMetadata_TypeMap {
  type: string;
  default: string;
}

export interface MarketplaceMetadata_SubMap {
  Options: MarketplaceMetadata_TypeMap | undefined;
}

export interface MarketplaceMetadata_InnerMap {
  Config: { [key: string]: MarketplaceMetadata_SubMap };
  Type: string;
}

export interface MarketplaceMetadata_InnerMap_ConfigEntry {
  key: string;
  value: MarketplaceMetadata_SubMap | undefined;
}

export interface MarketplaceMetadata_Variables {
  Values: { [key: string]: string };
  NestedValues: { [key: string]: MarketplaceMetadata_InnerMap };
  AdvancedValues: { [key: string]: any } | undefined;
}

export interface MarketplaceMetadata_Variables_ValuesEntry {
  key: string;
  value: string;
}

export interface MarketplaceMetadata_Variables_NestedValuesEntry {
  key: string;
  value: MarketplaceMetadata_InnerMap | undefined;
}

export interface MarketplaceMetadata_Defaultdeployment {
  Variables: MarketplaceMetadata_Variables | undefined;
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
    PostInstall: "",
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
    if (message.PostInstall !== "") {
      writer.uint32(138).string(message.PostInstall);
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
        case 17:
          if (tag !== 138) {
            break;
          }

          message.PostInstall = reader.string();
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
      PostInstall: isSet(object.PostInstall) ? String(object.PostInstall) : "",
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
    message.PostInstall !== undefined && (obj.PostInstall = message.PostInstall);
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
    message.PostInstall = object.PostInstall ?? "";
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

function createBaseMarketplaceMetadata_TypeMap(): MarketplaceMetadata_TypeMap {
  return { type: "", default: "" };
}

export const MarketplaceMetadata_TypeMap = {
  encode(message: MarketplaceMetadata_TypeMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.default !== "") {
      writer.uint32(18).string(message.default);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_TypeMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_TypeMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.type = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.default = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_TypeMap {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      default: isSet(object.default) ? String(object.default) : "",
    };
  },

  toJSON(message: MarketplaceMetadata_TypeMap): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.default !== undefined && (obj.default = message.default);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_TypeMap>, I>>(base?: I): MarketplaceMetadata_TypeMap {
    return MarketplaceMetadata_TypeMap.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_TypeMap>, I>>(object: I): MarketplaceMetadata_TypeMap {
    const message = createBaseMarketplaceMetadata_TypeMap();
    message.type = object.type ?? "";
    message.default = object.default ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_SubMap(): MarketplaceMetadata_SubMap {
  return { Options: undefined };
}

export const MarketplaceMetadata_SubMap = {
  encode(message: MarketplaceMetadata_SubMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Options !== undefined) {
      MarketplaceMetadata_TypeMap.encode(message.Options, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_SubMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_SubMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Options = MarketplaceMetadata_TypeMap.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_SubMap {
    return { Options: isSet(object.Options) ? MarketplaceMetadata_TypeMap.fromJSON(object.Options) : undefined };
  },

  toJSON(message: MarketplaceMetadata_SubMap): unknown {
    const obj: any = {};
    message.Options !== undefined &&
      (obj.Options = message.Options ? MarketplaceMetadata_TypeMap.toJSON(message.Options) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_SubMap>, I>>(base?: I): MarketplaceMetadata_SubMap {
    return MarketplaceMetadata_SubMap.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_SubMap>, I>>(object: I): MarketplaceMetadata_SubMap {
    const message = createBaseMarketplaceMetadata_SubMap();
    message.Options = (object.Options !== undefined && object.Options !== null)
      ? MarketplaceMetadata_TypeMap.fromPartial(object.Options)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_InnerMap(): MarketplaceMetadata_InnerMap {
  return { Config: {}, Type: "" };
}

export const MarketplaceMetadata_InnerMap = {
  encode(message: MarketplaceMetadata_InnerMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.Config).forEach(([key, value]) => {
      MarketplaceMetadata_InnerMap_ConfigEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    if (message.Type !== "") {
      writer.uint32(18).string(message.Type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_InnerMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_InnerMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = MarketplaceMetadata_InnerMap_ConfigEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.Config[entry1.key] = entry1.value;
          }
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Type = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_InnerMap {
    return {
      Config: isObject(object.Config)
        ? Object.entries(object.Config).reduce<{ [key: string]: MarketplaceMetadata_SubMap }>((acc, [key, value]) => {
          acc[key] = MarketplaceMetadata_SubMap.fromJSON(value);
          return acc;
        }, {})
        : {},
      Type: isSet(object.Type) ? String(object.Type) : "",
    };
  },

  toJSON(message: MarketplaceMetadata_InnerMap): unknown {
    const obj: any = {};
    obj.Config = {};
    if (message.Config) {
      Object.entries(message.Config).forEach(([k, v]) => {
        obj.Config[k] = MarketplaceMetadata_SubMap.toJSON(v);
      });
    }
    message.Type !== undefined && (obj.Type = message.Type);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_InnerMap>, I>>(base?: I): MarketplaceMetadata_InnerMap {
    return MarketplaceMetadata_InnerMap.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_InnerMap>, I>>(object: I): MarketplaceMetadata_InnerMap {
    const message = createBaseMarketplaceMetadata_InnerMap();
    message.Config = Object.entries(object.Config ?? {}).reduce<{ [key: string]: MarketplaceMetadata_SubMap }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = MarketplaceMetadata_SubMap.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.Type = object.Type ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_InnerMap_ConfigEntry(): MarketplaceMetadata_InnerMap_ConfigEntry {
  return { key: "", value: undefined };
}

export const MarketplaceMetadata_InnerMap_ConfigEntry = {
  encode(message: MarketplaceMetadata_InnerMap_ConfigEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      MarketplaceMetadata_SubMap.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_InnerMap_ConfigEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_InnerMap_ConfigEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = MarketplaceMetadata_SubMap.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_InnerMap_ConfigEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? MarketplaceMetadata_SubMap.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata_InnerMap_ConfigEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value ? MarketplaceMetadata_SubMap.toJSON(message.value) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_InnerMap_ConfigEntry>, I>>(
    base?: I,
  ): MarketplaceMetadata_InnerMap_ConfigEntry {
    return MarketplaceMetadata_InnerMap_ConfigEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_InnerMap_ConfigEntry>, I>>(
    object: I,
  ): MarketplaceMetadata_InnerMap_ConfigEntry {
    const message = createBaseMarketplaceMetadata_InnerMap_ConfigEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? MarketplaceMetadata_SubMap.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Variables(): MarketplaceMetadata_Variables {
  return { Values: {}, NestedValues: {}, AdvancedValues: undefined };
}

export const MarketplaceMetadata_Variables = {
  encode(message: MarketplaceMetadata_Variables, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.Values).forEach(([key, value]) => {
      MarketplaceMetadata_Variables_ValuesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    Object.entries(message.NestedValues).forEach(([key, value]) => {
      MarketplaceMetadata_Variables_NestedValuesEntry.encode({ key: key as any, value }, writer.uint32(18).fork())
        .ldelim();
    });
    if (message.AdvancedValues !== undefined) {
      Struct.encode(Struct.wrap(message.AdvancedValues), writer.uint32(26).fork()).ldelim();
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

          const entry1 = MarketplaceMetadata_Variables_ValuesEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.Values[entry1.key] = entry1.value;
          }
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          const entry2 = MarketplaceMetadata_Variables_NestedValuesEntry.decode(reader, reader.uint32());
          if (entry2.value !== undefined) {
            message.NestedValues[entry2.key] = entry2.value;
          }
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.AdvancedValues = Struct.unwrap(Struct.decode(reader, reader.uint32()));
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
    return {
      Values: isObject(object.Values)
        ? Object.entries(object.Values).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
      NestedValues: isObject(object.NestedValues)
        ? Object.entries(object.NestedValues).reduce<{ [key: string]: MarketplaceMetadata_InnerMap }>(
          (acc, [key, value]) => {
            acc[key] = MarketplaceMetadata_InnerMap.fromJSON(value);
            return acc;
          },
          {},
        )
        : {},
      AdvancedValues: isObject(object.AdvancedValues) ? object.AdvancedValues : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata_Variables): unknown {
    const obj: any = {};
    obj.Values = {};
    if (message.Values) {
      Object.entries(message.Values).forEach(([k, v]) => {
        obj.Values[k] = v;
      });
    }
    obj.NestedValues = {};
    if (message.NestedValues) {
      Object.entries(message.NestedValues).forEach(([k, v]) => {
        obj.NestedValues[k] = MarketplaceMetadata_InnerMap.toJSON(v);
      });
    }
    message.AdvancedValues !== undefined && (obj.AdvancedValues = message.AdvancedValues);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Variables>, I>>(base?: I): MarketplaceMetadata_Variables {
    return MarketplaceMetadata_Variables.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Variables>, I>>(
    object: I,
  ): MarketplaceMetadata_Variables {
    const message = createBaseMarketplaceMetadata_Variables();
    message.Values = Object.entries(object.Values ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = String(value);
      }
      return acc;
    }, {});
    message.NestedValues = Object.entries(object.NestedValues ?? {}).reduce<
      { [key: string]: MarketplaceMetadata_InnerMap }
    >((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = MarketplaceMetadata_InnerMap.fromPartial(value);
      }
      return acc;
    }, {});
    message.AdvancedValues = object.AdvancedValues ?? undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Variables_ValuesEntry(): MarketplaceMetadata_Variables_ValuesEntry {
  return { key: "", value: "" };
}

export const MarketplaceMetadata_Variables_ValuesEntry = {
  encode(message: MarketplaceMetadata_Variables_ValuesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Variables_ValuesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Variables_ValuesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Variables_ValuesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: MarketplaceMetadata_Variables_ValuesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Variables_ValuesEntry>, I>>(
    base?: I,
  ): MarketplaceMetadata_Variables_ValuesEntry {
    return MarketplaceMetadata_Variables_ValuesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Variables_ValuesEntry>, I>>(
    object: I,
  ): MarketplaceMetadata_Variables_ValuesEntry {
    const message = createBaseMarketplaceMetadata_Variables_ValuesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseMarketplaceMetadata_Variables_NestedValuesEntry(): MarketplaceMetadata_Variables_NestedValuesEntry {
  return { key: "", value: undefined };
}

export const MarketplaceMetadata_Variables_NestedValuesEntry = {
  encode(
    message: MarketplaceMetadata_Variables_NestedValuesEntry,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      MarketplaceMetadata_InnerMap.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_Variables_NestedValuesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_Variables_NestedValuesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = MarketplaceMetadata_InnerMap.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarketplaceMetadata_Variables_NestedValuesEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? MarketplaceMetadata_InnerMap.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: MarketplaceMetadata_Variables_NestedValuesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value ? MarketplaceMetadata_InnerMap.toJSON(message.value) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_Variables_NestedValuesEntry>, I>>(
    base?: I,
  ): MarketplaceMetadata_Variables_NestedValuesEntry {
    return MarketplaceMetadata_Variables_NestedValuesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_Variables_NestedValuesEntry>, I>>(
    object: I,
  ): MarketplaceMetadata_Variables_NestedValuesEntry {
    const message = createBaseMarketplaceMetadata_Variables_NestedValuesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? MarketplaceMetadata_InnerMap.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseMarketplaceMetadata_Defaultdeployment(): MarketplaceMetadata_Defaultdeployment {
  return { Variables: undefined };
}

export const MarketplaceMetadata_Defaultdeployment = {
  encode(message: MarketplaceMetadata_Defaultdeployment, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Variables !== undefined) {
      MarketplaceMetadata_Variables.encode(message.Variables, writer.uint32(10).fork()).ldelim();
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
    };
  },

  toJSON(message: MarketplaceMetadata_Defaultdeployment): unknown {
    const obj: any = {};
    message.Variables !== undefined &&
      (obj.Variables = message.Variables ? MarketplaceMetadata_Variables.toJSON(message.Variables) : undefined);
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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
