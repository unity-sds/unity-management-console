/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "";

export interface Install {
  applications: Install_Applications | undefined;
  extensions: Install_Extensions | undefined;
}

export interface Install_Applications {
  name: string;
  version: string;
  variables: { [key: string]: string };
}

export interface Install_Applications_VariablesEntry {
  key: string;
  value: string;
}

export interface Install_Extensions {
  eks: Install_Extensions_Eks | undefined;
}

export interface Install_Extensions_Nodegroups {
  name: string;
  instancetype: string;
  nodecount: string;
}

export interface Install_Extensions_Eks {
  clustername: string;
  owner: string;
  projectname: string;
  nodegroups: Install_Extensions_Nodegroups[];
}

function createBaseInstall(): Install {
  return { applications: undefined, extensions: undefined };
}

export const Install = {
  encode(message: Install, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.applications !== undefined) {
      Install_Applications.encode(message.applications, writer.uint32(10).fork()).ldelim();
    }
    if (message.extensions !== undefined) {
      Install_Extensions.encode(message.extensions, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.applications = Install_Applications.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.extensions = Install_Extensions.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Install {
    return {
      applications: isSet(object.applications) ? Install_Applications.fromJSON(object.applications) : undefined,
      extensions: isSet(object.extensions) ? Install_Extensions.fromJSON(object.extensions) : undefined,
    };
  },

  toJSON(message: Install): unknown {
    const obj: any = {};
    message.applications !== undefined &&
      (obj.applications = message.applications ? Install_Applications.toJSON(message.applications) : undefined);
    message.extensions !== undefined &&
      (obj.extensions = message.extensions ? Install_Extensions.toJSON(message.extensions) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install>, I>>(base?: I): Install {
    return Install.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install>, I>>(object: I): Install {
    const message = createBaseInstall();
    message.applications = (object.applications !== undefined && object.applications !== null)
      ? Install_Applications.fromPartial(object.applications)
      : undefined;
    message.extensions = (object.extensions !== undefined && object.extensions !== null)
      ? Install_Extensions.fromPartial(object.extensions)
      : undefined;
    return message;
  },
};

function createBaseInstall_Applications(): Install_Applications {
  return { name: "", version: "", variables: {} };
}

export const Install_Applications = {
  encode(message: Install_Applications, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }
    Object.entries(message.variables).forEach(([key, value]) => {
      Install_Applications_VariablesEntry.encode({ key: key as any, value }, writer.uint32(26).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Applications {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Applications();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.version = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          const entry3 = Install_Applications_VariablesEntry.decode(reader, reader.uint32());
          if (entry3.value !== undefined) {
            message.variables[entry3.key] = entry3.value;
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

  fromJSON(object: any): Install_Applications {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      version: isSet(object.version) ? String(object.version) : "",
      variables: isObject(object.variables)
        ? Object.entries(object.variables).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Install_Applications): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.version !== undefined && (obj.version = message.version);
    obj.variables = {};
    if (message.variables) {
      Object.entries(message.variables).forEach(([k, v]) => {
        obj.variables[k] = v;
      });
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Applications>, I>>(base?: I): Install_Applications {
    return Install_Applications.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Applications>, I>>(object: I): Install_Applications {
    const message = createBaseInstall_Applications();
    message.name = object.name ?? "";
    message.version = object.version ?? "";
    message.variables = Object.entries(object.variables ?? {}).reduce<{ [key: string]: string }>(
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

function createBaseInstall_Applications_VariablesEntry(): Install_Applications_VariablesEntry {
  return { key: "", value: "" };
}

export const Install_Applications_VariablesEntry = {
  encode(message: Install_Applications_VariablesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Applications_VariablesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Applications_VariablesEntry();
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

  fromJSON(object: any): Install_Applications_VariablesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: Install_Applications_VariablesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Applications_VariablesEntry>, I>>(
    base?: I,
  ): Install_Applications_VariablesEntry {
    return Install_Applications_VariablesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Applications_VariablesEntry>, I>>(
    object: I,
  ): Install_Applications_VariablesEntry {
    const message = createBaseInstall_Applications_VariablesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseInstall_Extensions(): Install_Extensions {
  return { eks: undefined };
}

export const Install_Extensions = {
  encode(message: Install_Extensions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.eks !== undefined) {
      Install_Extensions_Eks.encode(message.eks, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Extensions {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Extensions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.eks = Install_Extensions_Eks.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Install_Extensions {
    return { eks: isSet(object.eks) ? Install_Extensions_Eks.fromJSON(object.eks) : undefined };
  },

  toJSON(message: Install_Extensions): unknown {
    const obj: any = {};
    message.eks !== undefined && (obj.eks = message.eks ? Install_Extensions_Eks.toJSON(message.eks) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Extensions>, I>>(base?: I): Install_Extensions {
    return Install_Extensions.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Extensions>, I>>(object: I): Install_Extensions {
    const message = createBaseInstall_Extensions();
    message.eks = (object.eks !== undefined && object.eks !== null)
      ? Install_Extensions_Eks.fromPartial(object.eks)
      : undefined;
    return message;
  },
};

function createBaseInstall_Extensions_Nodegroups(): Install_Extensions_Nodegroups {
  return { name: "", instancetype: "", nodecount: "" };
}

export const Install_Extensions_Nodegroups = {
  encode(message: Install_Extensions_Nodegroups, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.instancetype !== "") {
      writer.uint32(18).string(message.instancetype);
    }
    if (message.nodecount !== "") {
      writer.uint32(26).string(message.nodecount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Extensions_Nodegroups {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Extensions_Nodegroups();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.instancetype = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.nodecount = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Install_Extensions_Nodegroups {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      instancetype: isSet(object.instancetype) ? String(object.instancetype) : "",
      nodecount: isSet(object.nodecount) ? String(object.nodecount) : "",
    };
  },

  toJSON(message: Install_Extensions_Nodegroups): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.instancetype !== undefined && (obj.instancetype = message.instancetype);
    message.nodecount !== undefined && (obj.nodecount = message.nodecount);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Extensions_Nodegroups>, I>>(base?: I): Install_Extensions_Nodegroups {
    return Install_Extensions_Nodegroups.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Extensions_Nodegroups>, I>>(
    object: I,
  ): Install_Extensions_Nodegroups {
    const message = createBaseInstall_Extensions_Nodegroups();
    message.name = object.name ?? "";
    message.instancetype = object.instancetype ?? "";
    message.nodecount = object.nodecount ?? "";
    return message;
  },
};

function createBaseInstall_Extensions_Eks(): Install_Extensions_Eks {
  return { clustername: "", owner: "", projectname: "", nodegroups: [] };
}

export const Install_Extensions_Eks = {
  encode(message: Install_Extensions_Eks, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clustername !== "") {
      writer.uint32(10).string(message.clustername);
    }
    if (message.owner !== "") {
      writer.uint32(18).string(message.owner);
    }
    if (message.projectname !== "") {
      writer.uint32(26).string(message.projectname);
    }
    for (const v of message.nodegroups) {
      Install_Extensions_Nodegroups.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Extensions_Eks {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Extensions_Eks();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.clustername = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.owner = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.projectname = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.nodegroups.push(Install_Extensions_Nodegroups.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Install_Extensions_Eks {
    return {
      clustername: isSet(object.clustername) ? String(object.clustername) : "",
      owner: isSet(object.owner) ? String(object.owner) : "",
      projectname: isSet(object.projectname) ? String(object.projectname) : "",
      nodegroups: Array.isArray(object?.nodegroups)
        ? object.nodegroups.map((e: any) => Install_Extensions_Nodegroups.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Install_Extensions_Eks): unknown {
    const obj: any = {};
    message.clustername !== undefined && (obj.clustername = message.clustername);
    message.owner !== undefined && (obj.owner = message.owner);
    message.projectname !== undefined && (obj.projectname = message.projectname);
    if (message.nodegroups) {
      obj.nodegroups = message.nodegroups.map((e) => e ? Install_Extensions_Nodegroups.toJSON(e) : undefined);
    } else {
      obj.nodegroups = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Extensions_Eks>, I>>(base?: I): Install_Extensions_Eks {
    return Install_Extensions_Eks.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Extensions_Eks>, I>>(object: I): Install_Extensions_Eks {
    const message = createBaseInstall_Extensions_Eks();
    message.clustername = object.clustername ?? "";
    message.owner = object.owner ?? "";
    message.projectname = object.projectname ?? "";
    message.nodegroups = object.nodegroups?.map((e) => Install_Extensions_Nodegroups.fromPartial(e)) || [];
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
