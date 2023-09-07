export type AlertOptions = object & {
    position: AlertPosition;
    duration: number;
    closable: boolean;
    focusable: boolean;
};
export type ShimOptions = object & {
    closeOnClick: boolean;
    onShimClick: Function;
};
export type ConfirmOptions = object & {
    callback: Function;
};
export type GraphQLOptions = object & {
    http: typeof fetcher;
    tokenGetterFunction: Function;
    expiredTokenCallback: Function;
    spinner: object;
    navigateTo: Function;
};
/**
 * Callback used to validate the values entered into a Prompter
 */
export type ValidatorFunc = (promptValues: any) => any;
export type Route = any;
export type Application = object & {
    routes: Array<Route>;
    targetElement: HTMLElement;
    router: Router;
    afterRoute: Function;
    beforeRoute: Function;
    injectParams: Function;
    onRenderComplete: Function;
    go: Function;
};
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
         */ success(message: string, callback: Function): void;
    /**
         * info displays an info alert. Use this for neutral messages.
         * @param {string} message
         * @param {function} callback
         * @returns {void}
         */ info(message: string, callback: Function): void;
    /**
         * warn displays a warning alert. Use this to warn users of something.
         * @param {string} message
         * @param {function} callback
         * @returns {void}
         */ warn(message: string, callback: Function): void;
    /**
         * error displays an error alert. Use this to warn users of something bad.
         * @param {string} message
         * @param {function} callback
         * @returns {void}
         */ error(message: string, callback: Function): void;
    /**
         * @param {string} message
         * @param {string} type
         * @param {function} callback
         * @returns {void}
         */ show(message: string, type: string, callback: Function): void;
    _setContent(card: any): void;
    /**
         * @param {AlertCard} card
         * @returns {void}
         */ _setIntroAnimation(card: AlertCard): void;
    /**
         * @param {AlertCard} card
         * @returns {void}
         */ _bindEvents(card: AlertCard): void;
    /**
         * @param {AlertCard} card
         * @returns {void}
         */ _autoDestroy(card: AlertCard, callback: any): void;
    /**
         * @param {AlertCard} card
         * @returns {void}
         */ _destroy(card: AlertCard, callback: any): void;
    _setup(): void;
}
/**
 * BaseView is the base class for all views in the application. It provides
 * a common set of functionality that all views can use. Your view JavaScript
 * components should extend this class and register themselves as custom elements.
 * @class BaseView
 * @extends {HTMLElement}
 * @property {string} title The title of the view. This is used to set the document title.
 * @property {object} params The parameters passed to the view.
 * @property {object} state The state of the view.
 */ export class BaseView extends HTMLElement {
    constructor(params: any, _onRenderComplete: any);
    _title: string;
    params: any;
    _state: {};
    _onRenderComplete: any;
    router: any;
    connectedCallback(): Promise<void>;
    disconnectedCallback(): void;
    /**
         * This method is called before the view is rendered. Override this method
         * to perform any actions before the view is rendered.
         * @returns {Promise<void>}
         */ beforeRender(): Promise<void>;
    /**
         * This method is called after the view is rendered. Override this method
         * to perform any actions after the view is rendered.
         * @returns {Promise<void>}
         */ afterRender(): Promise<void>;
    /**
         * This method is called when the view is unloaded. Override this method
         * to perform any actions when the view is unloaded.
         * @returns {Promise<void>}
         */ onUnload(): Promise<void>;
    /**
         * This method is called when the view is navigated to. Override this method
         * render your page contents.
         * @returns {Promise<void>}
         */ render(): Promise<void>;
    /**
         * Get the title for the current view.
         * @returns {string}
         */ get title(): string;
    /**
         * Get the HTML for the current view.
         * @returns {string}
         */ get html(): string;
    /**
         * Set the state for the current view.
         * @param {object} newState The new state for the view.
         * @returns {void}
         */ set state(arg: any);
    /**
         * Get the state for the current view.
         * @returns {object}
         */ get state(): any;
    /**
         * Get the value of a query parameter.
         * @param {string} paramName The name of the query parameter.
         * @returns {string}
         */ getQueryParam(paramName: string): string;
    /**
         * Navigate to a new URL.
         * @param {string} url The URL to navigate to.
         * @param {object} queryParams Query parameters to add to the URL.
         * @param {object} state The state to pass to the new view.
         * @returns {void}
         */ navigateTo(url: string, queryParams?: object, state?: object): void;
    _setDocumentTitle(): void;
}
export class ColorPicker extends HTMLElement {
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
 */ export class Confirmer {
    /**
         * confirm displays a confirmation dialog. It shows a message and a Close button.
         * @param {string} message
         * @param {function} callback
         * @returns {void}
         */ confirm(message: string, callback: Function): void;
    /**
         * yesNo displays a confirmation dialog. It shows a message and Yes and No buttons.
         * @param {string} message
         * @returns {Promise<boolean>}
         */ yesNo(message: string): Promise<boolean>;
    /**
         * show displays a confirmation dialog. This is a raw function that is normally
         * used by the yesNo and confirm functions.
         * @param {string} type
         * @param {string} message
         * @param {function} callback
         * @returns {void}
         */ show(type: string, message: string, callback: Function): void;
    _close(container: any, callback: any, callbackValue: any): void;
    _addButtons(container: any, type: any, shim: any, callback: any): void;
}
export namespace DateFormats {
    let IsoWithTimezone: string;
    let IsoWithoutTimezone: string;
    let InternationalWithSeconds: string;
    let International: string;
    let UsDateTimeWithSeconds: string;
    let UsDateTimeWithoutSeconds: string;
    let IsoDate: string;
    let UsDate: string;
}
/**
 * date-time-picker is a custom HTML element that allows the user to select a date and time.
 * It supports custom date formats.
 * @class DateTimePicker
 * @extends HTMLElement
 */ export class DateTimePicker extends HTMLElement {
    _daysOfTheWeek: string[];
    _months: string[];
    name: string;
    date: string | Date;
    dateFormat: string;
    showTimeSelector: boolean;
    twentyFourHourTime: boolean;
    timeIncrement: string;
    today: Date;
    inputEl: any;
    popupEl: any;
    headerEl: HTMLElement;
    bodyEl: HTMLElement;
    day: number;
    timeSelectorEl: HTMLSelectElement;
    selectedTimeIndex: number;
    yearBlockStart: number;
    connectedCallback(): void;
    /****************************************************
         * PUBLIC METHODS
         ****************************************************/
    /**
         * clear clears the date picker value.
         * @returns {void}
         */ clear(): void;
    /**
         * getDate returns the currently selected date.
         * @returns {string|Date}
         */ getDate(): string | Date;
    /**
         * Moves the calendar forward or backward one month. A positive number moves forward, a negative number moves backward.
         * @param {number} direction
         */ moveMonth(e: any, direction: number): void;
    /**
         * setDate sets the date picker value.
         * @param {Date} dt
         * @returns {void}
         */ setDate(dt: Date): void;
    /**
         * toggle shows or hides the calendar.
         * @returns {void}
         */ toggle(): void;
    hide(): void;
    /****************************************************
         * PRIVATE METHODS
         ****************************************************/ _drawHeaderEl(): void;
    _drawCalendarBody(): void;
    _drawMonthListBody(): void;
    _drawYearListBody(): void;
    _getMonth(): number;
    _getMonthName(): string;
    _getYear(): number;
    _getDay(): number;
    _getFirstDayOfMonth(): Date;
    _getLastDayOfMonth(): Date;
    _getHour(): number;
    _getMinute(): number;
    _getSecond(): number;
    /**
         * @param {number} day
         */ _onCalendarDayClick(e: any, day: number): void;
    _onHeaderMonthClick(e: any): void;
    /**
         * @returns {void}
         */ _onHeaderYearClick(e: any): void;
    /**
         * @param {number} monthIndex
         */ _onMonthClick(e: any, monthIndex: number): void;
    /**
         * @param {Event & { target: HTMLSelectElement }} e
         */ _onTimeChange(e: Event & {
        target: HTMLSelectElement;
    }): void;
    /**
         * @param {number} year
         */ _onYearClick(e: any, year: number): void;
    _onYearDownClick(e: any): void;
    _onYearUpClick(e: any): void;
    /**
         * @param {HTMLDivElement} newBody
         */ _replaceBodyEl(newBody: HTMLDivElement): void;
    _setInputDate(): void;
    /**********************************************************************
         * Methods to return invididual elements
         *********************************************************************/
    /**
         * @param {boolean} started
         * @param {number} dayIndex
         * @param {number} firstDayOfWeek
         * @returns {HTMLDivElement}
         */ _createCalendarBodyDayDiv(started: boolean, dayIndex: number, firstDayOfWeek: number): HTMLDivElement;
    _createCalendarBodyWeekDiv(): any;
    _createCurrentMonthButton(): any;
    _createCurrentYearButton(): any;
    _createInputEl(): any;
    _createInputLabel(): any;
    /**
         * @param {number} monthIndex
         * @returns {HTMLAnchorElement}
         */ _createMonthButton(monthIndex: number): HTMLAnchorElement;
    _createNextMonthButton(): any;
    _createOkButton(): any;
    _createPopupEl(): any;
    _createPreviousMonthButton(): any;
    _createTimeSelector(): void;
    _createTimeSelectorOptions(): void;
    /**
         * @param {number} year
         * @returns {HTMLAnchorElement}
         */ _createYearButton(year: number): HTMLAnchorElement;
    _createYearDownButton(): any;
    _createYearUpButton(): any;
}
export const ErrTokenExpired: "token expired";
export class GoogleLoginForm extends HTMLElement {
    loginPath: string;
    createAccountPath: string;
    signInButtonURL: string;
    connectedCallback(): void;
}
/** @typedef { object & { http: fetcher, tokenGetterFunction: function, expiredTokenCallback: function, spinner: object, navigateTo: function } } GraphQLOptions */
/**
 * This class is a wrapper around the fetcher function
 * that makes it easy to execute GraphQL queries and mutations.
 * It also handles expired tokens.
 * @class GraphQL
 * @param {string} queryURL The URL to the GraphQL API
 * @param {GraphQLOptions} options Options for the GraphQL class
 */ export class GraphQL {
    constructor(queryURL: any, options?: {
        http: typeof fetcher;
        tokenGetterFunction: any;
        expiredTokenCallback: any;
        spinner: any;
        navigateTo: any;
    });
    queryURL: any;
    http: typeof fetcher;
    tokenGetterFunction: any;
    expiredTokenCallback: any;
    spinner: any;
    navigateTo: any;
    /**
         * Executes a query against a GraphQL API
         * @param query string A graphql query. Omit the "query {}" portion.
         * @returns {Promise<object>} A promise that resolves to the fetch response
         */ query(query: any): Promise<object>;
    /**
         * Executes a mutation against a GraphQL API
         * @param query string A graphql mutation. Omit the "mutation {}" portion
         * @returns {Promise<object>} A promise that resolves to the fetch response
         */ mutation(query: any): Promise<object>;
}
export class MemberLoginBar extends HTMLElement {
    static get observedAttributes(): string[];
    memberService: MemberService;
    loginPath: string;
    attributedChangedCallback(name: any, _: any, newValue: any): void;
    connectedCallback(): Promise<void>;
    /*******************************************************************************
       * Event methods
       ******************************************************************************/
    /*******************************************************************************
       * UI elements
       ******************************************************************************/ createContainerEl(): HTMLDivElement;
    createAvatarEl(container: any, member: any): void;
    createTextEl(container: any, member: any): void;
    createPopupMenu(container: any, member: any): void;
    /*******************************************************************************
       * Private methods
       ******************************************************************************/ onMyAccountClick(): void;
    onLogOutClick(): void;
}
export class MemberService {
    constructor(spinnerEl: any);
    spinnerEl: any;
    getCurrentMember(): Promise<any>;
}
export class MessageBar extends HTMLElement {
    messageType: string;
    message: string;
    containerEl: HTMLDivElement;
    connectedCallback(): void;
    _createContainerEl(): HTMLDivElement;
    _createCloseButtonEl(): HTMLSpanElement;
    _createTextEl(): HTMLParagraphElement;
}
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
        */ toggle(e: Event): void;
    _hide(): void;
    _show(): void;
}
/**
 * Represents a popup menu item
 * @class PopupMenuItem
 * @extends HTMLElement
 */ export class PopupMenuItem extends HTMLElement {
    connectedCallback(): void;
    _render(): void;
}
/**
 * Converts a classic JS object to a Map
 * Copyright © 2023 App Nerds LLC
 * @param o object The object to convert
 * @returns {Map} A Map
 */
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
         */ addValidatorFunc(f: ValidatorFunc): void;
    _onCancelClick(): void;
    _onConfirmClick(): any;
    _renderWindow(): void;
    _clearAllInputs(): void;
}
/**
 * SessionService is a service used to manage session data. It makes use
 * of the browser's sessionStorage object. This class provides
 * only static methods.
 * @class SessionService
 */ export class SessionService {
    /**
         * Clears the member from session storage.
         * @static
         * @memberof SessionService
         * @returns {void}
         */
    static clearMember(): void;
    /**
         * Clears the token from session storage.
         * @static
         * @memberof SessionService
         * @returns {void}
         */ static clearToken(): void;
    /**
         * Gets the member from session storage.
         * @static
         * @memberof SessionService
         * @returns {object} The member object
         */ static getMember(): object;
    /**
         * Gets the token from session storage.
         * @static
         * @memberof SessionService
         * @returns {object} The token object
         */ static getToken(): object;
    /**
         * Determines if the session has a member.
         * @static
         * @memberof SessionService
         * @returns {boolean} True if the session has a member, otherwise false.
         */ static hasMember(): boolean;
    /**
         * Determines if the session has a token.
         * @static
         * @memberof SessionService
         * @returns {boolean} True if the session has a token, otherwise false.
         */ static hasToken(): boolean;
    /**
         * Navigates to the specified path if the token has expired.
         * This is determined by examing the error message. If the
         * error message is "token expired", the user is navigated
         * to the specified path.
         * @static
         * @param {object} e The error object
         * @param {string} path The path to navigate to
         * @param {function} navigateTo The function to use to navigate
         * @memberof SessionService
         */ static navigateOnTokenExpired(e: object, path: string, navigateTo: Function): void;
    /**
         * Sets the member in session storage.
         * @static
         * @param {object} member The member object
         * @memberof SessionService
         * @returns {void}
         */ static setMember(member: object): void;
    /**
         * Sets the token in session storage.
         * @static
         * @param {object} token The token object
         * @memberof SessionService
         * @returns {void}
         */ static setToken(token: object): void;
    /**
         * This is a function that can be used as a callback for the fetcher
         * class. It will check the response for a 401 status code. If it
         * finds one, it will clear the token and navigate to the specified
         * path.
         * @static
         * @param {object} httpResponse The HTTP response object
         * @param {string} path The path to navigate to
         * @param {function} navigateTo The function to use to navigate
         * @memberof SessionService
         * @returns {boolean} True if the token is valid, otherwise false.
         */ static tokenExpireFunc(httpResponse: object, path: string, navigateTo: Function): boolean;
}
/** @typedef {object & { closeOnClick: boolean, onShimClick: function }} ShimOptions */
/**
 * Shim displays a full screen shim to cover elements.
 * @param {ShimOptions} options
 */ export class Shim {
    constructor(closeOnClick: boolean, onShimClick: any);
    closeOnClick: boolean;
    onShimClick: any;
    /**
         * show displays the shim
         * @returns {void}
         */ show(): void;
    shim: Element | HTMLDivElement;
    /**
         * hide removes the shim
         * @returns {void}
         */ hide(callback: any): void;
    _destroy(): void;
}
export class Spinner {
    spinnerEl: HTMLDivElement;
    hide(): void;
    show(): void;
}
/** @typedef {import("./router.js").Route} Route */
/** @typedef {object & {routes: Array<Route>, targetElement: HTMLElement, router: Router, afterRoute: function, beforeRoute: function, injectParams: function, onRenderComplete: function, go: function }} Application */
/**
 * Creates a new single-page application.
 * @param {HTMLElement} targetElement The element to render the SPA into.
 * @param {Array<Route>} routes The routes to use for the SPA.
 * @param {BaseView} pageNotFoundView The view to use when a route is not found.
 * @returns {Application}
 */ export function application(targetElement: HTMLElement, routes: Array<Route>, pageNotFoundView?: BaseView): Application;
