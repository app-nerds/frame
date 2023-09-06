import { Shim } from "../shim/shim.js";

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
export class Confirmer {
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
		let shim = new Shim(true, () => { this._close(container, callback, false); });

		const container = Object.assign(document.createElement("dialog"), {
			className: "confirm-container",
			innerHTML: `<p>${message}</p>`,
		});

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
				const noB = Object.assign(document.createElement("button"), {
					innerText: "No",
					className: "cancel-button",
				});

				noB.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback, false)
				});

				const yesB = Object.assign(document.createElement("button"), {
					innerText: "Yes",
					className: "action-button",
				});

				yesB.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback, true);
				});

				buttons.push(noB, yesB);
				break;

			default:
				const b = Object.assign(document.createElement("button"), {
					innerText: "Close",
					className: "action-button",
				});

				b.addEventListener("click", (e) => {
					e.preventDefault();
					e.stopPropagation();

					shim.hide(false);
					this._close(container, callback);
				});

				buttons.push(b);
				break;
		}

		const buttonContainer = Object.assign(document.createElement("div"), { className: "button-row" });
		buttonContainer.append(...buttons);
		container.appendChild(buttonContainer);
	}
}
