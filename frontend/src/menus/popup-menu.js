/**
 * PopupMenu is a Web Component that displays a popup menu. It attaches to a trigger element
 * that, when clicked, will show a list of menu items. It supports icons through the Feather
 * Icons library (https://feathericons.com/).
 *
 * Styling is provided by popup-menu.css. It relies on variables:
 *   - --dialog-background-color
 *   - --prmiary-color (for the hover).
 *   - --border-color
 *
 * Usage example:
 *    <popup-menu trigger="#trigger">
 *       <popup-menu-item id="item1" text="Menu Item 1" icon="log-out"></popup-menu-item>
 *    </popup-menu>
 *
 * @class PopupMenu
 * @extends HTMLElement
 */
export class PopupMenu extends HTMLElement {
	constructor() {
		super();
		this._trigger = null;
		this.isVisible = false;
	}

	connectedCallback() {
		this._trigger = this.getAttribute("trigger");

		if (!this._trigger) {
			throw new Error(
				"You must provide a query selector for the element used to trigger this popup."
			);
		}

		this.classList.add("popup-menu-hidden");

		document.addEventListener("click", (e) => {
			if (e.target !== this && !this.contains(e.target)) {
				const t = document.querySelector(this._trigger);

				if (e.target !== t && !t.contains(e.target)) {
					this._hide();
				}
			}
		});

		document.querySelector(this._trigger).addEventListener("click", this.toggle.bind(this));
		const menuItemEls = document.querySelectorAll("popup-menu-item");

		menuItemEls.forEach((el) => {
			el.addEventListener("internal-menu-item-click", (e) => {
				this._hide();
				this.dispatchEvent(new CustomEvent("menu-item-click", {
					detail: {
						id: e.target.id,
						text: e.target.getAttribute("text"),
						data: e.target.getAttribute("data"),
					}
				}));
			});
		});
	}

	disconnectedCallback() {
		let el = document.querySelector(this._trigger);

		if (el) {
			el.removeEventListener("click", this.toggle.bind(this));
		}
	}

	/**
	* Toggles the visibility of the popup menu
	* @param {Event} e The click event
	* @returns {void}
	*/
	toggle(e) {
		if (e) {
			e.preventDefault();
		}

		if (!this.isVisible) {
			this._show();
		} else {
			this._hide();
		}
	}

	_hide() {
		this.isVisible = false;
		this.classList.add("popup-menu-hidden");
	}

	_show() {
		let triggerRect = document
			.querySelector(this._trigger)
			.getBoundingClientRect();
		let thisRect = this.getBoundingClientRect();
		let buffer = 3;

		let left = triggerRect.left;

		if (left + thisRect.width > window.innerWidth) {
			left = triggerRect.left - ((triggerRect.left + thisRect.width) - window.innerWidth) - buffer;
		}

		this.style.left = `${left}px`;
		this.style.top =
			"" + (triggerRect.y + triggerRect.height + buffer) + "px";

		this.isVisible = true;
		this.classList.remove("popup-menu-hidden");
	}
}

/**
 * Represents a popup menu item
 * @class PopupMenuItem
 * @extends HTMLElement
 */
export class PopupMenuItem extends HTMLElement {
	constructor() {
		super();
	}

	connectedCallback() {
		this._render();
	}

	_render() {
		let text = this.getAttribute("text");
		let icon = this.getAttribute("icon");
		let inner = "";

		if (icon) {
			inner += `<i class="${icon}"></i> `;
		}

		inner += text;

		const a = Object.assign(document.createElement("a"), {
			href: "javascript:void(0)",
			classList: ["popup-menu-item"],
			innerHTML: inner,
		});

		a.addEventListener("click", (e) => {
			e.preventDefault();
			e.stopPropagation();
			this.dispatchEvent(new CustomEvent("internal-menu-item-click", { detail: e }));
		});

		this.insertAdjacentElement("beforeend", a);
	}
}

/**
 * Shows a popup menu
 * @param {string} el The query selector for the popup menu
 * @returns {void}
 */
export const showPopup = (el) => {
	document.querySelector(el)._show();
};

/**
 * Hides a popup menu
 * @param {string} el The query selector for the popup menu
 */
export const hidePopup = (el) => {
	document.querySelector(el)._hide();
};

if (!customElements.get("popup-menu")) {
	customElements.define("popup-menu", PopupMenu);
}

if (!customElements.get("popup-menu-item")) {
	customElements.define("popup-menu-item", PopupMenuItem);
}

