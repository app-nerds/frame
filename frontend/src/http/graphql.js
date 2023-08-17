import { fetcher } from "./fetcher.js";

/** @typedef { object & { http: fetcher, tokenGetterFunction: function, expiredTokenCallback: function, spinner: object, navigateTo: function } } GraphQLOptions */

/**
 * This class is a wrapper around the fetcher function
 * that makes it easy to execute GraphQL queries and mutations.
 * It also handles expired tokens.
 * @class GraphQL
 * @param {string} queryURL The URL to the GraphQL API
 * @param {GraphQLOptions} options Options for the GraphQL class
 */
export class GraphQL {
	constructor(queryURL, options = {
		http: fetcher,
		tokenGetterFunction: null,
		expiredTokenCallback: null,
		spinner: null,
		navigateTo: null,
	}) {
		options = {
			http: fetcher,
			tokenGetterFunction: null,
			expiredTokenCallback: null,
			spinner: null,
			navigateTo: null,
			...options,
		};

		this.queryURL = queryURL;
		this.http = options.http;
		this.tokenGetterFunction = options.tokenGetterFunction;
		this.expiredTokenCallback = options.expiredTokenCallback;
		this.spinner = options.spinner
		this.navigateTo = options.navigateTo;
	}

	/**
	 * Executes a query against a GraphQL API
	 * @param query string A graphql query. Omit the "query {}" portion.
	 * @returns {Promise<object>} A promise that resolves to the fetch response
	 */
	async query(query) {
		if (this.expiredTokenCallback && !this.expiredTokenCallback(null, "/", this.navigateTo)) {
			return;
		}

		const token = (this.tokenGetterFunction) ? this.tokenGetterFunction() : "";

		query = `query {
			${query}
		}`;

		let options = {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({ query }),
		};

		if (token) {
			options.headers["Authorization"] = `Bearer ${token}`;
		}

		let response = await this.http(this.queryURL, options, this.spinner);

		if (response.status === 400 || response.status === 401) {
			if (this.expiredTokenCallback) {
				this.expiredTokenCallback(response, "/", this.navigateTo);
			}
			return;
		}

		let result = await response.json();

		if (!response.ok) {
			throw new Error(result.message);
		}

		return result;
	}

	/**
	 * Executes a mutation against a GraphQL API
	 * @param query string A graphql mutation. Omit the "mutation {}" portion
	 * @returns {Promise<object>} A promise that resolves to the fetch response
	 */
	async mutation(query) {
		if (this.expiredTokenCallback && !this.expiredTokenCallback(null, "/", this.navigateTo)) {
			return;
		}

		const token = (this.tokenGetterFunction) ? this.tokenGetterFunction() : "";

		query = `mutation {
			${query}
		}`;

		let options = {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({ query }),
		};

		if (token) {
			options.headers["Authorization"] = `Bearer ${token}`;
		}

		let response = await this.http(this.queryURL, options, this.spinner);

		if (response.status === 400 || response.status === 401) {
			if (this.expiredTokenCallback) {
				this.expiredTokenCallback(response, "/", this.navigateTo);
			}

			return;
		}

		let result = await response.json();

		if (!response.ok) {
			throw new Error(result.message);
		}

		return result;
	}
}

