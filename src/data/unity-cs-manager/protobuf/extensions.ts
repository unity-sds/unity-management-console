/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Struct } from "../../google/protobuf/struct";

export const protobufPackage = "";

export interface UnityWebsocketMessage {
  install?: Install | undefined;
  simplemessage?: SimpleMessage | undefined;
  connectionsetup?: ConnectionSetup | undefined;
  config?: Config | undefined;
  parameters?: Parameters | undefined;
  logs?: LogLine | undefined;
  deployments?: Deployments | undefined;
  uninstall?: Uninstall | undefined;
}

export interface Application {
  packageName: string;
  version: string;
  source: string;
  status: string;
  applicationName: string;
  displayName: string;
}

export interface Deployment {
  name: string;
  creator: string;
  creationdate: string;
  application: Application[];
}

export interface Deployments {
  deployment: Deployment[];
}

export interface ConnectionSetup {
  type: string;
  userID: string;
}

export interface Install {
  applications: Install_Applications | undefined;
  DeploymentName: string;
}

export interface Install_Variables {
  Values: { [key: string]: string };
  AdvancedValues: { [key: string]: any } | undefined;
}

export interface Install_Variables_ValuesEntry {
  key: string;
  value: string;
}

export interface Install_Applications {
  name: string;
  version: string;
  variables: Install_Variables | undefined;
  postinstall: string;
  preinstall: string;
  displayname: string;
  dependencies: { [key: string]: string };
}

export interface Install_Applications_DependenciesEntry {
  key: string;
  value: string;
}

export interface Uninstall {
  DeploymentName: string;
  Application: string;
  DisplayName: string;
  All: boolean;
}

export interface SimpleMessage {
  operation: string;
  payload: string;
}

export interface Config {
  applicationConfig: Config_ApplicationConfig | undefined;
  networkConfig: Config_NetworkConfig | undefined;
  lastupdated: string;
  updatedby: string;
  bootstrap: string;
}

export interface Config_ApplicationConfig {
  GithubToken: string;
  MarketplaceOwner: string;
  MarketplaceUser: string;
  Project: string;
  Venue: string;
}

export interface Config_NetworkConfig {
  publicsubnets: string[];
  privatesubnets: string[];
}

export interface Parameters {
  parameterlist: { [key: string]: Parameters_Parameter };
}

export interface Parameters_Parameter {
  name: string;
  value: string;
  type: string;
  tracked: boolean;
  insync: boolean;
}

export interface Parameters_ParameterlistEntry {
  key: string;
  value: Parameters_Parameter | undefined;
}

export interface LogLine {
  line: string;
  level: string;
  timestamp: string;
  type: string;
}

function createBaseUnityWebsocketMessage(): UnityWebsocketMessage {
  return {
    install: undefined,
    simplemessage: undefined,
    connectionsetup: undefined,
    config: undefined,
    parameters: undefined,
    logs: undefined,
    deployments: undefined,
    uninstall: undefined,
  };
}

