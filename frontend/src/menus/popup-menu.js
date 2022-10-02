/*
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
 * Copyright © 2022 App Nerds LLC
 */

export class PopupMenu extends HTMLElement {
  constructor() {
    super();
    this._trigger = null;
    this._el = null;
  }

  connectedCallback() {
    this._trigger = this.getAttribute("trigger");

    if (!this._trigger) {
      throw new Error(
        "You must provide a query selector for the element used to trigger this popup."
      );
    }

    this.style.visibility = "hidden";

    document
      .querySelector(this._trigger)
      .addEventListener("click", this.toggle.bind(this));
  }

  disconnectedCallback() {
    let el = document.querySelector(this._trigger);

    if (el) {
      el.removeEventListener("click", this.toggle.bind(this));
    }
  }

  toggle(e) {
    if (e) {
      e.preventDefault();
    }

    let triggerRect = document
      .querySelector(this._trigger)
      .getBoundingClientRect();
    let thisRect = this.getBoundingClientRect();
    let buffer = 3;

    if (thisRect.right > window.innerWidth) {
      this.style.left =
        "" +
        (triggerRect.x + (window.innerWidth - thisRect.right) - buffer) +
        "px";
    } else {
      this.style.left = "" + triggerRect.x + "px";
    }

    this.style.top =
      "" + (triggerRect.y + triggerRect.height + buffer) + "px";

    if (this.style.visibility === "hidden") {
      this.style.visibility = "visible";
    } else {
      this.style.visibility = "hidden";
    }
  }
}

export class PopupMenuItem extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  render() {
    let text = this.getAttribute("text");
    let icon = this.getAttribute("icon");

    const div = document.createElement("div");
    div.classList.add("popup-menu-item");

    let inner = "";

    if (icon) {
      inner += `<i data-feather="${icon}"></i> `;
    }

    inner += text;
    div.innerHTML = inner;

    div.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation();
      this.dispatchEvent(new CustomEvent("click", { detail: e }));

      const parent = e.target.parentElement.parentElement;
      parent.style.visibility = "hidden";
    });

    this.insertAdjacentElement("beforeend", div);
  }
}

export const showPopup = (el) => {
  document.querySelector(el).style.visibility = "visible";
};

export const hidePopup = (el) => {
  document.querySelector(el).style.visibility = "hidden";
};

if (!customElements.get("popup-menu")) {
  customElements.define("popup-menu", PopupMenu);
}

if (!customElements.get("popup-menu-item")) {
  customElements.define("popup-menu-item", PopupMenuItem);
}

