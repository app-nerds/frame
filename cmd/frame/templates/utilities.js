/* Copyright © 2023 App Nerds LLC v1.4.0 */
/**
 * Converts a classic JS object to a Map
 * Copyright © 2023 App Nerds LLC
 * @param o object The object to convert
 * @returns {Map} A Map
 */
const objectToMap = (o = {}) => {
	let result = new Map();

	for (const key in o) {
		result.set(key, o[key]);
	}

	return result;
};

/**
 * Debounces a function call. This is useful for things like
 * search boxes where you don't want to make a call to the
 * server for every keystroke.
 * Copyright © 2023 App Nerds LLC
 * @param {function} fn The function to debounce
 * @param {number} delay The delay in milliseconds. Default is 400
 * @returns {function} The debounced function
 */
const debounce = (fn, delay = 400) => {
	let id = null;

	return function() {
		let args = arguments;

		clearTimeout(id);

		id = setTimeout(() => {
			fn.apply(this, args);
		}, delay);
	};
};

export { debounce, objectToMap };
