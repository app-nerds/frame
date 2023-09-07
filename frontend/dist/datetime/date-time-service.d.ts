/**
 * parseDateTime parses a date/time string into a Date object.
 * @param {string|number|Date} dt
 * @returns {Date}
 */
export function parseDateTime(dt: string | number | Date): Date;
/**
 * formatDateTime formats a Date object into a string using the specified format.
 * @param {string|number|Date} dt
 * @param {string} format
 * @returns {string}
 */
export function formatDateTime(dt: string | number | Date, format: string): string;
export namespace DateFormats {
    let IsoWithTimezone: string;
    let IsoWithoutTimezone: string;
    let InternationalWithSeconds: string;
    let International: string;
    let UsDateTimeWithSeconds: string;
    let UsDateTimeWithoutSeconds: string;
    let IsoDate: string;
    let UsDate: string;
}
