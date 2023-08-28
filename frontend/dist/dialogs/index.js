/* Copyright © 2023 App Nerds LLC v1.4.2 */
/** @typedef {object & { position: AlertPosition, duration: number, closable: boolean, focusable: boolean }} AlertOptions */

/**
 * Constants for alert position.
 * @enum {AlertPosition}
 */
const AlertPosition = {
	TopLeft: "top-left",
	TopCenter: "top-center",
	TopRight: "top-right",
	BottomLeft: "bottom-left",
	BottomCenter: "bottom-center",
	BottomRight: "bottom-right"
};

const alertPositionIndex = [
	[AlertPosition.TopLeft, AlertPosition.TopCenter, AlertPosition.TopRight],
	[AlertPosition.BottomLeft, AlertPosition.BottomCenter, AlertPosition.BottomRight]
];

const svgs = {
	success: '<svg viewBox="0 0 426.667 426.667" width="18" height="18"><path d="M213.333 0C95.518 0 0 95.514 0 213.333s95.518 213.333 213.333 213.333c117.828 0 213.333-95.514 213.333-213.333S331.157 0 213.333 0zm-39.134 322.918l-93.935-93.931 31.309-31.309 62.626 62.622 140.894-140.898 31.309 31.309-172.203 172.207z" fill="#6ac259"></path></svg>',
	warn: '<svg viewBox="0 0 310.285 310.285" width=18 height=18> <path d="M264.845 45.441C235.542 16.139 196.583 0 155.142 0 113.702 0 74.743 16.139 45.44 45.441 16.138 74.743 0 113.703 0 155.144c0 41.439 16.138 80.399 45.44 109.701 29.303 29.303 68.262 45.44 109.702 45.44s80.399-16.138 109.702-45.44c29.303-29.302 45.44-68.262 45.44-109.701.001-41.441-16.137-80.401-45.439-109.703zm-132.673 3.895a12.587 12.587 0 0 1 9.119-3.873h28.04c3.482 0 6.72 1.403 9.114 3.888 2.395 2.485 3.643 5.804 3.514 9.284l-4.634 104.895c-.263 7.102-6.26 12.933-13.368 12.933H146.33c-7.112 0-13.099-5.839-13.345-12.945L128.64 58.594c-.121-3.48 1.133-6.773 3.532-9.258zm23.306 219.444c-16.266 0-28.532-12.844-28.532-29.876 0-17.223 12.122-30.211 28.196-30.211 16.602 0 28.196 12.423 28.196 30.211.001 17.591-11.456 29.876-27.86 29.876z" fill="#FFDA44" /> </svg>',
	info: '<svg viewBox="0 0 23.625 23.625" width=18 height=18> <path d="M11.812 0C5.289 0 0 5.289 0 11.812s5.289 11.813 11.812 11.813 11.813-5.29 11.813-11.813S18.335 0 11.812 0zm2.459 18.307c-.608.24-1.092.422-1.455.548a3.838 3.838 0 0 1-1.262.189c-.736 0-1.309-.18-1.717-.539s-.611-.814-.611-1.367c0-.215.015-.435.045-.659a8.23 8.23 0 0 1 .147-.759l.761-2.688c.067-.258.125-.503.171-.731.046-.23.068-.441.068-.633 0-.342-.071-.582-.212-.717-.143-.135-.412-.201-.813-.201-.196 0-.398.029-.605.09-.205.063-.383.12-.529.176l.201-.828c.498-.203.975-.377 1.43-.521a4.225 4.225 0 0 1 1.29-.218c.731 0 1.295.178 1.692.53.395.353.594.812.594 1.376 0 .117-.014.323-.041.617a4.129 4.129 0 0 1-.152.811l-.757 2.68a7.582 7.582 0 0 0-.167.736 3.892 3.892 0 0 0-.073.626c0 .356.079.599.239.728.158.129.435.194.827.194.185 0 .392-.033.626-.097.232-.064.4-.121.506-.17l-.203.827zm-.134-10.878a1.807 1.807 0 0 1-1.275.492c-.496 0-.924-.164-1.28-.492a1.57 1.57 0 0 1-.533-1.193c0-.465.18-.865.533-1.196a1.812 1.812 0 0 1 1.28-.497c.497 0 .923.165 1.275.497.353.331.53.731.53 1.196 0 .467-.177.865-.53 1.193z" fill="#006DF0" /> </svg>',
	error: '<svg viewBox="0 0 51.976 51.976" width=18 height=18> <path d="M44.373 7.603c-10.137-10.137-26.632-10.138-36.77 0-10.138 10.138-10.137 26.632 0 36.77s26.632 10.138 36.77 0c10.137-10.138 10.137-26.633 0-36.77zm-8.132 28.638a2 2 0 0 1-2.828 0l-7.425-7.425-7.778 7.778a2 2 0 1 1-2.828-2.828l7.778-7.778-7.425-7.425a2 2 0 1 1 2.828-2.828l7.425 7.425 7.071-7.071a2 2 0 1 1 2.828 2.828l-7.071 7.071 7.425 7.425a2 2 0 0 1 0 2.828z" fill="#D80027" /> </svg>'
};