/**
 * Debounces a function call. This is useful for things like
 * search boxes where you don't want to make a call to the
 * server for every keystroke.
 * Copyright © 2023 App Nerds LLC
 * @param {function} fn The function to debounce
 * @param {number} delay The delay in milliseconds. Default is 400
 * @returns {function} The debounced function
 */ export function debounce(fn: Function, delay?: number): Function;
/**
 * A wrapper around fetch that will show a spinner
 * while the request is being made. This is useful for
 * long running requests.
 * @param {string} url The URL to fetch
 * @param {object} options The fetch options
 * @param {object} spinner The spinner element to show
 * @param {number} msBeforeShowSpinner The number of milliseconds to wait before showing the spinner. Default is 1000
 * @returns {Promise<object>} A promise that resolves to the fetch response
 */ export function fetcher(url: string, options: object, spinner: object, msBeforeShowSpinner?: number): Promise<object>;
/**
 * formatDateTime formats a Date object into a string using the specified format.
 * @param {string|number|Date} dt
 * @param {string} format
 * @returns {string}
 */ export function formatDateTime(dt: string | number | Date, format: string): string;
/**
 * Shows a popup menu
 * @param {string} el The query selector for the popup menu
 * @returns {void}
 */ export function hidePopup(el: string): void;
/**
 * Debounces a function call. This is useful for things like
 * search boxes where you don't want to make a call to the
 * server for every keystroke.
 * Copyright © 2023 App Nerds LLC
 * @param {function} fn The function to debounce
 * @param {number} delay The delay in milliseconds. Default is 400
 * @returns {function} The debounced function
 */ export function objectToMap(o?: {}): Function;
