/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "";

export interface Config {
  applicationConfig: Config_ApplicationConfig | undefined;
  networkConfig: Config_NetworkConfig | undefined;
}

export interface Config_ApplicationConfig {
  GithubToken: string;
}

export interface Config_NetworkConfig {
  publicsubnets: string[];
  privatesubnets: string[];
}

export interface Parameters {
  parameterlist: { [key: string]: Parameters_Parameter };
}

export interface Parameters_Parameter {
  value: string;
  type: string;
  tracked: boolean;
  insync: boolean;
}

export interface Parameters_ParameterlistEntry {
  key: string;
  value: Parameters_Parameter | undefined;
}

function createBaseConfig(): Config {
  return { applicationConfig: undefined, networkConfig: undefined };
}

export const Config = {
  encode(message: Config, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.applicationConfig !== undefined) {
      Config_ApplicationConfig.encode(message.applicationConfig, writer.uint32(10).fork()).ldelim();
    }
    if (message.networkConfig !== undefined) {
      Config_NetworkConfig.encode(message.networkConfig, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Config {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.applicationConfig = Config_ApplicationConfig.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.networkConfig = Config_NetworkConfig.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Config {
    return {
      applicationConfig: isSet(object.applicationConfig)
        ? Config_ApplicationConfig.fromJSON(object.applicationConfig)
        : undefined,
      networkConfig: isSet(object.networkConfig) ? Config_NetworkConfig.fromJSON(object.networkConfig) : undefined,
    };
  },

  toJSON(message: Config): unknown {
    const obj: any = {};
    message.applicationConfig !== undefined && (obj.applicationConfig = message.applicationConfig
      ? Config_ApplicationConfig.toJSON(message.applicationConfig)
      : undefined);
    message.networkConfig !== undefined &&
      (obj.networkConfig = message.networkConfig ? Config_NetworkConfig.toJSON(message.networkConfig) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<Config>, I>>(base?: I): Config {
    return Config.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Config>, I>>(object: I): Config {
    const message = createBaseConfig();
    message.applicationConfig = (object.applicationConfig !== undefined && object.applicationConfig !== null)
      ? Config_ApplicationConfig.fromPartial(object.applicationConfig)
      : undefined;
    message.networkConfig = (object.networkConfig !== undefined && object.networkConfig !== null)
      ? Config_NetworkConfig.fromPartial(object.networkConfig)
      : undefined;
    return message;
  },
};

function createBaseConfig_ApplicationConfig(): Config_ApplicationConfig {
  return { GithubToken: "" };
}

export const Config_ApplicationConfig = {
  encode(message: Config_ApplicationConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.GithubToken !== "") {
      writer.uint32(10).string(message.GithubToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Config_ApplicationConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfig_ApplicationConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.GithubToken = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Config_ApplicationConfig {
    return { GithubToken: isSet(object.GithubToken) ? String(object.GithubToken) : "" };
  },

  toJSON(message: Config_ApplicationConfig): unknown {
    const obj: any = {};
    message.GithubToken !== undefined && (obj.GithubToken = message.GithubToken);
    return obj;
  },

  create<I extends Exact<DeepPartial<Config_ApplicationConfig>, I>>(base?: I): Config_ApplicationConfig {
    return Config_ApplicationConfig.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Config_ApplicationConfig>, I>>(object: I): Config_ApplicationConfig {
    const message = createBaseConfig_ApplicationConfig();
    message.GithubToken = object.GithubToken ?? "";
    return message;
  },
};

function createBaseConfig_NetworkConfig(): Config_NetworkConfig {
  return { publicsubnets: [], privatesubnets: [] };
}

export const Config_NetworkConfig = {
  encode(message: Config_NetworkConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.publicsubnets) {
      writer.uint32(10).string(v!);
    }
    for (const v of message.privatesubnets) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Config_NetworkConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfig_NetworkConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.publicsubnets.push(reader.string());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.privatesubnets.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Config_NetworkConfig {
    return {
      publicsubnets: Array.isArray(object?.publicsubnets) ? object.publicsubnets.map((e: any) => String(e)) : [],
      privatesubnets: Array.isArray(object?.privatesubnets) ? object.privatesubnets.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: Config_NetworkConfig): unknown {
    const obj: any = {};
    if (message.publicsubnets) {
      obj.publicsubnets = message.publicsubnets.map((e) => e);
    } else {
      obj.publicsubnets = [];
    }
    if (message.privatesubnets) {
      obj.privatesubnets = message.privatesubnets.map((e) => e);
    } else {
      obj.privatesubnets = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Config_NetworkConfig>, I>>(base?: I): Config_NetworkConfig {
    return Config_NetworkConfig.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Config_NetworkConfig>, I>>(object: I): Config_NetworkConfig {
    const message = createBaseConfig_NetworkConfig();
    message.publicsubnets = object.publicsubnets?.map((e) => e) || [];
    message.privatesubnets = object.privatesubnets?.map((e) => e) || [];
    return message;
  },
};

function createBaseParameters(): Parameters {
  return { parameterlist: {} };
}

export const Parameters = {
  encode(message: Parameters, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.parameterlist).forEach(([key, value]) => {
      Parameters_ParameterlistEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Parameters {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParameters();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = Parameters_ParameterlistEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.parameterlist[entry1.key] = entry1.value;
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

  fromJSON(object: any): Parameters {
    return {
      parameterlist: isObject(object.parameterlist)
        ? Object.entries(object.parameterlist).reduce<{ [key: string]: Parameters_Parameter }>((acc, [key, value]) => {
          acc[key] = Parameters_Parameter.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Parameters): unknown {
    const obj: any = {};
    obj.parameterlist = {};
    if (message.parameterlist) {
      Object.entries(message.parameterlist).forEach(([k, v]) => {
        obj.parameterlist[k] = Parameters_Parameter.toJSON(v);
      });
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Parameters>, I>>(base?: I): Parameters {
    return Parameters.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Parameters>, I>>(object: I): Parameters {
    const message = createBaseParameters();
    message.parameterlist = Object.entries(object.parameterlist ?? {}).reduce<{ [key: string]: Parameters_Parameter }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = Parameters_Parameter.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseParameters_Parameter(): Parameters_Parameter {
  return { value: "", type: "", tracked: false, insync: false };
}

export const Parameters_Parameter = {
  encode(message: Parameters_Parameter, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    if (message.type !== "") {
      writer.uint32(26).string(message.type);
    }
    if (message.tracked === true) {
      writer.uint32(32).bool(message.tracked);
    }
    if (message.insync === true) {
      writer.uint32(40).bool(message.insync);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Parameters_Parameter {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParameters_Parameter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.type = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.tracked = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.insync = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Parameters_Parameter {
    return {
      value: isSet(object.value) ? String(object.value) : "",
      type: isSet(object.type) ? String(object.type) : "",
      tracked: isSet(object.tracked) ? Boolean(object.tracked) : false,
      insync: isSet(object.insync) ? Boolean(object.insync) : false,
    };
  },

  toJSON(message: Parameters_Parameter): unknown {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value);
    message.type !== undefined && (obj.type = message.type);
    message.tracked !== undefined && (obj.tracked = message.tracked);
    message.insync !== undefined && (obj.insync = message.insync);
    return obj;
  },

  create<I extends Exact<DeepPartial<Parameters_Parameter>, I>>(base?: I): Parameters_Parameter {
    return Parameters_Parameter.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Parameters_Parameter>, I>>(object: I): Parameters_Parameter {
    const message = createBaseParameters_Parameter();
    message.value = object.value ?? "";
    message.type = object.type ?? "";
    message.tracked = object.tracked ?? false;
    message.insync = object.insync ?? false;
    return message;
  },
};

function createBaseParameters_ParameterlistEntry(): Parameters_ParameterlistEntry {
  return { key: "", value: undefined };
}

export const Parameters_ParameterlistEntry = {
  encode(message: Parameters_ParameterlistEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Parameters_Parameter.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Parameters_ParameterlistEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParameters_ParameterlistEntry();
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

          message.value = Parameters_Parameter.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Parameters_ParameterlistEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Parameters_Parameter.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Parameters_ParameterlistEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value ? Parameters_Parameter.toJSON(message.value) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<Parameters_ParameterlistEntry>, I>>(base?: I): Parameters_ParameterlistEntry {
    return Parameters_ParameterlistEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Parameters_ParameterlistEntry>, I>>(
    object: I,
  ): Parameters_ParameterlistEntry {
    const message = createBaseParameters_ParameterlistEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? Parameters_Parameter.fromPartial(object.value)
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
