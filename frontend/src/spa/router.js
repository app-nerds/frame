import { objectToMap } from "../utilities/objectToMap.js";
import { DefaultPageNotFound } from "./base-view.js";

/** @typedef {object & { path: string, view: BaseView }} Route */

/**
 * Router is responsible for routing requests to the correct view.
 * @class Router
 */
export class Router {
	/**
	 * Creates a new instance of Router.
	 * @param {string} targetEl The element to render the SPA into.
	 * @param {Array<Route>} routes The routes to use for the SPA.
	 * @param {BaseView} pageNotFoundView The view to use when a route is not found.
	 */
	constructor(targetEl, routes, pageNotFoundView = null) {
		this.targetEl = targetEl;
		this.routes = routes;
		this.pageNotFoundView = pageNotFoundView;

		this.beforeRoute = null;
		this.afterRoute = null;
		this.injectParams = null;
		this.onRenderComplete = null;

		if (this.pageNotFoundView) {
			this.routes.push({
				path: "/404notfound/:path",
				view: this.pageNotFoundView,
			});
		} else {
			this.routes.push({
				path: "/404notfound/:path",
				view: DefaultPageNotFound,
			});
		}
	}

	/**
	 * Retrieves a query parameter from the URL by name.
	 * @param {string} paramName The name of the query parameter to retrieve.
	 * @returns {string}
	 */
	getQueryParam(paramName) {
		let params = new URLSearchParams(location.search);
		return params.get(paramName);
	}

	/**
	 * Navigates to a URL.
	 * @param {string} url The URL to navigate to.
	 * @param {object} queryParams Query parameters to add to the URL.
	 * @param {object} state The state to pass to the new view.
	 * @returns {void}
	 */
	navigateTo(url, queryParams = {}, state = {}) {
		let q = "";

		if (Object.keys(queryParams).length > 0) {
			let m = objectToMap(queryParams);
			q += "?";

			for (const [key, value] of m) {
				let encodedKey = encodeURIComponent(key);
				let jsonValue = value;

				if (typeof value === "object") {
					jsonValue = JSON.stringify(value);
				}

				let encodedValue = encodeURIComponent(jsonValue);

				q += `${encodedKey}=${encodedValue}&`;
			}
		}

		history.pushState(state, null, `${url}${q}`);
		this._route({
			state: state,
		});
	}

	_pathToRegex(path) {
		return new RegExp(
			"^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$"
		);
	}

	_getParams(match) {
		let index = 0;

		const values = match.result.slice(1);
		const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(
			(result) => result[1]
		);

		let result = {};

		for (index = 0; index < values.length; index++) {
			result[keys[index]] = values[index];
		}

		if (this.injectParams) {
			const whatToInject = this.injectParams(match);

			for (const key in whatToInject) {
				result[key] = whatToInject[key];
			}
		}

		return result;
	}

	async _route(e) {
		let state = {};

		if (e.state) {
			state = e.state;
		}

		const potentialMatches = this.routes.map((route) => {
			return {
				route,
				result: location.pathname.match(this._pathToRegex(route.path)),
			};
		});

		let match = potentialMatches.find(
			(potentialMatch) => potentialMatch.result !== null
		);

		/*
		 * Route not found - return first route
		 */
		if (!match) {
			this.navigateTo(`/404notfound${location.pathname}`);
			return;
		}

		if (this.beforeRoute) {
			if (this.beforeRoute.apply(this, match.route) === false) {
				return;
			}
		}

		/*
		 * Get parameters, then initialie the view and render.
		 */
		const params = this._getParams(match);
		const view = new match.route.view(params);
		view.state = state;

		const el = document.querySelector(this.targetEl);
		el.innerHTML = "";
		el.appendChild(view);

		if (this.afterRoute) {
			this.afterRoute(match.route);
		}
	}
}
