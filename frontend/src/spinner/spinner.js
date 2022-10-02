/*
 * spinner is simple library for displaying a loading spinner. It makes use
 * of the whole page to display the spinner. The spinner is pure CSS, SVG, and JavaScript.
 * Copyright © 2022 App Nerds LLC
 */
export function spinner() {
  let spinner = undefined;

  function calculateHeight() {
    const body = document.body;
    const html = document.documentElement;

    const bodyHeight = body.getBoundingClientRect().height;
    const htmlHeight = html.getBoundingClientRect().height;

    const height = Math.max(bodyHeight, htmlHeight);
    return height;
  }

  function calculateTop() {
    const r = document.body.getBoundingClientRect();
    return Math.abs(r.top);
  }

  function show() {
    if (!spinner) {
      spinner = document.createElement("div");
      spinner.classList.add("spinner");
      spinner.style.top = `${calculateTop()}px`;
      spinner.style.height = `${calculateHeight()}px`;
      spinner.innerHTML = `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
          <circle cx="50" cy="50" r="45" />
        </svg>
      `;

      document.body.appendChild(spinner);
    }
  }

  function hide() {
    if (spinner) {
      spinner.remove();
      spinner = undefined;
    }
  }

  return {
    hide() {
      hide();
    },

    show() {
      show();
    },
  }
}

