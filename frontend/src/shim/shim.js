/*
 * shim displays a full screen shim to cover elements. CSS is provided by shim.css.
 * Copyright © 2022 App Nerds LLC
 */

export function shim(baseOptions = {
	closeOnClick: false,
}) {
	/*
	 * internal constructor function to build a new shim
	 */
	function newShim(options) {
		options = { ...baseOptions, ...options };
		let shim = undefined;

		function show() {
			if (!shim && !document.getElementsByClassName("shim").length) {
				shim = document.createElement("div");
				shim.classList.add("shim");

				if (options.closeOnClick) {
					shim.addEventListener("click", () => {
						destroy(true);
					});
				}

				document.body.appendChild(shim);
			} else if (document.getElementsByClassName("shim").length) {
				shim = document.getElementsByClassName("shim")[0];
			}
		}

		function destroy(fireCallback = true) {
			if (shim) {
				shim.remove();
				shim = undefined;

				if (typeof options.callback === "function" && fireCallback) {
					options.callback();
				}
			}
		}

		return {
			hide(fireCallback = true) {
				destroy(fireCallback);
			},

			show() {
				show();
			},
		}
	}

	/*
	 * Public functions
	 */
	return {
		new(options) {
			options = { ...{ callback: undefined }, ...options };
			return newShim(options);
		}
	};
}

