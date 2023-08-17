import { Router } from "./router.js";
import { DefaultPageNotFound } from "./base-view.js";

/** @typedef {import("./router.js").Route} Route */
/** @typedef {object & {routes: Array<Route>, targetElement: HTMLElement, router: Router, afterRoute: function, beforeRoute: function, injectParams: function, onRenderComplete: function, go: function }} Application */

/**
 * Creates a new single-page application.
 * @param {HTMLElement} targetElement The element to render the SPA into.
 * @param {Array<Route>} routes The routes to use for the SPA.
 * @param {BaseView} pageNotFoundView The view to use when a route is not found.
 * @returns {Application}
 */
export const application = (
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