/**
 * Alerter displays toast-like messages to users. It is inspired by vanilla-toast (
 * https://github.com/mehmetemineker/vanilla-toast)
 * @param {AlertOptions} options
 */
class Alerter {
	constructor(options = {
		position: AlertPosition.TopRight,
		duration: 3000,
		closable: true,
		focusable: true,
	}) {
		this.options = options;

		/*
		 * If the outer container doesn't exist, make it
		 */
		if (!document.getElementsByClassName("alert-container").length) {
			this._setup();
		}
	}

	/**
	 * success displays a success alert. Use this for positive messages.
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	success(message, callback) {
		this.show(message, "success", callback);
	}

	/**
	 * info displays an info alert. Use this for neutral messages.
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	info(message, callback) {
		this.show(message, "info", callback);
	}

	/**
	 * warn displays a warning alert. Use this to warn users of something.
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	warn(message, callback) {
		this.show(message, "warn", callback);
	}

	/**
	 * error displays an error alert. Use this to warn users of something bad.
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	error(message, callback) {
		this.show(message, "error", callback);
	}

	/**
	 * @param {string} message
	 * @param {string} type
	 * @param {function} callback
	 * @returns {void}
	 */
	show(message, type, callback) {
		const col = document.getElementsByClassName(this.options.position)[0];

		const card = document.createElement("div");
		card.className = `alert-card ${type}`;
		card.innerHTML += svgs[type];
		card.options = {
			...this.options, ...{
				message,
				type: type,
				yPos: this.options.position.indexOf("top") > -1 ? "top" : "bottom",
				inFocus: false,
			},
		};

		this._setContent(card);
		this._setIntroAnimation(card);
		this._bindEvents(card);
		this._autoDestroy(card, callback);

		col.appendChild(card);
	}

	_setContent(card) {
		const div = document.createElement("div");
		div.className = "text-group";

		if (card.options.title) {
			div.innerHTML = `<h4>${card.options.title}</h3>`;
		}

		div.innerHTML += `<p>${card.options.message}</p>`;
		card.appendChild(div);
	}

	/**
	 * @param {AlertCard} card
	 * @returns {void}
	 */
	_setIntroAnimation(card) {
		card.style.setProperty(`margin-${card.options.yPos}`, "-15px");
		card.style.setProperty(`opacity`, "0");

		setTimeout(() => {
			card.style.setProperty(`margin-${card.options.yPos}`, "15px");
			card.style.setProperty("opacity", "1");
		}, 50);
	}

	/**
	 * @param {AlertCard} card
	 * @returns {void}
	 */
	_bindEvents(card) {
		card.addEventListener("click", () => {
			if (card.options.closable) {
				this._destroy(card);
			}
		});

		card.addEventListener("mouseover", () => {
			card.options.inFocus = card.options.focusable;
		});

		card.addEventListener("mouseout", () => {
			card.options.inFocus = false;
			this._autoDestroy(card);
		});
	}

	/**
	 * @param {AlertCard} card
	 * @returns {void}
	 */
	_autoDestroy(card, callback) {
		if (card.options.duration !== 0) {
			setTimeout(() => {
				if (!card.options.inFocus) {
					this._destroy(card, callback);
				}
			}, card.options.duration);
		}
	}

	/**
	 * @param {AlertCard} card
	 * @returns {void}
	 */
	_destroy(card, callback) {
		card.style.setProperty(`margin-${card.options.yPos}`, `-${card.offsetHeight}px`);
		card.style.setProperty("opacity", "0");

		setTimeout(() => {
			card.remove();

			if (typeof callback === "function") {
				callback();
			}
		}, 500);
	}

