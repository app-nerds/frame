/* Copyright © 2023 App Nerds LLC v1.3.1 */
// @ts-check
const DateFormats = {
	"IsoWithTimezone": "yyyy-mm-ddthh:mm:ssz",
	"IsoWithoutTimezone": "yyyy-mm-ddthh:mm:ss",
	"InternationalWithSeconds": "yyyy-mm-dd hh:mm:ss",
	"International": "yyyy-mm-dd hh:mm",
	"UsDateTimeWithSeconds": "mm/dd/yyyy hh:mm:ss am/pm",
	"UsDateTimeWithoutSeconds": "mm/dd/yyyy hh:mm am/pm",
	"IsoDate": "yyyy-mm-dd",
	"UsDate": "mm/dd/yyyy",
};

const formatMap = {
	"yyyy-mm-ddthh:mm:ssz": format1,
	"yyyy-mm-ddthh:mm:ss": format2,
	"yyyy-mm-dd hh:mm:ss": format3,
	"yyyy-mm-dd hh:mm": format4,
	"mm/dd/yyyy hh:mm:ss am/pm": format5,
	"mm/dd/yyyy hh:mm am/pm": format6,
	"yyyy-mm-dd": format7,
	"mm/dd/yyyy": format8,
};

/**
 * parseDateTime parses a date/time string into a Date object.
 * @param {string|number|Date} dt
 * @returns {Date}
 */
function parseDateTime(dt) {
	if (typeof dt === "number") {
		return new Date(dt);
	}

	if (dt instanceof Date) {
		return dt;
	}

	if (typeof dt === "string") {
		return parseDateString(dt);
	}
}

/**
 * formatDateTime formats a Date object into a string using the specified format.
 * @param {string|number|Date} dt
 * @param {string} format
 * @returns {string}
 */
