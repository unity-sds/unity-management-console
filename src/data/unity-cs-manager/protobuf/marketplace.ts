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
  ManagedDependencies: { [key: string]: any } | undefined;
  Backend: string;
  Entrypoint: string;
  WorkDirectory: string;
  PostInstall: string;
  PreInstall: string;
  DefaultDeployment: MarketplaceMetadata_Defaultdeployment | undefined;
  Dependencies: { [key: string]: string };
}

export interface MarketplaceMetadata_Statement {
  Effect: string;
  Action: string[];
  Resource: string[];
}

export interface MarketplaceMetadata_Iamroles {
  Statement: MarketplaceMetadata_Statement[];
}

export interface MarketplaceMetadata_Variables {
  Values: { [key: string]: string };
  AdvancedValues: { [key: string]: any } | undefined;
}

export interface MarketplaceMetadata_Variables_ValuesEntry {
  key: string;
  value: string;
}

export interface MarketplaceMetadata_Defaultdeployment {
  Variables: MarketplaceMetadata_Variables | undefined;
}

export interface MarketplaceMetadata_DependenciesEntry {
  key: string;
  value: string;
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
    ManagedDependencies: undefined,
    Backend: "",
    Entrypoint: "",
    WorkDirectory: "",
    PostInstall: "",
    PreInstall: "",
    DefaultDeployment: undefined,
    Dependencies: {},
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
    if (message.ManagedDependencies !== undefined) {
      Struct.encode(Struct.wrap(message.ManagedDependencies), writer.uint32(90).fork()).ldelim();
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
    if (message.PreInstall !== "") {
      writer.uint32(146).string(message.PreInstall);
    }
    if (message.DefaultDeployment !== undefined) {
      MarketplaceMetadata_Defaultdeployment.encode(message.DefaultDeployment, writer.uint32(114).fork()).ldelim();
    }
    Object.entries(message.Dependencies).forEach(([key, value]) => {
      MarketplaceMetadata_DependenciesEntry.encode({ key: key as any, value }, writer.uint32(154).fork()).ldelim();
    });
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

          message.ManagedDependencies = Struct.unwrap(Struct.decode(reader, reader.uint32()));
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
        case 18:
          if (tag !== 146) {
            break;
          }

          message.PreInstall = reader.string();
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.DefaultDeployment = MarketplaceMetadata_Defaultdeployment.decode(reader, reader.uint32());
          continue;
        case 19:
          if (tag !== 154) {
            break;
          }

          const entry19 = MarketplaceMetadata_DependenciesEntry.decode(reader, reader.uint32());
          if (entry19.value !== undefined) {
            message.Dependencies[entry19.key] = entry19.value;
          }
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
      ManagedDependencies: isObject(object.ManagedDependencies) ? object.ManagedDependencies : undefined,
      Backend: isSet(object.Backend) ? String(object.Backend) : "",
      Entrypoint: isSet(object.Entrypoint) ? String(object.Entrypoint) : "",
      WorkDirectory: isSet(object.WorkDirectory) ? String(object.WorkDirectory) : "",
      PostInstall: isSet(object.PostInstall) ? String(object.PostInstall) : "",
      PreInstall: isSet(object.PreInstall) ? String(object.PreInstall) : "",
      DefaultDeployment: isSet(object.DefaultDeployment)
        ? MarketplaceMetadata_Defaultdeployment.fromJSON(object.DefaultDeployment)
        : undefined,
      Dependencies: isObject(object.Dependencies)
        ? Object.entries(object.Dependencies).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
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
    message.ManagedDependencies !== undefined && (obj.ManagedDependencies = message.ManagedDependencies);
    message.Backend !== undefined && (obj.Backend = message.Backend);
    message.Entrypoint !== undefined && (obj.Entrypoint = message.Entrypoint);
    message.WorkDirectory !== undefined && (obj.WorkDirectory = message.WorkDirectory);
    message.PostInstall !== undefined && (obj.PostInstall = message.PostInstall);
    message.PreInstall !== undefined && (obj.PreInstall = message.PreInstall);
    message.DefaultDeployment !== undefined && (obj.DefaultDeployment = message.DefaultDeployment
      ? MarketplaceMetadata_Defaultdeployment.toJSON(message.DefaultDeployment)
      : undefined);
    obj.Dependencies = {};
    if (message.Dependencies) {
      Object.entries(message.Dependencies).forEach(([k, v]) => {
        obj.Dependencies[k] = v;
      });
    }
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
    message.ManagedDependencies = object.ManagedDependencies ?? undefined;
    message.Backend = object.Backend ?? "";
    message.Entrypoint = object.Entrypoint ?? "";
    message.WorkDirectory = object.WorkDirectory ?? "";
    message.PostInstall = object.PostInstall ?? "";
    message.PreInstall = object.PreInstall ?? "";
    message.DefaultDeployment = (object.DefaultDeployment !== undefined && object.DefaultDeployment !== null)
      ? MarketplaceMetadata_Defaultdeployment.fromPartial(object.DefaultDeployment)
      : undefined;
    message.Dependencies = Object.entries(object.Dependencies ?? {}).reduce<{ [key: string]: string }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = String(value);
        }
        return acc;
      },
      {},
    );
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

function createBaseMarketplaceMetadata_Variables(): MarketplaceMetadata_Variables {
  return { Values: {}, AdvancedValues: undefined };
}

export const MarketplaceMetadata_Variables = {
  encode(message: MarketplaceMetadata_Variables, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.Values).forEach(([key, value]) => {
      MarketplaceMetadata_Variables_ValuesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
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

function createBaseMarketplaceMetadata_DependenciesEntry(): MarketplaceMetadata_DependenciesEntry {
  return { key: "", value: "" };
}

export const MarketplaceMetadata_DependenciesEntry = {
  encode(message: MarketplaceMetadata_DependenciesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketplaceMetadata_DependenciesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketplaceMetadata_DependenciesEntry();
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

  fromJSON(object: any): MarketplaceMetadata_DependenciesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: MarketplaceMetadata_DependenciesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create<I extends Exact<DeepPartial<MarketplaceMetadata_DependenciesEntry>, I>>(
    base?: I,
  ): MarketplaceMetadata_DependenciesEntry {
    return MarketplaceMetadata_DependenciesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<MarketplaceMetadata_DependenciesEntry>, I>>(
    object: I,
  ): MarketplaceMetadata_DependenciesEntry {
    const message = createBaseMarketplaceMetadata_DependenciesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
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
