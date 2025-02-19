import semver from 'semver';

type SpecialVersion = 'latest' | 'current' | 'stable';
type Version = string | SpecialVersion;

interface SemverResult {
  major: number;
  minor: number;
  patch: number;
  prerelease: string[];
  build: string[];
  version: string;
}

export default class ExtendedSemver {
  private readonly specialVersions: Record<SpecialVersion, string> = {
    latest: '999999.999999.999999',
    current: '999999.999999.999998',
    stable: '999999.999999.999997'
  };

  /**
   * Converts non-standard version to semver-compatible string
   * @param version The version string to normalize
   * @returns A semver-compatible version string
   * @throws Error if version is invalid
   */
  private normalize(version: Version): string {
    const versionLower = version.toLowerCase();

    // Handle special keywords
    if (versionLower in this.specialVersions) {
      return this.specialVersions[versionLower as SpecialVersion];
    }

    // Handle suffixed versions (e.g., "24.4-stable")
    if (version.includes('-')) {
      const [base, suffix] = version.split('-');
      const parts = base.split('.');

      // Pad with zeros for consistent comparison
      while (parts.length < 3) {
        parts.push('0');
      }

      // Handle stable suffix specially since it should be higher than base
      if (suffix.toLowerCase() === 'stable') {
        return parts.join('.');
      }

      return `${parts.join('.')}-${suffix}`;
    }

    // Handle versions without enough parts
    if (!version.includes('.')) {
      return `${version}.0.0`;
    }

    const parts = version.split('.');
    while (parts.length < 3) {
      parts.push('0');
    }

    return parts.join('.');
  }

  /**
   * Compares two versions
   * @returns -1 if v1 < v2, 0 if v1 === v2, 1 if v1 > v2
   */
  compare(version1: Version, version2: Version): number {
    const v1 = this.normalize(version1);
    const v2 = this.normalize(version2);

    return semver.compare(v1, v2);
  }

  /**
   * Sorts an array of versions in ascending order
   */
  sort(versions: Version[]): Version[] {
    const normalized = versions.map((v) => ({
      original: v,
      normalized: this.normalize(v)
    }));

    return normalized
      .sort((a, b) => semver.compare(a.normalized, b.normalized))
      .map((v) => v.original);
  }

  /**
   * Checks if a version satisfies a semver range
   */
  satisfies(version: Version, range: string): boolean {
    const normalized = this.normalize(version);
    return semver.satisfies(normalized, range);
  }

  /**
   * Checks if version1 is greater than version2
   */
  gt(version1: Version, version2: Version): boolean {
    return semver.gt(this.normalize(version1), this.normalize(version2));
  }

  /**
   * Checks if version1 is less than version2
   */
  lt(version1: Version, version2: Version): boolean {
    return semver.lt(this.normalize(version1), this.normalize(version2));
  }

  /**
   * Checks if two versions are equal
   */
  eq(version1: Version, version2: Version): boolean {
    return semver.eq(this.normalize(version1), this.normalize(version2));
  }

  /**
   * Parses a version string into its components
   * @returns Parsed version object or null if invalid
   */
  parse(version: Version): SemverResult | null {
    const parsed = semver.parse(this.normalize(version));
    if (!parsed) return null;
    
    return {
      major: parsed.major,
      minor: parsed.minor,
      patch: parsed.patch,
      prerelease: [...parsed.prerelease] as string[],
      build: [...parsed.build],
      version: parsed.version
    };
  }

  /**
   * Coerces a string into a semver version, handling both standard and non-standard formats
   * Returns null if the version string is invalid or cannot be coerced
   * @example
   * coerce('v2') => '2.0.0'
   * coerce('42.6.7.9.3-alpha') => '42.6.7'
   * coerce('latest') => '999999.999999.999999'
   * coerce('v24.4-stable') => '24.4.0'
   */
  coerce(version: string): string | null {
    // Handle special versions first
    const versionLower = version.toLowerCase();
    if (versionLower in this.specialVersions) {
      return this.specialVersions[versionLower as SpecialVersion];
    }

    // Remove 'v' prefix if it exists
    version = version.replace(/^v/, '');

    // Handle suffixed versions before coercing
    if (version.includes('-')) {
      const [base] = version.split('-');
      // Use semver's coerce on the base version
      const coerced = semver.coerce(base);
      if (!coerced) return null;
      return coerced.version;
    }

    // Use semver's built-in coerce for standard version strings
    const coerced = semver.coerce(version);
    return coerced ? coerced.version : null;
  }
}
