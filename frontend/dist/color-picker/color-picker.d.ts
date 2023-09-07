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
    _color: string;
    _colors: string;
    _name: string;
    input: any;
    _createOuterContainer(): any;
    _createColorGrid(colors: any, selectedColor: any): any;
    _createColorItem(color: any, selectedColor: any): any;
    _createInput(name: any, color: any): any;
    _onColorItemClicked(e: any): void;
    _clearGridSelectedClasses(): void;
}
