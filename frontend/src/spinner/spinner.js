/**
 * Spinner is a simple library for displaying a loading spinner. It makes use
 * of the whole page to display the spinner. The spinner is pure CSS, SVG, and JavaScript.
 * Copyright © 2022 App Nerds LLC
 * @class Spinner
 */
export class Spinner {
	constructor() {
		this.spinnerEl = null;
	}

	hide() {
		if (this.spinnerEl) {
			this.spinnerEl.remove();
			this.spinnerEl = null;
		}
	}

	show() {
		if (!this.spinnerEl) {
			this.spinnerEl = document.createElement("div");
			this.spinnerEl.classList.add("spinner");
			this.spinnerEl.innerHTML = `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
			 <circle cx="50" cy="50" r="45" />
		  </svg>
		`;

			document.body.appendChild(this.spinnerEl);
		}
	}
}

