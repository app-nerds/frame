/**
 * A wrapper around fetch that will show a spinner
 * while the request is being made. This is useful for
 * long running requests.
 * @param {string} url The URL to fetch
 * @param {object} options The fetch options
 * @param {object} spinner The spinner element to show
 * @param {number} msBeforeShowSpinner The number of milliseconds to wait before showing the spinner. Default is 1000
 * @returns {Promise<object>} A promise that resolves to the fetch response
 */
export async function fetcher(url, options, spinner, msBeforeShowSpinner = 1000) {
	let timerID;
	let response;

	if (spinner) {
		timerID = setTimeout(() => {
			spinner.show();
		}, msBeforeShowSpinner);
	}

	try {
		response = await fetch(url, options);
	} finally {
		if (spinner) {
			clearTimeout(timerID);
			spinner.hide();
		}
	}

	return response;
}

