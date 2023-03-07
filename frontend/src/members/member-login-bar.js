/*
 * MemberLoginBar is a component used to display a member dropdown in the header of websites. 
 * It displays either a user-uploaded image or the letter of the first initial of the user's name. 
 * When logged in the menu provides links to the user's account and log off. If the user is not logged 
 * in then a log in link is displayed.
 *
 * To work with member data this component requires service component that provides the following.
 *   - getCurrentMember - Must return an object with fields memberID, firstName, lastName, profilePictureURL
 *
 * This component uses Feather Icons. https://feathericons.com/
 * 
 * Copyright © 2022 App Nerds LLC
*/

import { MemberService } from "./member-service.js";

export default class MemberLoginBar extends HTMLElement {
  memberService;
  loginPath;

  constructor() {
    super();

    this.loginPath = this.getAttribute("login-path") || "/member/login";
    const spinner = this.getAttribute("spinner") || null;
    let spinnerEl = null;

    if (spinner) {
      spinnerEl = document.querySelector(spinner);
    }

    this.memberService = new MemberService(spinnerEl);
  }

  static get observedAttributes() {
    return ["login-path"];
  }

  set memberService(/** @type {any} */ service) {
    this.memberService = service;
  }

  attributedChangedCallback(name, _, newValue) {
    if (name === "login-path") {
      this.loginPath = newValue;
    }
  }

  async connectedCallback() {
    let member = null;

    member = await this.memberService.getCurrentMember();

    const containerEl = this.createContainerEl();
    this.createAvatarEl(containerEl, member);
    this.createTextEl(containerEl, member);
    this.createPopupMenu(containerEl, member);

    this.insertAdjacentElement("beforeend", containerEl);
    feather.replace();
  }

  /*******************************************************************************
   * Event methods
   ******************************************************************************/

  /*******************************************************************************
   * UI elements
   ******************************************************************************/

  createContainerEl() {
    const el = document.createElement("div");
    return el;
  }

  createAvatarEl(container, member) {
    let el;

    if (member && member.avatarURL) {
      el = document.createElement("img");
      el.classList.add("avatar");
      el.setAttribute("role", "img");
      el.src = member.avatarURL;
    } else {
      el = document.createElement("div");
      el.classList.add("avatar");
      el.setAttribute("role", "img");
      el.innerHTML = `<i data-feather="user"></i>`;
    }

    container.insertAdjacentElement("beforeend", el);
  }

  createTextEl(container, member) {
    let markup;

    const el = document.createElement("a");
    el.id = "member-link";

    if (member && member.email) {
      let name = "";

      el.href = "#";

      if (member.firstName) {
        name += member.firstName;
      }

      if (member.lastName) {
        name += ` ${member.lastName}`;
      }

      if (name === "") {
        name = "User";
      }

      markup = `${name} <i data-feather="chevron-down"></i>`;
    } else {
      el.href = this.loginPath;
      markup = "Log In";
    }

    el.innerHTML = markup;
    container.insertAdjacentElement("beforeend", el);
  }

  createPopupMenu(container, member) {
    if (member && member.email) {
      const el = document.createElement("popup-menu");
      el.setAttribute("trigger", "#member-link");

      const menuItems = [
        { id: "member-my-account-link", text: "My Account", icon: "home", handler: this.onMyAccountClick.bind(this) },
        { id: "member-log-out-link", text: "Log Out", icon: "log-out", handler: this.onLogOutClick.bind(this) },
      ];

      menuItems.forEach(data => {
        const menuItem = document.createElement("popup-menu-item");
        menuItem.setAttribute("id", data.id);
        menuItem.setAttribute("text", data.text);
        menuItem.setAttribute("icon", data.icon);
        menuItem.addEventListener("click", data.handler);

        el.insertAdjacentElement("beforeend", menuItem);
      });

      container.insertAdjacentElement("beforeend", el);
    }
  }

  /*******************************************************************************
   * Private methods
   ******************************************************************************/

  onMyAccountClick() {
    window.location = "/member/profile";
  }

  onLogOutClick() {
    window.location = "/api/member/logout";
  }
}

if (!customElements.get("member-login-bar")) {
  customElements.define("member-login-bar", MemberLoginBar);
}

