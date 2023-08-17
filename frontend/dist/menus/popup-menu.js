/* Copyright © 2023 App Nerds LLC v1.4.0 */
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
class PopupMenu extends HTMLElement {
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

		this.style.visibility = "hidden";

		document.addEventListener("click", (e) => {
			if (e.target !== this && !this.contains(e.target)) {
				if (e.target !== document.querySelector(this._trigger)) {
					this._hide();
				} else {
					this.toggle();
				}
			}
		});

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
		this.style.visibility = "hidden";
	}

	_show() {
		let triggerRect = document
			.querySelector(this._trigger)
			.getBoundingClientRect();
		let thisRect = this.getBoundingClientRect();
		let buffer = 3;

		if (thisRect.right > window.innerWidth) {
			this.style.left =
				"" +
				(triggerRect.x + (window.innerWidth - thisRect.right) - buffer) +
				"px";
		} else {
			this.style.left = "" + triggerRect.x + "px";
		}

		this.style.top =
			"" + (triggerRect.y + triggerRect.height + buffer) + "px";

		this.isVisible = true;
		this.style.visibility = "visible";
	}
}

/**
 * Represents a popup menu item
 * @class PopupMenuItem
 * @extends HTMLElement
 */
class PopupMenuItem extends HTMLElement {
	constructor() {
		super();
	}

	connectedCallback() {
		this._render();
	}

	_render() {
		let text = this.getAttribute("text");
		let icon = this.getAttribute("icon");

		const a = document.createElement("a");
		a.href = "javascript:void(0)";
		a.classList.add("popup-menu-item");

		let inner = "";

		if (icon) {
			inner += `<i class="${icon}"></i> `;
		}

		inner += text;
		a.innerHTML = inner;

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
const showPopup = (el) => {
	document.querySelector(el).style.visibility = "visible";
};

/**
 * Hides a popup menu
 * @param {string} el The query selector for the popup menu
 */
const hidePopup = (el) => {
	document.querySelector(el).style.visibility = "hidden";
};

if (!customElements.get("popup-menu")) {
	customElements.define("popup-menu", PopupMenu);
}

if (!customElements.get("popup-menu-item")) {
	customElements.define("popup-menu-item", PopupMenuItem);
}

export { PopupMenu, PopupMenuItem, hidePopup, showPopup };
