/** @typedef {object & { closeOnClick: boolean, onShimClick: function }} ShimOptions */

/**
 * Shim displays a full screen shim to cover elements.
 * @param {ShimOptions} options
 */
export class Shim {
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
