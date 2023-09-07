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
export class BaseView extends HTMLElement {
    constructor(params: any, _onRenderComplete: any);
    _title: string;
    params: any;
    _state: {};
    _onRenderComplete: any;
    router: any;
    connectedCallback(): Promise<void>;
    disconnectedCallback(): void;
    /**
     * This method is called before the view is rendered. Override this method
     * to perform any actions before the view is rendered.
     * @returns {Promise<void>}
     */
    beforeRender(): Promise<void>;
    /**
     * This method is called after the view is rendered. Override this method
     * to perform any actions after the view is rendered.
     * @returns {Promise<void>}
     */
    afterRender(): Promise<void>;
    /**
     * This method is called when the view is unloaded. Override this method
     * to perform any actions when the view is unloaded.
     * @returns {Promise<void>}
     */
    onUnload(): Promise<void>;
    /**
     * This method is called when the view is navigated to. Override this method
     * render your page contents.
     * @returns {Promise<void>}
     */
    render(): Promise<void>;
    /**
     * Get the title for the current view.
     * @returns {string}
     */
    get title(): string;
    /**
     * Get the HTML for the current view.
     * @returns {string}
     */
    get html(): string;
    /**
     * Set the state for the current view.
     * @param {object} newState The new state for the view.
     * @returns {void}
     */
    set state(arg: any);
    /**
     * Get the state for the current view.
     * @returns {object}
     */
    get state(): any;
    /**
     * Get the value of a query parameter.
     * @param {string} paramName The name of the query parameter.
     * @returns {string}
     */
    getQueryParam(paramName: string): string;
    /**
     * Navigate to a new URL.
     * @param {string} url The URL to navigate to.
     * @param {object} queryParams Query parameters to add to the URL.
     * @param {object} state The state to pass to the new view.
     * @returns {void}
     */
    navigateTo(url: string, queryParams?: object, state?: object): void;
    _setDocumentTitle(): void;
}
/**
 * DefaultPageNotFound is the default view to display when a page is not found.
 * @class DefaultPageNotFound
 * @extends {BaseView}
 */
export class DefaultPageNotFound extends BaseView {
    constructor(params: any);
    render(): Promise<string>;
}
