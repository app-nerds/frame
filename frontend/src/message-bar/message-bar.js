/*
 * MessageBar is a component used to display a message on the screen. It is typically used to display 
 * the results of submitting a form. It can also be used to provide informational breakout.
 *
 * Copyright © 2022 App Nerds LLC
*/

export default class MessageBar extends HTMLElement {
  constructor() {
    super();

    this.messageType = this.getAttribute("message-type") || "info";
    this.message = this.getAttribute("message") || "";

    this.containerEl = null;
  }

  connectedCallback() {
    this.containerEl = this.createContainerEl();
    const closeButtonEl = this.createCloseButtonEl();
    const textEl = this.createTextEl();

    this.containerEl.insertAdjacentElement("beforeend", closeButtonEl);
    this.containerEl.insertAdjacentElement("beforeend", textEl);

    this.insertAdjacentElement("beforeend", this.containerEl);
  }

  createContainerEl() {
    const el = document.createElement("div");
    el.classList.add("message-bar");

    switch (this.messageType) {
      case "error":
        el.classList.add("message-bar-error");
        break;

      case "warn":
        el.classList.add("message-bar-warn");
        break;

      case "info":
        el.classList.add("message-bar-info");
        break;

      case "success":
        el.classList.add("message-bar-success");
        break;
    }

    return el;
  }

  createCloseButtonEl() {
    const el = document.createElement("span");
    el.innerHTML = "&times;";

    el.addEventListener("click", () => {
      if (this.containerEl) {
        this.containerEl.remove();
      }
    });

    return el;
  }

  createTextEl() {
    const el = document.createElement("p");
    el.setAttribute("role", "alert");
    el.innerHTML = this.message;

    return el;
  }
}

customElements.define("message-bar", MessageBar);
