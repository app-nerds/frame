export const DateFormats = {
	"IsoWithTimezone": "yyyy-mm-ddthh:mm:ssz",
	"IsoWithoutTimezone": "yyyy-mm-ddthh:mm:ss",
	"InternationalWithSeconds": "yyyy-mm-dd hh:mm:ss",
	"International": "yyyy-mm-dd hh:mm",
	"UsDateTimeWithSeconds": "mm/dd/yyyy hh:mm:ss am/pm",
	"UsDateTimeWithoutSeconds": "mm/dd/yyyy hh:mm am/pm",
	"IsoDate": "yyyy-mm-dd",
	"UsDate": "mm/dd/yyyy",
}

const formatMap = {
	"yyyy-mm-ddthh:mm:ssz": format1,
	"yyyy-mm-ddthh:mm:ss": format2,
	"yyyy-mm-dd hh:mm:ss": format3,
	"yyyy-mm-dd hh:mm": format4,
	"mm/dd/yyyy hh:mm:ss am/pm": format5,
	"mm/dd/yyyy hh:mm am/pm": format6,
	"yyyy-mm-dd": format7,
	"mm/dd/yyyy": format8,
}

/**
 * parseDateTime parses a date/time string into a Date object.
 * @param {string|number|Date} dt
 * @returns {Date}
 */
export function parseDateTime(dt) {
	if (typeof dt === "number") {
		return new Date(dt);
	}

	if (dt instanceof Date) {
		return dt;
	}

	if (typeof dt === "string") {
		return parseDateString(dt);
	}
}

/**
 * formatDateTime formats a Date object into a string using the specified format.
 * @param {string|number|Date} dt
 * @param {string} format
 * @returns {string}
 */
export function formatDateTime(dt, format) {
	let date = parseDateTime(dt);
	let formatter = formatMap[format.toLowerCase()];

	if (!formatter) {
		throw new Error(`Invalid date format: ${format}`);
	}

	return formatter(date);
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format1(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}T${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}Z`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format2(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}T${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format3(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)} ${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format4(dt) {
	let { year, month, day, hour, minute } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)} ${zeroPad(hour)}:${zeroPad(minute)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format5(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	let meridian = "AM";

	if (hour > 12) {
		meridian = "PM";
		hour -= 12;
	}

	if (hour === 0) {
		hour = 12;
	}

	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)} ${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)} ${meridian}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format6(dt) {
	let { year, month, day, hour, minute } = breakDownDate(dt);
	let meridian = "AM";

	if (hour > 12) {
		meridian = "PM";
		hour -= 12;
	}

	if (hour === 0) {
		hour = 12;
	}

	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)} ${zeroPad(hour)}:${zeroPad(minute)} ${meridian}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format7(dt) {
	let { year, month, day } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format8(dt) {
	let { year, month, day } = breakDownDate(dt);
	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)}`;
}

/**
 * @param {Date} dt
 * @returns {object}
 */
function breakDownDate(dt) {
	let year = dt.getFullYear();
	let month = dt.getMonth() + 1;
	let day = dt.getDate();
	let hour = dt.getHours();
	let minute = dt.getMinutes();
	let second = dt.getSeconds();

	return {
		year,
		month,
		day,
		hour,
		minute,
		second,
	};
}

/**
 * @param {number} num
 * @returns {string}
 */
function zeroPad(num) {
	return num.toString().padStart(2, "0");
}

/**
 * @param {string} dt
 * @returns {Date}
 */
function parseDateString(dt) {
	const formatRegexes = [
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})t(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})z/i,                         /* YYYY-MM-DDTHH-mm-ssZ */
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})t(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})/i,                          /* YYYY-MM-DDTHH-mm-ss */
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})\s+(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})/i,                        /* YYYY-MM-DD HH-mm-ss */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})\s+(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})\s*(?<meridian>am|pm)/i, /* MM/DD/YYYY HH:mm:ss am/pm */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})\s+(?<hour>\d{2}):(?<minute>\d{2})\s*(?<meridian>am|pm)/i,                  /* MM/DD/YYYY HH:mm am/pm */
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})/,                                                                            /* YYYY-MM-DD */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})/,                                                                          /* MM/DD/YYYY */
	];

	for (let i = 0; i < formatRegexes.length; i++) {
		let match = dt.match(formatRegexes[i]);

		if (match) {
			let year = parseInt(match.groups.year) || new Date().getFullYear();
			let month = parseInt(match.groups.month) || new Date().getMonth();
			let day = parseInt(match.groups.day) || new Date().getDate();
			let hour = parseInt(match.groups.hour) || 0;
			let minute = parseInt(match.groups.minute) || 0;
			let second = parseInt(match.groups.second) || 0;
			let meridian = match.groups.meridian || "";

			if (meridian !== "") {
				if (meridian.toLowerCase() === "pm" && hour < 12) {
					hour += 12;
				}
			}

			return new Date(year, month - 1, day, hour, minute, second);
		}
	}

	throw new Error(`no pattern match for ${dt}`);
}

