/**
 * ColorPicker is a component used to display a color picker on the screen.
 * If the color the user wants is not there, they can type a hex code into the box to get
 * the color they want.
 * @class ColorPicker
 * @extends {HTMLElement}
 * @property {string} color The currently selected color.
 * @property {string} colors A comma-separated list of colors to display in the grid. These must be valid hex codes.
 * @property {string} name The name of the input field.
 */
export default class ColorPicker extends HTMLElement {
	constructor() {
		super();

		this._color = this.getAttribute("color") || "";
		this._colors = this.getAttribute("colors") || "#ffffff,#858585,#000000,#fc1303,#8f0b01,#fc5e03,#943701,#fcc600,#8f7000,#37fc00,#1e8701,#03fcdf,#018778,#05c5fa,#017291,#0349fc,#002582,#7e00fc,#47018c,#fc03f4,#8a0085,#fa009a,#8a0055";
		this._name = this.getAttribute("name") || "color";

		const colorOptions = this._colors.split(",");

		const outerContainer = this._createOuterContainer();
		const colorGrid = this._createColorGrid(colorOptions, this._color);
		this.input = this._createInput(this._name, this._color);

		outerContainer.insertAdjacentElement("beforeend", colorGrid);
		outerContainer.insertAdjacentElement("beforeend", this.input);
		this.appendChild(outerContainer);
	}

	_createOuterContainer() {
		return Object.assign(document.createElement("div"), { className: "color-picker" });
	}

	_createColorGrid(colors, selectedColor) {
		const grid = Object.assign(document.createElement("div"), { className: "grid" });

		colors.forEach(color => {
			const el = this._createColorItem(color, selectedColor);
			grid.insertAdjacentElement("beforeend", el);
		});

		return grid;
	}

	_createColorItem(color, selectedColor) {
		const el = Object.assign(document.createElement("div"), {
			className: "grid-item",
			style: `background-color: ${color}`,
		});

		el.setAttribute("data-color", color);

		if (selectedColor === color) {
			el.classList.add("grid-item-selected");
		}

		el.addEventListener("click", this._onColorItemClicked.bind(this));
		return el;
	}

	_createInput(name, color) {
		const el = Object.assign(document.createElement("input"), {
			type: "text",
			name: name,
			"aria-label": "Selected color hexidecimal value",
			autocomplete: "on",
			className: "color-input",
			value: color,
		});

		return el;
	}

	_onColorItemClicked(e) {
		this._clearGridSelectedClasses();

		const color = e.target.getAttribute("data-color");
		e.target.classList.add("grid-item-selected");

		this.input.value = color;
		this.dispatchEvent(new CustomEvent("color-selected", { detail: color }));
	}

	_clearGridSelectedClasses() {
		const gridItems = document.querySelectorAll(".grid-item");

		for (let i = 0; i < gridItems.length; i++) {
			gridItems[i].classList.remove("grid-item-selected");
		}
	}
}

if (!customElements.get("color-picker")) {
	customElements.define("color-picker", ColorPicker);
}