/**
 * parseDateTime parses a date/time string into a Date object.
 * @param {string|number|Date} dt
 * @returns {Date}
 */
export function parseDateTime(dt: string | number | Date): Date;
/**
 * Shows a popup menu
 * @param {string} el The query selector for the popup menu
 * @returns {void}
 */ export function showPopup(el: string): void;
declare class Router {
    /**
         * Creates a new instance of Router.
         * @param {string} targetEl The element to render the SPA into.
         * @param {Array<Route>} routes The routes to use for the SPA.
         * @param {BaseView} pageNotFoundView The view to use when a route is not found.
         */
    constructor(targetEl: string, routes: Array<Route>, pageNotFoundView?: BaseView);
    targetEl: string;
    routes: any[];
    pageNotFoundView: BaseView;
    beforeRoute: any;
    afterRoute: any;
    injectParams: any;
    onRenderComplete: any;
    /**
         * Retrieves a query parameter from the URL by name.
         * @param {string} paramName The name of the query parameter to retrieve.
         * @returns {string}
         */ getQueryParam(paramName: string): string;
    /**
         * Navigates to a URL.
         * @param {string} url The URL to navigate to.
         * @param {object} queryParams Query parameters to add to the URL.
         * @param {object} state The state to pass to the new view.
         * @returns {void}
         */ navigateTo(url: string, queryParams?: object, state?: object): void;
    _pathToRegex(path: any): RegExp;
    _getParams(match: any): {};
    _route(e: any): Promise<void>;
}
export {};
