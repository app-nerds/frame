/*
 * Copyright © 2022 App Nerds LLC
 */

export default class GoogleLoginForm extends HTMLElement {
  loginPath;
  createAccountPath;
  signInButtonURL;

  constructor() {
    super();

    this.loginPath = this.getAttribute("login-path") || "/auth/google";
    this.createAccountPath = this.getAttribute("create-account-path") || "/create-account";
    this.signInButtonURL = this.getAttribute("sign-in-button-url") || "/static/images/sign-in-with-google.jpg";
  }

  connectedCallback() {
    const sectionEl = document.createElement("section");
    sectionEl.classList.add("google-login-form");

    const footerEl = document.createElement("div");
    footerEl.classList.add("sign-up-footer");

    sectionEl.innerHTML = `
      <a href="${this.loginPath}"><img src="${this.signInButtonURL}" alt="Sign in with Google" style="width:100%;" /></a>
    `;

    footerEl.innerHTML = `
      <p>
        Don't have an account? Click <a href="${this.createAccountPath}">here</a> to create one.
      </p>
    `;

    sectionEl.insertAdjacentElement("beforeend", footerEl);
    this.insertAdjacentElement("beforeend", sectionEl);
  }
}

if (!customElements.get("google-login-form")) {
  customElements.define("google-login-form", GoogleLoginForm);
}

