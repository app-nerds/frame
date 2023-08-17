/**
 * Converts a classic JS object to a Map
 * Copyright © 2023 App Nerds LLC
 * @param o object The object to convert
 * @returns {Map} A Map
 */
export const objectToMap = (o = {}) => {
	let result = new Map();

	for (const key in o) {
		result.set(key, o[key]);
	}

	return result;
};

