/*
 * Copyright © 2022 App Nerds LLC
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

  _setDocumentTitle() {
    let titles = this.querySelectorAll("title");

    if (titles && titles.length > 0) {
      this._title = titles[0].innerText;
      document.title = this._title;
      this.removeChild(titles[0]);
      return;
    }
  }

  async beforeRender() { }
  async afterRender() { }
  async onUnload() { }

  async render() {
    throw new Error("not implemented");
  }

  get title() {
    return this._title;
  }

  get html() {
    return this._html;
  }

  get state() {
    return this._state;
  }

  set state(newState) {
    this._state = newState;
  }

  getQueryParam(paramName) {
    return this.router.getQueryParam(paramName);
  }

  navigateTo(url, queryParams = {}, state = {}) {
    this.router.navigateTo(url, queryParams, state);
  }
}

// Used when a route cannot be found
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
