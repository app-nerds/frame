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
    /**
     * confirm displays a confirmation dialog. It shows a message and a Close button.
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    confirm(message: string, callback: Function): void;
    /**
     * yesNo displays a confirmation dialog. It shows a message and Yes and No buttons.
     * @param {string} message
     * @returns {Promise<boolean>}
     */
    yesNo(message: string): Promise<boolean>;
    /**
     * show displays a confirmation dialog. This is a raw function that is normally
     * used by the yesNo and confirm functions.
     * @param {string} type
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    show(type: string, message: string, callback: Function): void;
    _close(container: any, callback: any, callbackValue: any): void;
    _addButtons(container: any, type: any, shim: any, callback: any): void;
}
export type ConfirmOptions = object & {
    callback: Function;
};
