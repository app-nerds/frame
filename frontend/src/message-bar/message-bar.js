/**
 * MessageBar is a component used to display a message on the screen.
 * @class MessageBar
 * @extends {HTMLElement}
 * @property {string} messageType The type of message to display. Valid values are "error", "warn", "info", and "success".
 * @property {string} message The message to display.
 */
export default class MessageBar extends HTMLElement {
	constructor() {
		super();

		this.messageType = this.getAttribute("message-type") || "info";
		this.message = this.getAttribute("message") || "";

		this.containerEl = null;
	}

	connectedCallback() {
		this.containerEl = this._createContainerEl();
		const closeButtonEl = this._createCloseButtonEl();
		const textEl = this._createTextEl();

		this.containerEl.insertAdjacentElement("beforeend", closeButtonEl);
		this.containerEl.insertAdjacentElement("beforeend", textEl);

		this.insertAdjacentElement("beforeend", this.containerEl);
	}

	_createContainerEl() {
		const el = document.createElement("div");
		el.classList.add("message-bar");

		switch (this.messageType) {
			case "error":
				el.classList.add("message-bar-error");
				break;

			case "warn":
				el.classList.add("message-bar-warn");
				break;

			case "info":
				el.classList.add("message-bar-info");
				break;

			case "success":
				el.classList.add("message-bar-success");
				break;
		}

		return el;
	}

	_createCloseButtonEl() {
		const el = document.createElement("span");
		el.innerHTML = "&times;";

		el.addEventListener("click", () => {
			if (this.containerEl) {
				this.containerEl.remove();
			}
		});

		return el;
	}

	_createTextEl() {
		const el = document.createElement("p");
		el.setAttribute("role", "alert");
		el.innerHTML = this.message;

		return el;
	}
}

if (!customElements.get("message-bar")) {
	customElements.define("message-bar", MessageBar);
}