function formatDateTime(dt, format) {
	let date = parseDateTime(dt);
	let formatter = formatMap[format.toLowerCase()];

	if (!formatter) {
		throw new Error(`Invalid date format: ${format}`);
	}

	return formatter(date);
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format1(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}T${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}Z`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format2(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}T${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format3(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)} ${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format4(dt) {
	let { year, month, day, hour, minute } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)} ${zeroPad(hour)}:${zeroPad(minute)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format5(dt) {
	let { year, month, day, hour, minute, second } = breakDownDate(dt);
	let meridian = "AM";

	if (hour > 12) {
		meridian = "PM";
		hour -= 12;
	}

	if (hour === 0) {
		hour = 12;
	}

	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)} ${zeroPad(hour)}:${zeroPad(minute)}:${zeroPad(second)} ${meridian}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format6(dt) {
	let { year, month, day, hour, minute } = breakDownDate(dt);
	let meridian = "AM";

	if (hour > 12) {
		meridian = "PM";
		hour -= 12;
	}

	if (hour === 0) {
		hour = 12;
	}

	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)} ${zeroPad(hour)}:${zeroPad(minute)} ${meridian}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format7(dt) {
	let { year, month, day } = breakDownDate(dt);
	return `${zeroPad(year)}-${zeroPad(month)}-${zeroPad(day)}`;
}

/**
 * @param {Date} dt
 * @returns {string}
 */
function format8(dt) {
	let { year, month, day } = breakDownDate(dt);
	return `${zeroPad(month)}/${zeroPad(day)}/${zeroPad(year)}`;
}

/**
 * @param {Date} dt
 * @returns {object}
 */
function breakDownDate(dt) {
	let year = dt.getFullYear();
	let month = dt.getMonth() + 1;
	let day = dt.getDate();
	let hour = dt.getHours();
	let minute = dt.getMinutes();
	let second = dt.getSeconds();

	return {
		year,
		month,
		day,
		hour,
		minute,
		second,
	};
}

/**
 * @param {number} num
 * @returns {string}
 */
function zeroPad(num) {
	return num.toString().padStart(2, "0");
}

/**
 * @param {string} dt
 * @returns {Date}
 */
function parseDateString(dt) {
	const formatRegexes = [
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})t(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})z/i,                         /* YYYY-MM-DDTHH-mm-ssZ */
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})t(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})/i,                          /* YYYY-MM-DDTHH-mm-ss */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})\s+(?<hour>\d{2}):(?<minute>\d{2}):(?<second>\d{2})\s*(?<meridian>am|pm)/i, /* MM/DD/YYYY HH:mm:ss am/pm */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})\s+(?<hour>\d{2}):(?<minute>\d{2})\s*(?<meridian>am|pm)/i,                  /* MM/DD/YYYY HH:mm am/pm */
		/(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})/,                                                                            /* YYYY-MM-DD */
		/(?<month>\d{2})\/(?<day>\d{2})\/(?<year>\d{4})/,                                                                          /* MM/DD/YYYY */
	];

	for (let i = 0; i < formatRegexes.length; i++) {
		let match = dt.match(formatRegexes[i]);

		if (match) {
			let year = parseInt(match.groups.year) || new Date().getFullYear();
			let month = parseInt(match.groups.month) || new Date().getMonth();
			let day = parseInt(match.groups.day) || new Date().getDate();
			let hour = parseInt(match.groups.hour) || 0;
			let minute = parseInt(match.groups.minute) || 0;
			let second = parseInt(match.groups.second) || 0;
			let meridian = match.groups.meridian || "";

			if (meridian !== "") {
				if (meridian.toLowerCase() === "pm" && hour < 12) {
					hour += 12;
				}
			}

			return new Date(year, month - 1, day, hour, minute, second);
		}
	}

	throw new Error(`no pattern match for ${dt}`);
}

/*
 * date-time-picker is a custom HTML element that allows the user to select a date and time.
 * It supports custom date formats.
 *
 * Copyright © 2023 App Nerds LLC
 */

/**
 * date-time-picker is a custom HTML element that allows the user to select a date and time.
 * @class DateTimePicker
 * @extends HTMLElement
 */
class DateTimePicker extends HTMLElement {
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
		this.date = (this.getAttribute("date") !== "") ? parseDateTime(this.getAttribute("date")) : "";

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
		this.setAttribute("name", `${this.name}-datepicker`);
		this.setAttribute("aria-label", "Date Picker");

		this.inputEl = this._createInputEl();
		let formatP = this._createInputLabel();

		this.popupEl = this._createPopupEl();
		this._drawHeaderEl();
		this._drawCalendarBody();

		this.insertAdjacentElement("beforeend", this.inputEl);
		this.insertAdjacentElement("beforeend", formatP);
		this.insertAdjacentElement("beforeend", this.popupEl);
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
	moveMonth(direction) {
		let newDate = new Date(this.date);
		newDate.setMonth(newDate.getMonth() + direction);

		this.date = newDate;
		/** @type {HTMLAnchorElement} */ (this.popupEl.querySelector("header a:nth-child(2)")).innerText = this._getMonthName();
		/** @type {HTMLAnchorElement} */ (this.popupEl.querySelector("header a:nth-child(3)")).innerText = this._getYear().toString();
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
	 * toggleCalendar shows or hides the calendar.
	 * @returns {void}
	 */
	toggleCalendar() {
		this.popupEl.classList.toggle("calendar-hidden");
		this.inputEl.focus();
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
	_onCalendarDayClick(day) {
		this.day = day;
		this._setInputDate();

		if (!this.showTimeSelector) {
			this.toggleCalendar();
			this.inputEl.focus();
		} else {
			this._createTimeSelectorOptions();
		}
	}

	_onHeaderMonthClick() {
		this._drawMonthListBody();
	}

	/**
	 * @returns {void}
	 */
	_onHeaderYearClick() {
		this._drawYearListBody();
	}

	/**
	 * @param {number} monthIndex
	 */
	_onMonthClick(monthIndex) {
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
	_onYearClick(year) {
		this.date = new Date(year, this._getMonth(), 1);
		this._setInputDate();
		this._drawHeaderEl();
		this._drawCalendarBody();
	}

	_onYearDownClick() {
		this.yearBlockStart += 10;
		this._drawYearListBody();
	}

	_onYearUpClick() {
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
		let el = document.createElement("div");
		el.classList.add("day");

		if (started) {
			let d = dayIndex - firstDayOfWeek + 1;

			let a = document.createElement("a");
			a.href = "javascript:void(0)";
			a.innerText = `${d}`;
			a.addEventListener("click", this._onCalendarDayClick.bind(this, d));

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
		let el = document.createElement("div");
		el.classList.add("week");
		return el;
	}

	_createCurrentMonthButton() {
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = this._getMonthName();
		el.addEventListener("click", this._onHeaderMonthClick.bind(this));
		return el;
	}

	_createCurrentYearButton() {
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = this._getYear().toString();
		el.addEventListener("click", this._onHeaderYearClick.bind(this));
		return el;
	}

	_createInputEl() {
		let el = document.createElement("input");
		el.setAttribute("type", "datetime");
		el.setAttribute("name", this.name);
		el.setAttribute("aria-describedby", `${this.name}-format`);

		if (this.date instanceof Date) {
			el.value = formatDateTime(this.date, this.dateFormat);
		}

		el.addEventListener("click", () => {
			if (this.date === "") {
				this.date = new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate(), 0, 0, 0);
				this._drawHeaderEl();
				this._drawCalendarBody();

			}

			this.toggleCalendar();
		});

		return el;
	}

	_createInputLabel() {
		let el = document.createElement("p");
		el.innerText = `(${this.dateFormat})`;
		el.id = `${this.name}-format`;
		return el;
	}

	/**
	 * @param {number} monthIndex
	 * @returns {HTMLAnchorElement}
	 */
	_createMonthButton(monthIndex) {
		let month = document.createElement("a");
		month.href = "javascript:void(0)";
		month.innerText = this._months[monthIndex];
		month.addEventListener("click", this._onMonthClick.bind(this, monthIndex));
		return month;
	}

	_createNextMonthButton() {
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = `<i class="icon--mdi icon--mdi--arrow-right"></i>`;
		el.addEventListener("click", this.moveMonth.bind(this, 1));
		return el;
	}

	_createOkButton() {
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerText = "OK";
		el.classList.add("ok");
		el.addEventListener("click", this.toggleCalendar.bind(this));
		return el;
	}

	_createPopupEl() {
		let el = document.createElement("div");
		this.headerEl = document.createElement("header");
		this.bodyEl = document.createElement("section");

		el.classList.add("date-time-picker-popup", "calendar-hidden");
		el.setAttribute("role", "dialog");
		el.setAttribute("aria-modal", "true");
		el.setAttribute("aria-label", `Choose Date`);

		el.insertAdjacentElement("beforeend", this.headerEl);
		el.insertAdjacentElement("beforeend", this.bodyEl);

		return el;
	}

	_createPreviousMonthButton() {
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = `<i class="icon--mdi icon--mdi--arrow-left"></i>`;
		el.addEventListener("click", this.moveMonth.bind(this, -1));
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
		let el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerText = year.toString();
		el.addEventListener("click", this._onYearClick.bind(this, year));
		return el;
	}

	_createYearDownButton() {
		const el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = `<i class="icon--mdi icon--mdi--arrow-down"></i>`;
		el.addEventListener("click", this._onYearDownClick.bind(this));
		return el;
	}

	_createYearUpButton() {
		const el = document.createElement("a");
		el.href = "javascript:void(0)";
		el.innerHTML = `<i class="icon--mdi icon--mdi--arrow-up"></i>`;
		el.addEventListener("click", this._onYearUpClick.bind(this));
		return el;
	}

}

customElements.define("date-time-picker", DateTimePicker);

export { DateTimePicker };
