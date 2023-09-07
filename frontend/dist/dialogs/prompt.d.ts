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
    windowEl: HTMLDivElement;
    shim: Shim;
    width: string;
    height: string;
    actionButtonID: string;
    cancelButtonID: string;
    /** @type {ValidatorFunc} */ validatorFunc: ValidatorFunc;
    connectedCallback(): void;
    hide(): void;
    show(): void;
    /**
     * Add a validation function to the prompter. This function will be called when
     * the confirm button is clicked.
     * @param {ValidatorFunc} f
     * @returns {void}
     */
    addValidatorFunc(f: ValidatorFunc): void;
    _onCancelClick(): void;
    _onConfirmClick(): void;
    _renderWindow(): void;
    _clearAllInputs(): void;
}
/**
 * Callback used to validate the values entered into a Prompter
 */
export type ValidatorFunc = (promptValues: any) => any;
import { Shim } from "../shim/shim.js";