	_setup() {
		const container = document.createElement("div");
		container.className = "alert-container";

		for (const rowIndex of [0, 1]) {
			const row = document.createElement("div");
			row.className = "alert-row";

			for (const colIndex of [0, 1, 2]) {
				const col = document.createElement("div");
				col.className = `alert-col ${alertPositionIndex[rowIndex][colIndex]}`;

				row.appendChild(col);
			}

			container.appendChild(row);
		}

		document.body.appendChild(container);
	}
}

/** @typedef {object & { closeOnClick: boolean, onShimClick: function }} ShimOptions */

/**
 * Shim displays a full screen shim to cover elements.
 * @param {ShimOptions} options
 */
class Shim {
	constructor(closeOnClick = false, onShimClick) {
		this.closeOnClick = closeOnClick;
		this.onShimClick = onShimClick;

		this.shim = undefined;
	}

	/**
	 * show displays the shim
	 * @returns {void}
	 */
	show() {
		if (!this.shim && !document.getElementsByClassName("shim").length) {
			this.shim = document.createElement("div");
			this.shim.classList.add("shim");

			if (this.closeOnClick) {
				this.shim.addEventListener("click", () => {
					this.hide(this.onShimClick);
				});
			}

			document.body.appendChild(this.shim);
		} else if (document.getElementsByClassName("shim").length) {
			this.shim = document.getElementsByClassName("shim")[0];
		}
	}

	/**
	 * hide removes the shim
	 * @returns {void}
	 */
	hide(callback) {
		this._destroy();

		if (typeof callback === "function") {
			callback();
		}
	}

	_destroy() {
		if (this.shim) {
			this.shim.remove();
			this.shim = undefined;
		}
	}
}

/** @typedef {object & { callback: Function }} ConfirmOptions */

/**
 * Confirmer displays a confirmation dialog. It has two mode: "yesno", "other".
 * "yesno" mode will display two buttons: Yes and No. "other" will only display a Close button.
 * The result of the click will be returned in a promise value.
 *
 * Styling is provided by confirm.css. It relies on variables:
 *   - --dialog-background-color
 *   - --border-color
 *
 * Example:
 *    const confirmer = new Confirmer();
 *    const result = await confirmer.yesNo("Are you sure?");
 */
class Confirmer {
	constructor() {
	}

	/**
	 * confirm displays a confirmation dialog. It shows a message and a Close button.
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	confirm(message, callback) {
		this.show("confirm", message, callback);
	}

	/**
	 * yesNo displays a confirmation dialog. It shows a message and Yes and No buttons.
	 * @param {string} message
	 * @returns {Promise<boolean>}
	 */
	yesNo(message) {
		return new Promise((resolve) => {
			const cb = (result) => {
				return resolve(result);
			};

			this.show("yesno", message, cb);
		});
	}

	/**
	 * show displays a confirmation dialog. This is a raw function that is normally
	 * used by the yesNo and confirm functions.
	 * @param {string} type
	 * @param {string} message
	 * @param {function} callback
	 * @returns {void}
	 */
	show(type, message, callback) {
		const container = document.createElement("dialog");
		container.classList.add("confirm-container");

		let shim = new Shim(true, () => { this._close(container, callback, false); });

		container.innerHTML += `<p>${message}</p>`;
		this._addButtons(container, type, shim, callback);

		shim.show();
		document.body.appendChild(container);
	}

	_close(container, callback, callbackValue) {
		container.remove();
		if (typeof callback === "function") {
			callback(callbackValue);
		}
	}

	_addButtons(container, type, shim, callback) {
		let buttons = [];

		switch (type) {
			case "yesno":
				const noB = document.createElement("button");
				noB.innerText = "No";
				noB.classList.add("cancel-button");
				noB.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback, false);
				});

				const yesB = document.createElement("button");
				yesB.innerText = "Yes";
				yesB.classList.add("action-button");
				yesB.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback, true);
				});

				buttons.push(noB);
				buttons.push(yesB);
				break;

			default:
				const b = document.createElement("button");
				b.innerText = "Close";
				b.classList.add("action-button");
				b.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback);
				});

				buttons.push(b);
				break;
		}

		const buttonContainer = document.createElement("div");
		buttonContainer.classList.add("button-row");

		buttons.forEach((button) => { buttonContainer.appendChild(button); });
		container.appendChild(buttonContainer);
	}
}

