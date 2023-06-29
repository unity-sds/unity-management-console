/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "";

export interface Install {
  applications: Install_Applications | undefined;
  extensions: Install_Extensions | undefined;
  DeploymentName: string;
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
  apigateway: Install_Extensions_Apigateway | undefined;
}

export interface Install_Extensions_Apigateway {
  name: string[];
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

export interface ActionMeta {
  metadataVersion: string;
  exectarget: string;
  deploymentName: string;
  services: ActionMeta_Services[];
  extensions: ActionMeta_Extensions | undefined;
  deploymentType: string;
}

export interface ActionMeta_Services {
  name: string;
  source: string;
  version: string;
  branch: string;
}

export interface ActionMeta_Extensions {
  eks: ActionMeta_Extensions_Eks | undefined;
  apis: ActionMeta_Extensions_Apigateway | undefined;
}

export interface ActionMeta_Extensions_Apis {
  name: string;
}

export interface ActionMeta_Extensions_Apigateway {
  apis: ActionMeta_Extensions_Apis[];
}

export interface ActionMeta_Extensions_Nodegroups {
  name: string;
  instancetype: string;
  nodecount: string;
}

export interface ActionMeta_Extensions_Eks {
  clustername: string;
  owner: string;
  projectname: string;
  nodegroups: ActionMeta_Extensions_Nodegroups[];
}

function createBaseInstall(): Install {
  return { applications: undefined, extensions: undefined, DeploymentName: "" };
}

export const Install = {
  encode(message: Install, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.applications !== undefined) {
      Install_Applications.encode(message.applications, writer.uint32(10).fork()).ldelim();
    }
    if (message.extensions !== undefined) {
      Install_Extensions.encode(message.extensions, writer.uint32(18).fork()).ldelim();
    }
    if (message.DeploymentName !== "") {
      writer.uint32(26).string(message.DeploymentName);
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
        case 3:
          if (tag !== 26) {
            break;
          }

          message.DeploymentName = reader.string();
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
      DeploymentName: isSet(object.DeploymentName) ? String(object.DeploymentName) : "",
    };
  },

