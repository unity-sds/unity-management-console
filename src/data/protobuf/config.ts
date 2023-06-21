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
