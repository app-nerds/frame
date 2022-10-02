/*
 * confirm is a function to display a confirmation dialog. It has two mode: "yesno", "other". 
 * "yesno" mode will display two buttons: Yes and No. "other" will only display a Close button.
 * The result of the click will be returned in a promise value.
 *
 * Styling is provided by confirm.css. It relies on variables: 
 *   - --dialog-background-color 
 *   - --border-color
 *
 * Example:
 *    const confirmation = confirm();
 *    const result = await confirmation.yesNo("Are you sure?");
 *
 * Copyright © 2022 App Nerds LLC
 */

import { shim } from "../shim/shim.js";

export function confirm(baseOptions = {
  width: "25%",
  height: "25%",
  callback: undefined,
}) {
  const shimBuilder = shim({ closeOnClick: true });
  let _shim;

  function calculateTop() {
    const r = document.body.getBoundingClientRect();
    return Math.abs(r.top);
  }

  function show(type, message, options) {
    options = { ...baseOptions, ...options };

    const container = document.createElement("div");
    container.classList.add("confirm-container");
    container.style.setProperty("width", options.width);
    container.style.setProperty("height", options.height);
    container.style.setProperty("top", `calc((50% - ${options.height}) + ${calculateTop()}px)`);
    container.style.setProperty("left", `calc(50% - (${options.width}/2))`);

    _shim = shimBuilder.new({ callback: () => { close(container, options.callback, false); } });

    setContent(container, message);
    addButtons(container, type, options.callback);

    _shim.show();
    document.body.appendChild(container);
  }

  function setContent(container, message) {
    container.innerHTML += `<p>${message}</p>`;
  }

  function addButtons(container, type, callback) {
    let buttons = [];

    switch (type) {
      case "yesno":
        const noB = document.createElement("button");
        noB.innerText = "No";
        noB.classList.add("cancel-button");
        noB.addEventListener("click", (e) => {
          e.preventDefault();
          e.stopPropagation();

          _shim.hide(false);
          close(container, callback, false)
        });

        const yesB = document.createElement("button");
        yesB.innerText = "Yes";
        yesB.classList.add("action-button");
        yesB.addEventListener("click", (e) => {
          e.preventDefault();
          e.stopPropagation();

          _shim.hide(false);
          close(container, callback, true);
        });

        buttons.push(noB);
        buttons.push(yesB);
        break;

      default:
        const b = document.createElement("button");
        b.innerText = "Close";
        b.classList.add("action-button");
        b.addEventListener("click", (e) => {
          e.preventDefault();
          e.stopPropagation();

          _shim.hide(false);
          close(container, callback);
        });

        buttons.push(b);
        break;
    }

    const buttonContainer = document.createElement("div");
    buttonContainer.classList.add("button-row");

    buttons.forEach((button) => { buttonContainer.appendChild(button); });
    container.appendChild(buttonContainer);
  }

  function close(container, callback, callbackValue) {
    container.remove();
    if (typeof callback === "function") {
      callback(callbackValue);
    }
  }

  return {
    confirm(message, options) {
      show("confirm", message, options);
    },

    yesNo(message, options) {
      return new Promise((resolve) => {
        const cb = (result) => {
          return resolve(result);
        };

        options = { ...{ callback: cb }, ...options };
        show("yesno", message, options);
      });
    },
  };
}