/**
 * Callback used to validate the values entered into a Prompter
 * @callback ValidatorFunc
 * @param {Object} promptValues - The values entered into the prompter
 * @return {Object} { validationErrors: Array, isValid: boolean }
 */

/**
 * Prompter displays a modal dialog using the contents provided in the web component slots.
 * It allows you to put whatever elements you want into the dialog, and then retrieve the
 * contents of the dialog when the user clicks the confirm button.
 * @class Prompter
 * @extends {HTMLElement}
 */
class Prompter extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: "open" });

		this.windowEl = null;
		this.shim = new Shim(false);
		this.width = this.getAttribute("width") || "";
		this.height = this.getAttribute("height") || "";
		this.actionButtonID = this.getAttribute("action-button") || "";
		this.cancelButtonID = this.getAttribute("cancel-button") || "";
		/** @type {ValidatorFunc} */ this.validatorFunc = null;

		if (!this.actionButtonID) {
			throw new Error("Prompter requires an action button ID");
		}

		if (!this.cancelButtonID) {
			throw new Error("Prompter requires a cancel button ID");
		}

		this.classList.add("hidden");

		this.shadowRoot.innerHTML = `
			<div id="window" part="prompter" role="dialog" aria-modal="true" aria-label="Prompt" style="width: ${this.width}; height: ${this.height};">
				<slot name="title"></slot>
				<slot name="body"></slot>
				<nav part="buttons">
					<slot name="buttons"></slot>
				</nav>
			</div>
		`;

	}

	connectedCallback() {
		this.querySelector(this.cancelButtonID).addEventListener("click", this._onCancelClick.bind(this));
		this.querySelector(this.actionButtonID).addEventListener("click", this._onConfirmClick.bind(this));
	}

	hide() {
		this.classList.add("hidden");
		this.shim.hide();
		this._clearAllInputs();
	}

	show() {
		this.shim.show();
		this.classList.remove("hidden");
		this.querySelector(`div[slot="body"]>input, div[slot="body"]>select, div[slot="body"]>textarea, div[slot="body"]>form>input,div[slot="body"]>form>select,div[slot="body"]>form>textarea`).focus();
	}

	/**
	 * Add a validation function to the prompter. This function will be called when
	 * the confirm button is clicked.
	 * @param {ValidatorFunc} f
	 * @returns {void}
	 */
	addValidatorFunc(f) {
		this.validatorFunc = f;
	}

	_onCancelClick() {
		this.hide();
		this.dispatchEvent(new CustomEvent("cancel"));
	}

	_onConfirmClick() {
		let result = {};

		this.querySelectorAll("input, select, textarea").forEach((el) => {
			let key = "";

			if (el.hasAttribute("name")) {
				key = el.getAttribute("name");
			} else if (el.hasAttribute("id")) {
				key = el.getAttribute("id");
			}

			result[key] = el.value;
		});

		if (this.validatorFunc) {
			const { validationErrors, isValid } = this.validatorFunc(result);

			if (!isValid) {
				this.dispatchEvent(new CustomEvent("validation-failed", {
					detail: {
						result,
						validationErrors,
					}
				}));

				return;
			}
		}

		this.hide();
		this.dispatchEvent(new CustomEvent("confirm", { detail: result }));
	}

	_renderWindow() {
		this.windowEl = document.createElement("div");
		this.windowEl.classList.add("prompter");
		this.windowEl.setAttribute("role", "dialog");
		this.windowEl.setAttribute("aria-modal", "true");
		this.windowEl.setAttribute("aria-label", "Prompt");
		this.windowEl.style.width = this.width;
		this.windowEl.style.height = this.height;

		this.windowEl.innerHTML = `
			<slot name="title"></slot>
			<slot name="body"></slot>
		`;
	}

	_clearAllInputs() {
		this.querySelectorAll("input, select, textarea").forEach((el) => {
			el.value = "";
		});
	}
}

if (!customElements.get("prompter-ui")) {
	customElements.define("prompter-ui", Prompter);
}

export { AlertPosition, Alerter, Confirmer, Prompter };
