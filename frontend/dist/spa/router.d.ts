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
    constructor(targetEl: string, routes: Array<Route>, pageNotFoundView?: BaseView);
    targetEl: string;
    routes: any[];
    pageNotFoundView: BaseView;
    beforeRoute: any;
    afterRoute: any;
    injectParams: any;
    onRenderComplete: any;
    /**
     * Retrieves a query parameter from the URL by name.
     * @param {string} paramName The name of the query parameter to retrieve.
     * @returns {string}
     */
    getQueryParam(paramName: string): string;
    /**
     * Navigates to a URL.
     * @param {string} url The URL to navigate to.
     * @param {object} queryParams Query parameters to add to the URL.
     * @param {object} state The state to pass to the new view.
     * @returns {void}
     */
    navigateTo(url: string, queryParams?: object, state?: object): void;
    _pathToRegex(path: any): RegExp;
    _getParams(match: any): {};
    _route(e: any): Promise<void>;
}
export type Route = object & {
    path: string;
    view: BaseView;
};
