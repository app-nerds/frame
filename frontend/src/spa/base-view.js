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
export class DefaultPageNotFound extends BaseView {
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
