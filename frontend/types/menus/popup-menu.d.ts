/**
 * PopupMenu is a Web Component that displays a popup menu. It attaches to a trigger element
 * that, when clicked, will show a list of menu items. It supports icons through the Feather
 * Icons library (https://feathericons.com/).
 *
 * Styling is provided by popup-menu.css. It relies on variables:
 *   - --dialog-background-color
 *   - --prmiary-color (for the hover).
 *   - --border-color
 *
 * Usage example:
 *    <popup-menu trigger="#trigger">
 *       <popup-menu-item id="item1" text="Menu Item 1" icon="log-out"></popup-menu-item>
 *    </popup-menu>
 *
 * @class PopupMenu
 * @extends HTMLElement
 */
export class PopupMenu extends HTMLElement {
    _trigger: string;
    isVisible: boolean;
    connectedCallback(): void;
    disconnectedCallback(): void;
    /**
    * Toggles the visibility of the popup menu
    * @param {Event} e The click event
    * @returns {void}
    */
    toggle(e: Event): void;
    _hide(): void;
    _show(): void;
}
/**
 * Represents a popup menu item
 * @class PopupMenuItem
 * @extends HTMLElement
 */
export class PopupMenuItem extends HTMLElement {
    connectedCallback(): void;
    _render(): void;
}
export function showPopup(el: string): void;
export function hidePopup(el: string): void;
