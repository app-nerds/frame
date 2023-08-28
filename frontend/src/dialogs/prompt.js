import { Shim } from "../shim/shim.js";

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
export class Prompter extends HTMLElement {
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
