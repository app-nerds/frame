import { parseDateTime, formatDateTime, DateFormats } from "./date-time-service.js";

/**
 * date-time-picker is a custom HTML element that allows the user to select a date and time.
 * It supports custom date formats.
 * @class DateTimePicker
 * @extends HTMLElement
 */
export class DateTimePicker extends HTMLElement {
	constructor() {
		super();

		this._daysOfTheWeek = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
		this._months = [
			"January",
			"February",
			"March",
			"April",
			"May",
			"June",
			"July",
			"August",
			"September",
			"October",
			"November",
			"December",
		];

		this.name = this.getAttribute("name") || "dateTime";

		/* Get the date from attributes. If one isn't passed in, use now, but zero out the time. */
		this.date = (this.getAttribute("date")) ? parseDateTime(this.getAttribute("date")) : "";

		this.dateFormat = this.getAttribute("date-format") || DateFormats.IsoWithTimezone;
		this.showTimeSelector = this.dateFormat === DateFormats.IsoWithTimezone || this.dateFormat === DateFormats.IsoWithoutTimezone ||
			this.dateFormat === DateFormats.UsDateTimeWithSeconds || this.dateFormat === DateFormats.UsDateTimeWithoutSeconds ||
			this.dateFormat === DateFormats.InternationalWithSeconds || this.dateFormat === DateFormats.International;
		this.twentyFourHourTime = this.dateFormat === DateFormats.IsoWithTimezone || this.dateFormat === DateFormats.IsoWithoutTimezone ||
			this.dateFormat === DateFormats.InternationalWithSeconds || this.dateFormat === DateFormats.International;
		this.timeIncrement = this.getAttribute("time-increment") || "hour"; // valid values are hour,30minute,15minute,10minute,5minute,1minute
		this.today = new Date();
		this.inputEl = null;
		this.popupEl = null;
		this.headerEl = null;
		this.bodyEl = null;
		this.day = 0;
		this.timeSelectorEl = null;
		this.selectedTimeIndex = 0;
		this.yearBlockStart = this._getYear() - 5;
	}

	connectedCallback() {
		this.setAttribute("name", `${this.name}-datepicker`)
		this.setAttribute("aria-label", "Date Picker");

		this.inputEl = this._createInputEl();
		let formatP = this._createInputLabel();

		if (this.date === "") {
			this.date = new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate(), 0, 0, 0);
			this.yearBlockStart = this._getYear() - 5;
		}

		this.popupEl = this._createPopupEl();
		this._drawHeaderEl();
		this._drawCalendarBody();

		this.insertAdjacentElement("beforeend", this.inputEl);
		this.insertAdjacentElement("beforeend", formatP);
		this.insertAdjacentElement("beforeend", this.popupEl);

