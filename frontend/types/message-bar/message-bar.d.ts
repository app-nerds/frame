/**
 * MessageBar is a component used to display a message on the screen.
 * @class MessageBar
 * @extends {HTMLElement}
 * @property {string} messageType The type of message to display. Valid values are "error", "warn", "info", and "success".
 * @property {string} message The message to display.
 */
export default class MessageBar extends HTMLElement {
    messageType: string;
    message: string;
    containerEl: HTMLDivElement;
    connectedCallback(): void;
    _createContainerEl(): HTMLDivElement;
    _createCloseButtonEl(): HTMLSpanElement;
    _createTextEl(): HTMLParagraphElement;
}
