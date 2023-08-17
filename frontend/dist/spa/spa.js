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
 * BaseView is the base class for all views in the application. It provides
 * a common set of functionality that all views can use. Your view JavaScript
 * components should extend this class and register themselves as custom elements.
 * @class BaseView
 * @extends {HTMLElement}
 * @property {string} title The title of the view. This is used to set the document title.
 * @property {object} params The parameters passed to the view.
 * @property {object} state The state of the view.
 */
class BaseView extends HTMLElement {
	constructor(params, _onRenderComplete) {
		super();

		this._title = "";
		this.params = params;
		this._state = {};
		this._onRenderComplete = window._router.onRenderComplete || null;

		this.router = window._router;
	}

	async connectedCallback() {
		await this.beforeRender();
		await this.render();
		this._setDocumentTitle();
		await this.afterRender();

		if (this._onRenderComplete) {
			this._onRenderComplete(this);
		}
	}

	disconnectedCallback() {
		this.onUnload();
	}

	/**
	 * This method is called before the view is rendered. Override this method
	 * to perform any actions before the view is rendered.
	 * @returns {Promise<void>}
	 */
	async beforeRender() { }

	/**
	 * This method is called after the view is rendered. Override this method
	 * to perform any actions after the view is rendered.
	 * @returns {Promise<void>}
	 */
	async afterRender() { }

	/**
	 * This method is called when the view is unloaded. Override this method
	 * to perform any actions when the view is unloaded.
	 * @returns {Promise<void>}
	 */
	async onUnload() { }

	/**
	 * This method is called when the view is navigated to. Override this method
	 * render your page contents.
	 * @returns {Promise<void>}
	 */
	async render() {
		throw new Error("not implemented");
	}

	/**
	 * Get the title for the current view.
	 * @returns {string}
	 */
	get title() {
		return this._title;
	}

	/**
	 * Get the HTML for the current view.
	 * @returns {string}
	 */
	get html() {
		return this._html;
	}

	/**
	 * Get the state for the current view.
	 * @returns {object}
	 */
	get state() {
		return this._state;
	}

	/**
	 * Set the state for the current view.
	 * @param {object} newState The new state for the view.
	 * @returns {void}
	 */
	set state(newState) {
		this._state = newState;
	}

	/**
	 * Get the value of a query parameter.
	 * @param {string} paramName The name of the query parameter.
	 * @returns {string}
	 */
	getQueryParam(paramName) {
		return this.router.getQueryParam(paramName);
	}

	/**
	 * Navigate to a new URL.
	 * @param {string} url The URL to navigate to.
	 * @param {object} queryParams Query parameters to add to the URL.
	 * @param {object} state The state to pass to the new view.
	 * @returns {void}
	 */
	navigateTo(url, queryParams = {}, state = {}) {
		this.router.navigateTo(url, queryParams, state);
	}

	_setDocumentTitle() {
		let titles = this.querySelectorAll("title");

		if (titles && titles.length > 0) {
			this._title = titles[0].innerText;
			document.title = this._title;
			this.removeChild(titles[0]);
			return;
		}
	}
}

/**
 * DefaultPageNotFound is the default view to display when a page is not found.
 * @class DefaultPageNotFound
 * @extends {BaseView}
 */
class DefaultPageNotFound extends BaseView {
	constructor(params) {
		super(params);
	}

	async render() {
		return `
			<title>Page Not Found</title>
			<p>The page ${this.params.path} could not be found.</p>
		`;
	}
}

if (!customElements.get("default-page-not-found")) {
	customElements.define("default-page-not-found", DefaultPageNotFound);
}

/** @typedef {object & { path: string, view: BaseView }} Route */

/**
 * Router is responsible for routing requests to the correct view.
 * @class Router
 */
class Router {
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

/** @typedef {import("./router.js").Route} Route */
/** @typedef {object & {routes: Array<Route>, targetElement: HTMLElement, router: Router, afterRoute: function, beforeRoute: function, injectParams: function, onRenderComplete: function, go: function }} Application */

/**
 * Creates a new single-page application.
 * @param {HTMLElement} targetElement The element to render the SPA into.
 * @param {Array<Route>} routes The routes to use for the SPA.
 * @param {BaseView} pageNotFoundView The view to use when a route is not found.
 * @returns {Application}
 */
const application = (
	targetElement,
	routes,
	pageNotFoundView = DefaultPageNotFound
) => {
	window._router = new Router(targetElement, routes, pageNotFoundView);
	window.navigateTo = window._router.navigateTo.bind(window._router);

	window.addEventListener("popstate", (e) => {
		window._router._route({
			state: e.state,
		});
	});

	return {
		routes: routes,
		targetElement: targetElement,
		router: window._router,

		afterRoute: (f) => {
			window._router.afterRoute = f.bind(window._router);
		},

		beforeRoute: (f) => {
			window._router.beforeRoute = f.bind(window._router);
		},

		injectParams: (f) => {
			window._router.injectParams = f.bind(window._router);
		},

		onRenderComplete: (f) => {
			window._router.onRenderComplete = f.bind(window._router);
		},

		go: () => {
			window._router._route({});
		},
	};
};

export { application };
