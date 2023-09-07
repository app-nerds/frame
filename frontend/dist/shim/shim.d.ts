/** @typedef {object & { closeOnClick: boolean, onShimClick: function }} ShimOptions */
/**
 * Shim displays a full screen shim to cover elements.
 * @param {ShimOptions} options
 */
export class Shim {
    constructor(closeOnClick: boolean, onShimClick: any);
    closeOnClick: boolean;
    onShimClick: any;
    shim: Element | HTMLDivElement;
    /**
     * show displays the shim
     * @returns {void}
     */
    show(): void;
    /**
     * hide removes the shim
     * @returns {void}
     */
    hide(callback: any): void;
    _destroy(): void;
}
export type ShimOptions = object & {
    closeOnClick: boolean;
    onShimClick: Function;
};
