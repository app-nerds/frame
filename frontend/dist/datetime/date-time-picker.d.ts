/**
 * date-time-picker is a custom HTML element that allows the user to select a date and time.
 * It supports custom date formats.
 * @class DateTimePicker
 * @extends HTMLElement
 */
export class DateTimePicker extends HTMLElement {
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
     */
    clear(): void;
    /**
     * getDate returns the currently selected date.
     * @returns {string|Date}
     */
    getDate(): string | Date;
    /**
     * Moves the calendar forward or backward one month. A positive number moves forward, a negative number moves backward.
     * @param {number} direction
     */
    moveMonth(e: any, direction: number): void;
    /**
     * setDate sets the date picker value.
     * @param {Date} dt
     * @returns {void}
     */
    setDate(dt: Date): void;
    /**
     * toggle shows or hides the calendar.
     * @returns {void}
     */
    toggle(): void;
    hide(): void;
    /****************************************************
     * PRIVATE METHODS
     ****************************************************/
    _drawHeaderEl(): void;
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
     */
    _onCalendarDayClick(e: any, day: number): void;
    _onHeaderMonthClick(e: any): void;
    /**
     * @returns {void}
     */
    _onHeaderYearClick(e: any): void;
    /**
     * @param {number} monthIndex
     */
    _onMonthClick(e: any, monthIndex: number): void;
    /**
     * @param {Event & { target: HTMLSelectElement }} e
     */
    _onTimeChange(e: Event & {
        target: HTMLSelectElement;
    }): void;
    /**
     * @param {number} year
     */
    _onYearClick(e: any, year: number): void;
    _onYearDownClick(e: any): void;
    _onYearUpClick(e: any): void;
    /**
     * @param {HTMLDivElement} newBody
     */
    _replaceBodyEl(newBody: HTMLDivElement): void;
    _setInputDate(): void;
    /**********************************************************************
     * Methods to return invididual elements
     *********************************************************************/
    /**
     * @param {boolean} started
     * @param {number} dayIndex
     * @param {number} firstDayOfWeek
     * @returns {HTMLDivElement}
     */
    _createCalendarBodyDayDiv(started: boolean, dayIndex: number, firstDayOfWeek: number): HTMLDivElement;
    _createCalendarBodyWeekDiv(): any;
    _createCurrentMonthButton(): any;
    _createCurrentYearButton(): any;
    _createInputEl(): any;
    _createInputLabel(): any;
    /**
     * @param {number} monthIndex
     * @returns {HTMLAnchorElement}
     */
    _createMonthButton(monthIndex: number): HTMLAnchorElement;
    _createNextMonthButton(): any;
    _createOkButton(): any;
    _createPopupEl(): any;
    _createPreviousMonthButton(): any;
    _createTimeSelector(): void;
    _createTimeSelectorOptions(): void;
    /**
     * @param {number} year
     * @returns {HTMLAnchorElement}
     */
    _createYearButton(year: number): HTMLAnchorElement;
    _createYearDownButton(): any;
    _createYearUpButton(): any;
}