		document.addEventListener("click", e => {
			if (!e.target.contains(this.popupEl) && !e.target.contains(this.inputEl) && !this.popupEl.contains(e.target)) {
				this.hide();
			}
		})
	}

	/****************************************************
	 * PUBLIC METHODS
	 ****************************************************/

	/**
	 * clear clears the date picker value.
	 * @returns {void}
	 */
	clear() {
		this.inputEl.value = "";
	}

	/**
	 * getDate returns the currently selected date.
	 * @returns {string|Date}
	 */
	getDate() {
		return this.date;
	}

	/**
	 * Moves the calendar forward or backward one month. A positive number moves forward, a negative number moves backward.
	 * @param {number} direction
	 */
	moveMonth(e, direction) {
		e.preventDefault();
		e.stopPropagation();

		let newDate = new Date(this.date);
		newDate.setMonth(newDate.getMonth() + direction);

		this.date = newDate;
		this.popupEl.querySelector("header button:nth-child(2)").innerText = this._getMonthName();
		this.popupEl.querySelector("header button:nth-child(3)").innerText = this._getYear().toString();
		this.popupEl.querySelector(".calendar-body").remove();
		this._drawCalendarBody();
	}

	/**
	 * setDate sets the date picker value.
	 * @param {Date} dt
	 * @returns {void}
	 */
	setDate(dt) {
		this.date = dt;
		this.day = dt.getDate();
		this._setInputDate();
	}

	/**
	 * toggle shows or hides the calendar.
	 * @returns {void}
	 */
	toggle() {
		this.popupEl.classList.toggle("calendar-hidden");
		this.inputEl.focus();
	}

	hide() {
		this.popupEl.classList.add("calendar-hidden");
	}

	/****************************************************
	 * PRIVATE METHODS
	 ****************************************************/

	_drawHeaderEl() {
		let previousMonthEl = this._createPreviousMonthButton();
		let nextMonthEl = this._createNextMonthButton();
		let currentMonthEl = this._createCurrentMonthButton();
		let currentYearEl = this._createCurrentYearButton();

		this.headerEl.innerHTML = "";
		this.headerEl.insertAdjacentElement("beforeend", previousMonthEl);
		this.headerEl.insertAdjacentElement("beforeend", currentMonthEl);
		this.headerEl.insertAdjacentElement("beforeend", currentYearEl);
		this.headerEl.insertAdjacentElement("beforeend", nextMonthEl);
	}

	_drawCalendarBody() {
		let bodyDiv = document.createElement("div");
		let weekDiv = this._createCalendarBodyWeekDiv();

		let firstDate = this._getFirstDayOfMonth();
		let firstDayOfWeek = firstDate.getDay();
		let lastDate = this._getLastDayOfMonth();
		let lastDay = lastDate.getDate();
		let started = false;

		bodyDiv.classList.add("calendar-body");

		for (let dayIndex = 0; dayIndex < lastDay + firstDayOfWeek; dayIndex++) {
			/*
			 * Basically we want to not render day numbers until we hit the
			 * first day of the month on the correct day of the week.
			 */
			if (!started) {
				if (dayIndex === firstDayOfWeek) {
					started = true;
				}
			}

			let dayDiv = this._createCalendarBodyDayDiv(started, dayIndex, firstDayOfWeek);
			weekDiv.insertAdjacentElement("beforeend", dayDiv);

			/*
			 * Create a new week div every 7 days.
			 */
			if (!((dayIndex + 1) % 7)) {
				bodyDiv.insertAdjacentElement("beforeend", weekDiv);
				weekDiv = this._createCalendarBodyWeekDiv();
			}
		}

		if (weekDiv.innerHTML !== "") {
			bodyDiv.insertAdjacentElement("beforeend", weekDiv);
		}

		if (this.showTimeSelector) {
			this._createTimeSelector();
			bodyDiv.insertAdjacentElement("beforeend", this.timeSelectorEl);

			let okButton = this._createOkButton();
			bodyDiv.insertAdjacentElement("beforeend", okButton);
		}

		this._replaceBodyEl(bodyDiv);
	}

	_drawMonthListBody() {
		const body = document.createElement("div");
		body.classList.add("month-list-body");

		for (let monthIndex = 0; monthIndex < this._months.length; monthIndex++) {
			let month = this._createMonthButton(monthIndex);
			body.insertAdjacentElement("beforeend", month);
		}

		this._replaceBodyEl(body);
	}

	_drawYearListBody() {
		const body = document.createElement("div");
		body.classList.add("year-list-body");

		const yearList = document.createElement("div");
		yearList.classList.add("year-list");

		const yearUp = this._createYearUpButton();
		const yearDown = this._createYearDownButton();

		body.insertAdjacentElement("beforeend", yearUp);

		for (let yearIndex = this.yearBlockStart; yearIndex < this.yearBlockStart + 10; yearIndex++) {
			let yearButton = this._createYearButton(yearIndex);
			yearList.insertAdjacentElement("beforeend", yearButton);
		}

		body.insertAdjacentElement("beforeend", yearList);
		body.insertAdjacentElement("beforeend", yearDown);

		this._replaceBodyEl(body);
	}

	_getMonth() { return new Date(this.date).getMonth(); }
	_getMonthName() { return this._months[this._getMonth()]; }
	_getYear() { return new Date(this.date).getFullYear(); }
	_getDay() { return new Date(this.date).getDate(); }
	_getFirstDayOfMonth() {
		let result = new Date(this._getYear(), this._getMonth(), 1);
		return result;
	}
	_getLastDayOfMonth() {
		let result = new Date(this._getYear(), this._getMonth() + 1, 0);
		return result;
	}
	_getHour() { return new Date(this.date).getHours(); }
	_getMinute() { return new Date(this.date).getMinutes(); }
	_getSecond() { return new Date(this.date).getSeconds(); }

	/**
	 * @param {number} day
	 */
	_onCalendarDayClick(e, day) {
		e.preventDefault();
		e.stopPropagation();

		this.day = day;
		this._setInputDate();

		if (!this.showTimeSelector) {
			this.toggle();
			this.inputEl.focus();
		} else {
			this._createTimeSelectorOptions();
		}
	}

	_onHeaderMonthClick(e) {
		e.preventDefault();
		e.stopPropagation();
		this._drawMonthListBody();
	}

	/**
	 * @returns {void}
	 */
	_onHeaderYearClick(e) {
		e.preventDefault();
		e.stopPropagation();
		this._drawYearListBody();
	}

	/**
	 * @param {number} monthIndex
	 */
	_onMonthClick(e, monthIndex) {
		e.preventDefault();
		e.stopPropagation();

		this.date = new Date(this._getYear(), monthIndex, 1);
		this._setInputDate();
		this._drawHeaderEl();
		this._drawCalendarBody();
	}

	/**
	 * @param {Event & { target: HTMLSelectElement }} e
	 */
	_onTimeChange(e) {
		let selected = e.target.value;
		this.selectedTimeIndex = e.target.selectedIndex;
		this.date = selected;
		this._setInputDate();
	}

	/**
	 * @param {number} year
	 */
	_onYearClick(e, year) {
		e.preventDefault();
		e.stopPropagation();

		this.date = new Date(year, this._getMonth(), 1);
		this._setInputDate();
		this._drawHeaderEl();
		this._drawCalendarBody();
	}

	_onYearDownClick(e) {
		e.preventDefault();
		e.stopPropagation();

		this.yearBlockStart += 10;
		this._drawYearListBody();
	}

	_onYearUpClick(e) {
		e.preventDefault();
		e.stopPropagation();

		this.yearBlockStart -= 10;
		this._drawYearListBody();
	}

	/**
	 * @param {HTMLDivElement} newBody
	 */
	_replaceBodyEl(newBody) {
		this.bodyEl.innerHTML = "";
		this.bodyEl.insertAdjacentElement("beforeend", newBody);
	}

	_setInputDate() {
		let selected = new Date(this._getYear(), this._getMonth(), this.day, this._getHour(), this._getMinute(), this._getSecond());

		this.inputEl.value = formatDateTime(selected, this.dateFormat);
		this.dispatchEvent(new CustomEvent("change", { detail: { value: selected } }));
	}


	/**********************************************************************
	 * Methods to return invididual elements
	 *********************************************************************/

	/**
	 * @param {boolean} started
	 * @param {number} dayIndex
	 * @param {number} firstDayOfWeek
	 * @returns {HTMLDivElement}
	 */
	_createCalendarBodyDayDiv(started, dayIndex, firstDayOfWeek) {
		const el = Object.assign(document.createElement("div"), { className: "day" });

		if (started) {
			let d = dayIndex - firstDayOfWeek + 1;

			const a = Object.assign(document.createElement("button"), { innerText: `${d}`, type: "button" });
			a.addEventListener("click", e => this._onCalendarDayClick.call(this, e, d));

			let thisDay = new Date(this._getYear(), this._getMonth(), d);
			if (thisDay === this.today) {
				a.classList.add("today");
			}

			el.insertAdjacentElement("beforeend", a);
		} else {
			el.classList.add("disabled");
		}

		return el;
	}

	_createCalendarBodyWeekDiv() {
		return Object.assign(document.createElement("div"), { className: "week" });
	}

	_createCurrentMonthButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: this._getMonthName(),
			type: "button",
		});

		el.addEventListener("click", e => this._onHeaderMonthClick.call(this, e));
		return el;
	}

	_createCurrentYearButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: this._getYear().toString(),
			type: "button",
		});

		el.addEventListener("click", e => this._onHeaderYearClick.call(this, e));
		return el;
	}

	_createInputEl() {
		const el = Object.assign(document.createElement("input"), {
			type: "text",
			name: this.name,
			"aria-describedby": `${this.name}-format`,
			value: this.date instanceof Date ? formatDateTime(this.date, this.dateFormat) : "",
		});

		el.addEventListener("click", () => {
			if (this.date === "") {
				this.date = new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate(), 0, 0, 0);
				this._drawHeaderEl();
				this._drawCalendarBody();

			}

			this.toggle();
		});

		return el;
	}

	_createInputLabel() {
		const el = Object.assign(document.createElement("p"), { innerText: `(${this.dateFormat})`, id: `${this.name}-format` });
		return el;
	}

	/**
	 * @param {number} monthIndex
	 * @returns {HTMLAnchorElement}
	 */
	_createMonthButton(monthIndex) {
		const month = Object.assign(document.createElement("button"), {
			innerText: this._months[monthIndex],
			type: "button",
		});

		month.addEventListener("click", e => this._onMonthClick.call(this, e, monthIndex));
		return month;
	}

	_createNextMonthButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: `<i class="icon--mdi icon--mdi--arrow-right"></i>`,
			type: "button",
		});

		el.addEventListener("click", e => this.moveMonth.call(this, e, 1));
		return el;
	}

	_createOkButton() {
		const el = Object.assign(document.createElement("button"), {
			innerText: "OK",
			type: "button",
			className: "ok",
		});

		el.addEventListener("click", this.toggle.bind(this));
		return el;
	}

	_createPopupEl() {
		this.headerEl = document.createElement("header");
		this.bodyEl = document.createElement("section");

		const el = Object.assign(document.createElement("div"), {
			className: "date-time-picker-popup calendar-hidden",
			role: "dialog",
			"aria-modal": "true",
			"aria-label": `Choose Date`,
		});

		el.insertAdjacentElement("beforeend", this.headerEl);
		el.insertAdjacentElement("beforeend", this.bodyEl);

		return el;
	}

	_createPreviousMonthButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: `<i class="icon--mdi icon--mdi--arrow-left"></i>`,
			type: "button",
		});

		el.addEventListener("click", e => this.moveMonth.call(this, e, -1));
		return el;
	}

	_createTimeSelector() {
		this.timeSelectorEl = document.createElement("select");
		this._createTimeSelectorOptions();
		this.timeSelectorEl.addEventListener("change", this._onTimeChange.bind(this));
	}

	_createTimeSelectorOptions() {
		this.timeSelectorEl.innerHTML = "";

		let increment = 1;

		if (this.timeIncrement === "5minute") {
			increment = 5;
		}

		if (this.timeIncrement === "10minute") {
			increment = 10;
		}

		if (this.timeIncrement === "15minute") {
			increment = 15;
		}

		if (this.timeIncrement === "30minute") {
			increment = 30;
		}

		if (this.timeIncrement === "hour") {
			increment = 60;
		}

		let start = new Date(this._getYear(), this._getMonth(), this._getDay(), 0, 0, 0);
		let index = 0;

		for (let i = 0; i < 1440; i += increment) {
			let option = document.createElement("option");
			option.value = formatDateTime(start, this.dateFormat);

			let selected = index === this.selectedTimeIndex ? true : false;

			if (selected) {
				option.setAttribute("selected", "selected");
			}

			if (this.twentyFourHourTime) {
				option.innerText = `${start.getHours().toString().padStart(2, "0")}:${start.getMinutes().toString().padStart(2, "0")}`;
			} else {
				let hours = start.getHours();
				let ampm = "AM";

				if (hours > 12) {
					hours -= 12;
					ampm = "PM";
				}

				if (hours === 12) {
					ampm = "PM";
				}

				if (hours === 0) {
					hours = 12;
				}

				option.innerText = `${hours.toString().padStart(2, "0")}:${start.getMinutes().toString().padStart(2, "0")} ${ampm}`;
			}

			this.timeSelectorEl.insertAdjacentElement("beforeend", option);
			start = new Date(start.getTime() + increment * 60000);
			index++;
		}
	}

	/**
	 * @param {number} year
	 * @returns {HTMLAnchorElement}
	 */
	_createYearButton(year) {
		const el = Object.assign(document.createElement("button"), {
			innerText: year.toString(),
			type: "button",
		});

		el.addEventListener("click", e => this._onYearClick.call(this, e, year));
		return el;
	}

	_createYearDownButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: `<i class="icon--mdi icon--mdi--arrow-down"></i>`,
			type: "button",
		});

		el.addEventListener("click", e => this._onYearDownClick.call(this, e));
		return el;
	}

	_createYearUpButton() {
		const el = Object.assign(document.createElement("button"), {
			innerHTML: `<i class="icon--mdi icon--mdi--arrow-up"></i>`,
			type: "button",
		});

		el.addEventListener("click", e => this._onYearUpClick.call(this, e));
		return el;
	}

}

customElements.define("date-time-picker", DateTimePicker);