  toJSON(message: Install): unknown {
    const obj: any = {};
    message.applications !== undefined &&
      (obj.applications = message.applications ? Install_Applications.toJSON(message.applications) : undefined);
    message.extensions !== undefined &&
      (obj.extensions = message.extensions ? Install_Extensions.toJSON(message.extensions) : undefined);
    message.DeploymentName !== undefined && (obj.DeploymentName = message.DeploymentName);
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
    message.DeploymentName = object.DeploymentName ?? "";
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
  return { eks: undefined, apigateway: undefined };
}

export const Install_Extensions = {
  encode(message: Install_Extensions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.eks !== undefined) {
      Install_Extensions_Eks.encode(message.eks, writer.uint32(10).fork()).ldelim();
    }
    if (message.apigateway !== undefined) {
      Install_Extensions_Apigateway.encode(message.apigateway, writer.uint32(18).fork()).ldelim();
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
        case 2:
          if (tag !== 18) {
            break;
          }

          message.apigateway = Install_Extensions_Apigateway.decode(reader, reader.uint32());
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
    return {
      eks: isSet(object.eks) ? Install_Extensions_Eks.fromJSON(object.eks) : undefined,
      apigateway: isSet(object.apigateway) ? Install_Extensions_Apigateway.fromJSON(object.apigateway) : undefined,
    };
  },

  toJSON(message: Install_Extensions): unknown {
    const obj: any = {};
    message.eks !== undefined && (obj.eks = message.eks ? Install_Extensions_Eks.toJSON(message.eks) : undefined);
    message.apigateway !== undefined &&
      (obj.apigateway = message.apigateway ? Install_Extensions_Apigateway.toJSON(message.apigateway) : undefined);
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
    message.apigateway = (object.apigateway !== undefined && object.apigateway !== null)
      ? Install_Extensions_Apigateway.fromPartial(object.apigateway)
      : undefined;
    return message;
  },
};

function createBaseInstall_Extensions_Apigateway(): Install_Extensions_Apigateway {
  return { name: [] };
}

export const Install_Extensions_Apigateway = {
  encode(message: Install_Extensions_Apigateway, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.name) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Extensions_Apigateway {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Extensions_Apigateway();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Install_Extensions_Apigateway {
    return { name: Array.isArray(object?.name) ? object.name.map((e: any) => String(e)) : [] };
  },

  toJSON(message: Install_Extensions_Apigateway): unknown {
    const obj: any = {};
    if (message.name) {
      obj.name = message.name.map((e) => e);
    } else {
      obj.name = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Extensions_Apigateway>, I>>(base?: I): Install_Extensions_Apigateway {
    return Install_Extensions_Apigateway.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Extensions_Apigateway>, I>>(
    object: I,
  ): Install_Extensions_Apigateway {
    const message = createBaseInstall_Extensions_Apigateway();
    message.name = object.name?.map((e) => e) || [];
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

function createBaseActionMeta(): ActionMeta {
  return {
    metadataVersion: "",
    exectarget: "",
    deploymentName: "",
    services: [],
    extensions: undefined,
    deploymentType: "",
  };
}

export const ActionMeta = {
  encode(message: ActionMeta, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.metadataVersion !== "") {
      writer.uint32(10).string(message.metadataVersion);
    }
    if (message.exectarget !== "") {
      writer.uint32(18).string(message.exectarget);
    }
    if (message.deploymentName !== "") {
      writer.uint32(26).string(message.deploymentName);
    }
    for (const v of message.services) {
      ActionMeta_Services.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.extensions !== undefined) {
      ActionMeta_Extensions.encode(message.extensions, writer.uint32(42).fork()).ldelim();
    }
    if (message.deploymentType !== "") {
      writer.uint32(50).string(message.deploymentType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.metadataVersion = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.exectarget = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.deploymentName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.services.push(ActionMeta_Services.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.extensions = ActionMeta_Extensions.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.deploymentType = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta {
    return {
      metadataVersion: isSet(object.metadataVersion) ? String(object.metadataVersion) : "",
      exectarget: isSet(object.exectarget) ? String(object.exectarget) : "",
      deploymentName: isSet(object.deploymentName) ? String(object.deploymentName) : "",
      services: Array.isArray(object?.services) ? object.services.map((e: any) => ActionMeta_Services.fromJSON(e)) : [],
      extensions: isSet(object.extensions) ? ActionMeta_Extensions.fromJSON(object.extensions) : undefined,
      deploymentType: isSet(object.deploymentType) ? String(object.deploymentType) : "",
    };
  },

  toJSON(message: ActionMeta): unknown {
    const obj: any = {};
    message.metadataVersion !== undefined && (obj.metadataVersion = message.metadataVersion);
    message.exectarget !== undefined && (obj.exectarget = message.exectarget);
    message.deploymentName !== undefined && (obj.deploymentName = message.deploymentName);
    if (message.services) {
      obj.services = message.services.map((e) => e ? ActionMeta_Services.toJSON(e) : undefined);
    } else {
      obj.services = [];
    }
    message.extensions !== undefined &&
      (obj.extensions = message.extensions ? ActionMeta_Extensions.toJSON(message.extensions) : undefined);
    message.deploymentType !== undefined && (obj.deploymentType = message.deploymentType);
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta>, I>>(base?: I): ActionMeta {
    return ActionMeta.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta>, I>>(object: I): ActionMeta {
    const message = createBaseActionMeta();
    message.metadataVersion = object.metadataVersion ?? "";
    message.exectarget = object.exectarget ?? "";
    message.deploymentName = object.deploymentName ?? "";
    message.services = object.services?.map((e) => ActionMeta_Services.fromPartial(e)) || [];
    message.extensions = (object.extensions !== undefined && object.extensions !== null)
      ? ActionMeta_Extensions.fromPartial(object.extensions)
      : undefined;
    message.deploymentType = object.deploymentType ?? "";
    return message;
  },
};

function createBaseActionMeta_Services(): ActionMeta_Services {
  return { name: "", source: "", version: "", branch: "" };
}

export const ActionMeta_Services = {
  encode(message: ActionMeta_Services, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.source !== "") {
      writer.uint32(18).string(message.source);
    }
    if (message.version !== "") {
      writer.uint32(26).string(message.version);
    }
    if (message.branch !== "") {
      writer.uint32(34).string(message.branch);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Services {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Services();
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

          message.source = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.version = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.branch = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta_Services {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      source: isSet(object.source) ? String(object.source) : "",
      version: isSet(object.version) ? String(object.version) : "",
      branch: isSet(object.branch) ? String(object.branch) : "",
    };
  },

  toJSON(message: ActionMeta_Services): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.source !== undefined && (obj.source = message.source);
    message.version !== undefined && (obj.version = message.version);
    message.branch !== undefined && (obj.branch = message.branch);
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Services>, I>>(base?: I): ActionMeta_Services {
    return ActionMeta_Services.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Services>, I>>(object: I): ActionMeta_Services {
    const message = createBaseActionMeta_Services();
    message.name = object.name ?? "";
    message.source = object.source ?? "";
    message.version = object.version ?? "";
    message.branch = object.branch ?? "";
    return message;
  },
};

function createBaseActionMeta_Extensions(): ActionMeta_Extensions {
  return { eks: undefined, apis: undefined };
}

export const ActionMeta_Extensions = {
  encode(message: ActionMeta_Extensions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.eks !== undefined) {
      ActionMeta_Extensions_Eks.encode(message.eks, writer.uint32(10).fork()).ldelim();
    }
    if (message.apis !== undefined) {
      ActionMeta_Extensions_Apigateway.encode(message.apis, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Extensions {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Extensions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.eks = ActionMeta_Extensions_Eks.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.apis = ActionMeta_Extensions_Apigateway.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta_Extensions {
    return {
      eks: isSet(object.eks) ? ActionMeta_Extensions_Eks.fromJSON(object.eks) : undefined,
      apis: isSet(object.apis) ? ActionMeta_Extensions_Apigateway.fromJSON(object.apis) : undefined,
    };
  },

  toJSON(message: ActionMeta_Extensions): unknown {
    const obj: any = {};
    message.eks !== undefined && (obj.eks = message.eks ? ActionMeta_Extensions_Eks.toJSON(message.eks) : undefined);
    message.apis !== undefined &&
      (obj.apis = message.apis ? ActionMeta_Extensions_Apigateway.toJSON(message.apis) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Extensions>, I>>(base?: I): ActionMeta_Extensions {
    return ActionMeta_Extensions.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Extensions>, I>>(object: I): ActionMeta_Extensions {
    const message = createBaseActionMeta_Extensions();
    message.eks = (object.eks !== undefined && object.eks !== null)
      ? ActionMeta_Extensions_Eks.fromPartial(object.eks)
      : undefined;
    message.apis = (object.apis !== undefined && object.apis !== null)
      ? ActionMeta_Extensions_Apigateway.fromPartial(object.apis)
      : undefined;
    return message;
  },
};

function createBaseActionMeta_Extensions_Apis(): ActionMeta_Extensions_Apis {
  return { name: "" };
}

export const ActionMeta_Extensions_Apis = {
  encode(message: ActionMeta_Extensions_Apis, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Extensions_Apis {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Extensions_Apis();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta_Extensions_Apis {
    return { name: isSet(object.name) ? String(object.name) : "" };
  },

  toJSON(message: ActionMeta_Extensions_Apis): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Extensions_Apis>, I>>(base?: I): ActionMeta_Extensions_Apis {
    return ActionMeta_Extensions_Apis.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Extensions_Apis>, I>>(object: I): ActionMeta_Extensions_Apis {
    const message = createBaseActionMeta_Extensions_Apis();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseActionMeta_Extensions_Apigateway(): ActionMeta_Extensions_Apigateway {
  return { apis: [] };
}

export const ActionMeta_Extensions_Apigateway = {
  encode(message: ActionMeta_Extensions_Apigateway, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.apis) {
      ActionMeta_Extensions_Apis.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Extensions_Apigateway {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Extensions_Apigateway();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.apis.push(ActionMeta_Extensions_Apis.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta_Extensions_Apigateway {
    return {
      apis: Array.isArray(object?.apis) ? object.apis.map((e: any) => ActionMeta_Extensions_Apis.fromJSON(e)) : [],
    };
  },

  toJSON(message: ActionMeta_Extensions_Apigateway): unknown {
    const obj: any = {};
    if (message.apis) {
      obj.apis = message.apis.map((e) => e ? ActionMeta_Extensions_Apis.toJSON(e) : undefined);
    } else {
      obj.apis = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Extensions_Apigateway>, I>>(
    base?: I,
  ): ActionMeta_Extensions_Apigateway {
    return ActionMeta_Extensions_Apigateway.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Extensions_Apigateway>, I>>(
    object: I,
  ): ActionMeta_Extensions_Apigateway {
    const message = createBaseActionMeta_Extensions_Apigateway();
    message.apis = object.apis?.map((e) => ActionMeta_Extensions_Apis.fromPartial(e)) || [];
    return message;
  },
};

function createBaseActionMeta_Extensions_Nodegroups(): ActionMeta_Extensions_Nodegroups {
  return { name: "", instancetype: "", nodecount: "" };
}

export const ActionMeta_Extensions_Nodegroups = {
  encode(message: ActionMeta_Extensions_Nodegroups, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Extensions_Nodegroups {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Extensions_Nodegroups();
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

  fromJSON(object: any): ActionMeta_Extensions_Nodegroups {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      instancetype: isSet(object.instancetype) ? String(object.instancetype) : "",
      nodecount: isSet(object.nodecount) ? String(object.nodecount) : "",
    };
  },

  toJSON(message: ActionMeta_Extensions_Nodegroups): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.instancetype !== undefined && (obj.instancetype = message.instancetype);
    message.nodecount !== undefined && (obj.nodecount = message.nodecount);
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Extensions_Nodegroups>, I>>(
    base?: I,
  ): ActionMeta_Extensions_Nodegroups {
    return ActionMeta_Extensions_Nodegroups.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Extensions_Nodegroups>, I>>(
    object: I,
  ): ActionMeta_Extensions_Nodegroups {
    const message = createBaseActionMeta_Extensions_Nodegroups();
    message.name = object.name ?? "";
    message.instancetype = object.instancetype ?? "";
    message.nodecount = object.nodecount ?? "";
    return message;
  },
};

function createBaseActionMeta_Extensions_Eks(): ActionMeta_Extensions_Eks {
  return { clustername: "", owner: "", projectname: "", nodegroups: [] };
}

export const ActionMeta_Extensions_Eks = {
  encode(message: ActionMeta_Extensions_Eks, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
      ActionMeta_Extensions_Nodegroups.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionMeta_Extensions_Eks {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionMeta_Extensions_Eks();
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

          message.nodegroups.push(ActionMeta_Extensions_Nodegroups.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionMeta_Extensions_Eks {
    return {
      clustername: isSet(object.clustername) ? String(object.clustername) : "",
      owner: isSet(object.owner) ? String(object.owner) : "",
      projectname: isSet(object.projectname) ? String(object.projectname) : "",
      nodegroups: Array.isArray(object?.nodegroups)
        ? object.nodegroups.map((e: any) => ActionMeta_Extensions_Nodegroups.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ActionMeta_Extensions_Eks): unknown {
    const obj: any = {};
    message.clustername !== undefined && (obj.clustername = message.clustername);
    message.owner !== undefined && (obj.owner = message.owner);
    message.projectname !== undefined && (obj.projectname = message.projectname);
    if (message.nodegroups) {
      obj.nodegroups = message.nodegroups.map((e) => e ? ActionMeta_Extensions_Nodegroups.toJSON(e) : undefined);
    } else {
      obj.nodegroups = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ActionMeta_Extensions_Eks>, I>>(base?: I): ActionMeta_Extensions_Eks {
    return ActionMeta_Extensions_Eks.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ActionMeta_Extensions_Eks>, I>>(object: I): ActionMeta_Extensions_Eks {
    const message = createBaseActionMeta_Extensions_Eks();
    message.clustername = object.clustername ?? "";
    message.owner = object.owner ?? "";
    message.projectname = object.projectname ?? "";
    message.nodegroups = object.nodegroups?.map((e) => ActionMeta_Extensions_Nodegroups.fromPartial(e)) || [];
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