export const UnityWebsocketMessage = {
  encode(message: UnityWebsocketMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.install !== undefined) {
      Install.encode(message.install, writer.uint32(10).fork()).ldelim();
    }
    if (message.simplemessage !== undefined) {
      SimpleMessage.encode(message.simplemessage, writer.uint32(18).fork()).ldelim();
    }
    if (message.connectionsetup !== undefined) {
      ConnectionSetup.encode(message.connectionsetup, writer.uint32(26).fork()).ldelim();
    }
    if (message.config !== undefined) {
      Config.encode(message.config, writer.uint32(34).fork()).ldelim();
    }
    if (message.parameters !== undefined) {
      Parameters.encode(message.parameters, writer.uint32(42).fork()).ldelim();
    }
    if (message.logs !== undefined) {
      LogLine.encode(message.logs, writer.uint32(50).fork()).ldelim();
    }
    if (message.deployments !== undefined) {
      Deployments.encode(message.deployments, writer.uint32(58).fork()).ldelim();
    }
    if (message.uninstall !== undefined) {
      Uninstall.encode(message.uninstall, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnityWebsocketMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnityWebsocketMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.install = Install.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.simplemessage = SimpleMessage.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.connectionsetup = ConnectionSetup.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.config = Config.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.parameters = Parameters.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.logs = LogLine.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.deployments = Deployments.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.uninstall = Uninstall.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UnityWebsocketMessage {
    return {
      install: isSet(object.install) ? Install.fromJSON(object.install) : undefined,
      simplemessage: isSet(object.simplemessage) ? SimpleMessage.fromJSON(object.simplemessage) : undefined,
      connectionsetup: isSet(object.connectionsetup) ? ConnectionSetup.fromJSON(object.connectionsetup) : undefined,
      config: isSet(object.config) ? Config.fromJSON(object.config) : undefined,
      parameters: isSet(object.parameters) ? Parameters.fromJSON(object.parameters) : undefined,
      logs: isSet(object.logs) ? LogLine.fromJSON(object.logs) : undefined,
      deployments: isSet(object.deployments) ? Deployments.fromJSON(object.deployments) : undefined,
      uninstall: isSet(object.uninstall) ? Uninstall.fromJSON(object.uninstall) : undefined,
    };
  },

  toJSON(message: UnityWebsocketMessage): unknown {
    const obj: any = {};
    message.install !== undefined && (obj.install = message.install ? Install.toJSON(message.install) : undefined);
    message.simplemessage !== undefined &&
      (obj.simplemessage = message.simplemessage ? SimpleMessage.toJSON(message.simplemessage) : undefined);
    message.connectionsetup !== undefined &&
      (obj.connectionsetup = message.connectionsetup ? ConnectionSetup.toJSON(message.connectionsetup) : undefined);
    message.config !== undefined && (obj.config = message.config ? Config.toJSON(message.config) : undefined);
    message.parameters !== undefined &&
      (obj.parameters = message.parameters ? Parameters.toJSON(message.parameters) : undefined);
    message.logs !== undefined && (obj.logs = message.logs ? LogLine.toJSON(message.logs) : undefined);
    message.deployments !== undefined &&
      (obj.deployments = message.deployments ? Deployments.toJSON(message.deployments) : undefined);
    message.uninstall !== undefined &&
      (obj.uninstall = message.uninstall ? Uninstall.toJSON(message.uninstall) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<UnityWebsocketMessage>, I>>(base?: I): UnityWebsocketMessage {
    return UnityWebsocketMessage.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<UnityWebsocketMessage>, I>>(object: I): UnityWebsocketMessage {
    const message = createBaseUnityWebsocketMessage();
    message.install = (object.install !== undefined && object.install !== null)
      ? Install.fromPartial(object.install)
      : undefined;
    message.simplemessage = (object.simplemessage !== undefined && object.simplemessage !== null)
      ? SimpleMessage.fromPartial(object.simplemessage)
      : undefined;
    message.connectionsetup = (object.connectionsetup !== undefined && object.connectionsetup !== null)
      ? ConnectionSetup.fromPartial(object.connectionsetup)
      : undefined;
    message.config = (object.config !== undefined && object.config !== null)
      ? Config.fromPartial(object.config)
      : undefined;
    message.parameters = (object.parameters !== undefined && object.parameters !== null)
      ? Parameters.fromPartial(object.parameters)
      : undefined;
    message.logs = (object.logs !== undefined && object.logs !== null) ? LogLine.fromPartial(object.logs) : undefined;
    message.deployments = (object.deployments !== undefined && object.deployments !== null)
      ? Deployments.fromPartial(object.deployments)
      : undefined;
    message.uninstall = (object.uninstall !== undefined && object.uninstall !== null)
      ? Uninstall.fromPartial(object.uninstall)
      : undefined;
    return message;
  },
};

function createBaseApplication(): Application {
  return { packageName: "", version: "", source: "", status: "", applicationName: "", displayName: "" };
}

export const Application = {
  encode(message: Application, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.packageName !== "") {
      writer.uint32(10).string(message.packageName);
    }
    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }
    if (message.source !== "") {
      writer.uint32(26).string(message.source);
    }
    if (message.status !== "") {
      writer.uint32(34).string(message.status);
    }
    if (message.applicationName !== "") {
      writer.uint32(42).string(message.applicationName);
    }
    if (message.displayName !== "") {
      writer.uint32(50).string(message.displayName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Application {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApplication();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.packageName = reader.string();
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

          message.source = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.status = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.applicationName = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.displayName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Application {
    return {
      packageName: isSet(object.packageName) ? String(object.packageName) : "",
      version: isSet(object.version) ? String(object.version) : "",
      source: isSet(object.source) ? String(object.source) : "",
      status: isSet(object.status) ? String(object.status) : "",
      applicationName: isSet(object.applicationName) ? String(object.applicationName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
    };
  },

  toJSON(message: Application): unknown {
    const obj: any = {};
    message.packageName !== undefined && (obj.packageName = message.packageName);
    message.version !== undefined && (obj.version = message.version);
    message.source !== undefined && (obj.source = message.source);
    message.status !== undefined && (obj.status = message.status);
    message.applicationName !== undefined && (obj.applicationName = message.applicationName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    return obj;
  },

  create<I extends Exact<DeepPartial<Application>, I>>(base?: I): Application {
    return Application.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Application>, I>>(object: I): Application {
    const message = createBaseApplication();
    message.packageName = object.packageName ?? "";
    message.version = object.version ?? "";
    message.source = object.source ?? "";
    message.status = object.status ?? "";
    message.applicationName = object.applicationName ?? "";
    message.displayName = object.displayName ?? "";
    return message;
  },
};

function createBaseDeployment(): Deployment {
  return { name: "", creator: "", creationdate: "", application: [] };
}

export const Deployment = {
  encode(message: Deployment, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.creator !== "") {
      writer.uint32(18).string(message.creator);
    }
    if (message.creationdate !== "") {
      writer.uint32(26).string(message.creationdate);
    }
    for (const v of message.application) {
      Application.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Deployment {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeployment();
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

          message.creator = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.creationdate = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.application.push(Application.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Deployment {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      creator: isSet(object.creator) ? String(object.creator) : "",
      creationdate: isSet(object.creationdate) ? String(object.creationdate) : "",
      application: Array.isArray(object?.application)
        ? object.application.map((e: any) => Application.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Deployment): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.creator !== undefined && (obj.creator = message.creator);
    message.creationdate !== undefined && (obj.creationdate = message.creationdate);
    if (message.application) {
      obj.application = message.application.map((e) => e ? Application.toJSON(e) : undefined);
    } else {
      obj.application = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Deployment>, I>>(base?: I): Deployment {
    return Deployment.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Deployment>, I>>(object: I): Deployment {
    const message = createBaseDeployment();
    message.name = object.name ?? "";
    message.creator = object.creator ?? "";
    message.creationdate = object.creationdate ?? "";
    message.application = object.application?.map((e) => Application.fromPartial(e)) || [];
    return message;
  },
};

function createBaseDeployments(): Deployments {
  return { deployment: [] };
}

export const Deployments = {
  encode(message: Deployments, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.deployment) {
      Deployment.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Deployments {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeployments();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.deployment.push(Deployment.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Deployments {
    return {
      deployment: Array.isArray(object?.deployment) ? object.deployment.map((e: any) => Deployment.fromJSON(e)) : [],
    };
  },

  toJSON(message: Deployments): unknown {
    const obj: any = {};
    if (message.deployment) {
      obj.deployment = message.deployment.map((e) => e ? Deployment.toJSON(e) : undefined);
    } else {
      obj.deployment = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Deployments>, I>>(base?: I): Deployments {
    return Deployments.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Deployments>, I>>(object: I): Deployments {
    const message = createBaseDeployments();
    message.deployment = object.deployment?.map((e) => Deployment.fromPartial(e)) || [];
    return message;
  },
};

function createBaseConnectionSetup(): ConnectionSetup {
  return { type: "", userID: "" };
}

export const ConnectionSetup = {
  encode(message: ConnectionSetup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.userID !== "") {
      writer.uint32(18).string(message.userID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ConnectionSetup {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConnectionSetup();
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

          message.userID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ConnectionSetup {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      userID: isSet(object.userID) ? String(object.userID) : "",
    };
  },

  toJSON(message: ConnectionSetup): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.userID !== undefined && (obj.userID = message.userID);
    return obj;
  },

  create<I extends Exact<DeepPartial<ConnectionSetup>, I>>(base?: I): ConnectionSetup {
    return ConnectionSetup.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ConnectionSetup>, I>>(object: I): ConnectionSetup {
    const message = createBaseConnectionSetup();
    message.type = object.type ?? "";
    message.userID = object.userID ?? "";
    return message;
  },
};

function createBaseInstall(): Install {
  return { applications: undefined, DeploymentName: "" };
}

export const Install = {
  encode(message: Install, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.applications !== undefined) {
      Install_Applications.encode(message.applications, writer.uint32(10).fork()).ldelim();
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
      DeploymentName: isSet(object.DeploymentName) ? String(object.DeploymentName) : "",
    };
  },

  toJSON(message: Install): unknown {
    const obj: any = {};
    message.applications !== undefined &&
      (obj.applications = message.applications ? Install_Applications.toJSON(message.applications) : undefined);
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
    message.DeploymentName = object.DeploymentName ?? "";
    return message;
  },
};

function createBaseInstall_Variables(): Install_Variables {
  return { Values: {}, AdvancedValues: undefined };
}

export const Install_Variables = {
  encode(message: Install_Variables, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.Values).forEach(([key, value]) => {
      Install_Variables_ValuesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    if (message.AdvancedValues !== undefined) {
      Struct.encode(Struct.wrap(message.AdvancedValues), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Variables {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Variables();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = Install_Variables_ValuesEntry.decode(reader, reader.uint32());
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

  fromJSON(object: any): Install_Variables {
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

  toJSON(message: Install_Variables): unknown {
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

  create<I extends Exact<DeepPartial<Install_Variables>, I>>(base?: I): Install_Variables {
    return Install_Variables.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Variables>, I>>(object: I): Install_Variables {
    const message = createBaseInstall_Variables();
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

function createBaseInstall_Variables_ValuesEntry(): Install_Variables_ValuesEntry {
  return { key: "", value: "" };
}

export const Install_Variables_ValuesEntry = {
  encode(message: Install_Variables_ValuesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Variables_ValuesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Variables_ValuesEntry();
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

  fromJSON(object: any): Install_Variables_ValuesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: Install_Variables_ValuesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Variables_ValuesEntry>, I>>(base?: I): Install_Variables_ValuesEntry {
    return Install_Variables_ValuesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Variables_ValuesEntry>, I>>(
    object: I,
  ): Install_Variables_ValuesEntry {
    const message = createBaseInstall_Variables_ValuesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseInstall_Applications(): Install_Applications {
  return {
    name: "",
    version: "",
    variables: undefined,
    postinstall: "",
    preinstall: "",
    displayname: "",
    dependencies: {},
  };
}

export const Install_Applications = {
  encode(message: Install_Applications, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }
    if (message.variables !== undefined) {
      Install_Variables.encode(message.variables, writer.uint32(26).fork()).ldelim();
    }
    if (message.postinstall !== "") {
      writer.uint32(34).string(message.postinstall);
    }
    if (message.preinstall !== "") {
      writer.uint32(42).string(message.preinstall);
    }
    if (message.displayname !== "") {
      writer.uint32(50).string(message.displayname);
    }
    Object.entries(message.dependencies).forEach(([key, value]) => {
      Install_Applications_DependenciesEntry.encode({ key: key as any, value }, writer.uint32(58).fork()).ldelim();
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

          message.variables = Install_Variables.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.postinstall = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.preinstall = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.displayname = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          const entry7 = Install_Applications_DependenciesEntry.decode(reader, reader.uint32());
          if (entry7.value !== undefined) {
            message.dependencies[entry7.key] = entry7.value;
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
      variables: isSet(object.variables) ? Install_Variables.fromJSON(object.variables) : undefined,
      postinstall: isSet(object.postinstall) ? String(object.postinstall) : "",
      preinstall: isSet(object.preinstall) ? String(object.preinstall) : "",
      displayname: isSet(object.displayname) ? String(object.displayname) : "",
      dependencies: isObject(object.dependencies)
        ? Object.entries(object.dependencies).reduce<{ [key: string]: string }>((acc, [key, value]) => {
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
    message.variables !== undefined &&
      (obj.variables = message.variables ? Install_Variables.toJSON(message.variables) : undefined);
    message.postinstall !== undefined && (obj.postinstall = message.postinstall);
    message.preinstall !== undefined && (obj.preinstall = message.preinstall);
    message.displayname !== undefined && (obj.displayname = message.displayname);
    obj.dependencies = {};
    if (message.dependencies) {
      Object.entries(message.dependencies).forEach(([k, v]) => {
        obj.dependencies[k] = v;
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
    message.variables = (object.variables !== undefined && object.variables !== null)
      ? Install_Variables.fromPartial(object.variables)
      : undefined;
    message.postinstall = object.postinstall ?? "";
    message.preinstall = object.preinstall ?? "";
    message.displayname = object.displayname ?? "";
    message.dependencies = Object.entries(object.dependencies ?? {}).reduce<{ [key: string]: string }>(
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

function createBaseInstall_Applications_DependenciesEntry(): Install_Applications_DependenciesEntry {
  return { key: "", value: "" };
}

export const Install_Applications_DependenciesEntry = {
  encode(message: Install_Applications_DependenciesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Install_Applications_DependenciesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstall_Applications_DependenciesEntry();
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

  fromJSON(object: any): Install_Applications_DependenciesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: Install_Applications_DependenciesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create<I extends Exact<DeepPartial<Install_Applications_DependenciesEntry>, I>>(
    base?: I,
  ): Install_Applications_DependenciesEntry {
    return Install_Applications_DependenciesEntry.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Install_Applications_DependenciesEntry>, I>>(
    object: I,
  ): Install_Applications_DependenciesEntry {
    const message = createBaseInstall_Applications_DependenciesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseUninstall(): Uninstall {
  return { DeploymentName: "", Application: "", DisplayName: "", All: false };
}

export const Uninstall = {
  encode(message: Uninstall, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.DeploymentName !== "") {
      writer.uint32(10).string(message.DeploymentName);
    }
    if (message.Application !== "") {
      writer.uint32(18).string(message.Application);
    }
    if (message.DisplayName !== "") {
      writer.uint32(34).string(message.DisplayName);
    }
    if (message.All === true) {
      writer.uint32(24).bool(message.All);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Uninstall {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUninstall();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.DeploymentName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Application = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.DisplayName = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.All = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Uninstall {
    return {
      DeploymentName: isSet(object.DeploymentName) ? String(object.DeploymentName) : "",
      Application: isSet(object.Application) ? String(object.Application) : "",
      DisplayName: isSet(object.DisplayName) ? String(object.DisplayName) : "",
      All: isSet(object.All) ? Boolean(object.All) : false,
    };
  },

  toJSON(message: Uninstall): unknown {
    const obj: any = {};
    message.DeploymentName !== undefined && (obj.DeploymentName = message.DeploymentName);
    message.Application !== undefined && (obj.Application = message.Application);
    message.DisplayName !== undefined && (obj.DisplayName = message.DisplayName);
    message.All !== undefined && (obj.All = message.All);
    return obj;
  },

  create<I extends Exact<DeepPartial<Uninstall>, I>>(base?: I): Uninstall {
    return Uninstall.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Uninstall>, I>>(object: I): Uninstall {
    const message = createBaseUninstall();
    message.DeploymentName = object.DeploymentName ?? "";
    message.Application = object.Application ?? "";
    message.DisplayName = object.DisplayName ?? "";
    message.All = object.All ?? false;
    return message;
  },
};

function createBaseSimpleMessage(): SimpleMessage {
  return { operation: "", payload: "" };
}

export const SimpleMessage = {
  encode(message: SimpleMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.operation !== "") {
      writer.uint32(10).string(message.operation);
    }
    if (message.payload !== "") {
      writer.uint32(18).string(message.payload);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimpleMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimpleMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.operation = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.payload = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SimpleMessage {
    return {
      operation: isSet(object.operation) ? String(object.operation) : "",
      payload: isSet(object.payload) ? String(object.payload) : "",
    };
  },

  toJSON(message: SimpleMessage): unknown {
    const obj: any = {};
    message.operation !== undefined && (obj.operation = message.operation);
    message.payload !== undefined && (obj.payload = message.payload);
    return obj;
  },

  create<I extends Exact<DeepPartial<SimpleMessage>, I>>(base?: I): SimpleMessage {
    return SimpleMessage.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<SimpleMessage>, I>>(object: I): SimpleMessage {
    const message = createBaseSimpleMessage();
    message.operation = object.operation ?? "";
    message.payload = object.payload ?? "";
    return message;
  },
};

function createBaseConfig(): Config {
  return { applicationConfig: undefined, networkConfig: undefined, lastupdated: "", updatedby: "", bootstrap: "" };
}

export const Config = {
  encode(message: Config, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.applicationConfig !== undefined) {
      Config_ApplicationConfig.encode(message.applicationConfig, writer.uint32(10).fork()).ldelim();
    }
    if (message.networkConfig !== undefined) {
      Config_NetworkConfig.encode(message.networkConfig, writer.uint32(18).fork()).ldelim();
    }
    if (message.lastupdated !== "") {
      writer.uint32(26).string(message.lastupdated);
    }
    if (message.updatedby !== "") {
      writer.uint32(34).string(message.updatedby);
    }
    if (message.bootstrap !== "") {
      writer.uint32(42).string(message.bootstrap);
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
        case 3:
          if (tag !== 26) {
            break;
          }

          message.lastupdated = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.updatedby = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.bootstrap = reader.string();
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
      lastupdated: isSet(object.lastupdated) ? String(object.lastupdated) : "",
      updatedby: isSet(object.updatedby) ? String(object.updatedby) : "",
      bootstrap: isSet(object.bootstrap) ? String(object.bootstrap) : "",
    };
  },

  toJSON(message: Config): unknown {
    const obj: any = {};
    message.applicationConfig !== undefined && (obj.applicationConfig = message.applicationConfig
      ? Config_ApplicationConfig.toJSON(message.applicationConfig)
      : undefined);
    message.networkConfig !== undefined &&
      (obj.networkConfig = message.networkConfig ? Config_NetworkConfig.toJSON(message.networkConfig) : undefined);
    message.lastupdated !== undefined && (obj.lastupdated = message.lastupdated);
    message.updatedby !== undefined && (obj.updatedby = message.updatedby);
    message.bootstrap !== undefined && (obj.bootstrap = message.bootstrap);
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
    message.lastupdated = object.lastupdated ?? "";
    message.updatedby = object.updatedby ?? "";
    message.bootstrap = object.bootstrap ?? "";
    return message;
  },
};

function createBaseConfig_ApplicationConfig(): Config_ApplicationConfig {
  return { GithubToken: "", MarketplaceOwner: "", MarketplaceUser: "", Project: "", Venue: "" };
}

export const Config_ApplicationConfig = {
  encode(message: Config_ApplicationConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.GithubToken !== "") {
      writer.uint32(10).string(message.GithubToken);
    }
    if (message.MarketplaceOwner !== "") {
      writer.uint32(18).string(message.MarketplaceOwner);
    }
    if (message.MarketplaceUser !== "") {
      writer.uint32(26).string(message.MarketplaceUser);
    }
    if (message.Project !== "") {
      writer.uint32(34).string(message.Project);
    }
    if (message.Venue !== "") {
      writer.uint32(42).string(message.Venue);
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
        case 2:
          if (tag !== 18) {
            break;
          }

          message.MarketplaceOwner = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.MarketplaceUser = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.Project = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.Venue = reader.string();
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
    return {
      GithubToken: isSet(object.GithubToken) ? String(object.GithubToken) : "",
      MarketplaceOwner: isSet(object.MarketplaceOwner) ? String(object.MarketplaceOwner) : "",
      MarketplaceUser: isSet(object.MarketplaceUser) ? String(object.MarketplaceUser) : "",
      Project: isSet(object.Project) ? String(object.Project) : "",
      Venue: isSet(object.Venue) ? String(object.Venue) : "",
    };
  },

  toJSON(message: Config_ApplicationConfig): unknown {
    const obj: any = {};
    message.GithubToken !== undefined && (obj.GithubToken = message.GithubToken);
    message.MarketplaceOwner !== undefined && (obj.MarketplaceOwner = message.MarketplaceOwner);
    message.MarketplaceUser !== undefined && (obj.MarketplaceUser = message.MarketplaceUser);
    message.Project !== undefined && (obj.Project = message.Project);
    message.Venue !== undefined && (obj.Venue = message.Venue);
    return obj;
  },

  create<I extends Exact<DeepPartial<Config_ApplicationConfig>, I>>(base?: I): Config_ApplicationConfig {
    return Config_ApplicationConfig.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Config_ApplicationConfig>, I>>(object: I): Config_ApplicationConfig {
    const message = createBaseConfig_ApplicationConfig();
    message.GithubToken = object.GithubToken ?? "";
    message.MarketplaceOwner = object.MarketplaceOwner ?? "";
    message.MarketplaceUser = object.MarketplaceUser ?? "";
    message.Project = object.Project ?? "";
    message.Venue = object.Venue ?? "";
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
  return { name: "", value: "", type: "", tracked: false, insync: false };
}

export const Parameters_Parameter = {
  encode(message: Parameters_Parameter, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
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
      name: isSet(object.name) ? String(object.name) : "",
      value: isSet(object.value) ? String(object.value) : "",
      type: isSet(object.type) ? String(object.type) : "",
      tracked: isSet(object.tracked) ? Boolean(object.tracked) : false,
      insync: isSet(object.insync) ? Boolean(object.insync) : false,
    };
  },

  toJSON(message: Parameters_Parameter): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
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
    message.name = object.name ?? "";
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

function createBaseLogLine(): LogLine {
  return { line: "", level: "", timestamp: "", type: "" };
}

export const LogLine = {
  encode(message: LogLine, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.line !== "") {
      writer.uint32(10).string(message.line);
    }
    if (message.level !== "") {
      writer.uint32(18).string(message.level);
    }
    if (message.timestamp !== "") {
      writer.uint32(26).string(message.timestamp);
    }
    if (message.type !== "") {
      writer.uint32(34).string(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LogLine {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLogLine();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.line = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.level = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.timestamp = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.type = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LogLine {
    return {
      line: isSet(object.line) ? String(object.line) : "",
      level: isSet(object.level) ? String(object.level) : "",
      timestamp: isSet(object.timestamp) ? String(object.timestamp) : "",
      type: isSet(object.type) ? String(object.type) : "",
    };
  },

  toJSON(message: LogLine): unknown {
    const obj: any = {};
    message.line !== undefined && (obj.line = message.line);
    message.level !== undefined && (obj.level = message.level);
    message.timestamp !== undefined && (obj.timestamp = message.timestamp);
    message.type !== undefined && (obj.type = message.type);
    return obj;
  },

  create<I extends Exact<DeepPartial<LogLine>, I>>(base?: I): LogLine {
    return LogLine.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<LogLine>, I>>(object: I): LogLine {
    const message = createBaseLogLine();
    message.line = object.line ?? "";
    message.level = object.level ?? "";
    message.timestamp = object.timestamp ?? "";
    message.type = object.type ?? "";
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
