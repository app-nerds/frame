/**
 * Constants for alert position.
 */
export type AlertPosition = AlertPosition;
export namespace AlertPosition {
    let TopLeft: string;
    let TopCenter: string;
    let TopRight: string;
    let BottomLeft: string;
    let BottomCenter: string;
    let BottomRight: string;
}
/**
 * Alerter displays toast-like messages to users. It is inspired by vanilla-toast (
 * https://github.com/mehmetemineker/vanilla-toast)
 * @param {AlertOptions} options
 */
export class Alerter {
    constructor(options?: {
        position: string;
        duration: number;
        closable: boolean;
        focusable: boolean;
    });
    options: {
        position: string;
        duration: number;
        closable: boolean;
        focusable: boolean;
    };
    /**
     * success displays a success alert. Use this for positive messages.
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    success(message: string, callback: Function): void;
    /**
     * info displays an info alert. Use this for neutral messages.
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    info(message: string, callback: Function): void;
    /**
     * warn displays a warning alert. Use this to warn users of something.
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    warn(message: string, callback: Function): void;
    /**
     * error displays an error alert. Use this to warn users of something bad.
     * @param {string} message
     * @param {function} callback
     * @returns {void}
     */
    error(message: string, callback: Function): void;
    /**
     * @param {string} message
     * @param {string} type
     * @param {function} callback
     * @returns {void}
     */
    show(message: string, type: string, callback: Function): void;
    _setContent(card: any): void;
    /**
     * @param {AlertCard} card
     * @returns {void}
     */
    _setIntroAnimation(card: AlertCard): void;
    /**
     * @param {AlertCard} card
     * @returns {void}
     */
    _bindEvents(card: AlertCard): void;
    /**
     * @param {AlertCard} card
     * @returns {void}
     */
    _autoDestroy(card: AlertCard, callback: any): void;
    /**
     * @param {AlertCard} card
     * @returns {void}
     */
    _destroy(card: AlertCard, callback: any): void;
    _setup(): void;
}
export type AlertOptions = object & {
    position: AlertPosition;
    duration: number;
    closable: boolean;
    focusable: boolean;
};
