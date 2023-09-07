interface NestedValue {
	Config: {
		[key: string]: ConfigValue;
	};
}

interface ConfigValue {
	Options: {
		type: string;
		default: string;
	};
}

export function getConfigValue(configValue: unknown): ConfigValue {
	return configValue as ConfigValue;
}

export function getNestedValue(nestedValue: unknown): NestedValue {
	return nestedValue as NestedValue;
}

export function getEntries(value: unknown): [string, unknown][] {
	if (value != undefined) {
		return Object.entries(value as Record<string, unknown>);
	} else {
		return [];
	}
}
