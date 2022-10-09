import frame from "../frame.min.js";
import { PendingApproval, Active, Inactive } from "../constants/member-constants.js";

export default class MembersTable extends HTMLElement {
  constructor() {
    super();

    this._tbody = null;
    this._page = 1;
  }

  async connectedCallback() {
    const members = await this.getMembers();
    const table = this.createTable(members);

    this.insertAdjacentElement("beforeend", table);
    feather.replace();
  }

  createTable(members) {
    const el = document.createElement("table");
    const caption = document.createElement("caption");
    const head = this.createTableHead();
    const body = this.createTableBody(members);

    caption.innerText = "Site Members";
    el.insertAdjacentElement("beforeend", caption);
    el.insertAdjacentElement("beforeend", head);
    el.insertAdjacentElement("beforeend", body);

    return el;
  }

  createTableHead() {
    const head = document.createElement("thead");
    const tr = document.createElement("tr");
    const th1 = document.createElement("th");
    const th2 = document.createElement("th");
    const th3 = document.createElement("th");
    const th4 = document.createElement("th");
    const th5 = document.createElement("th");

    th1.setAttribute("scope", "col");
    th1.innerText = "Name";

    th2.setAttribute("scope", "col");
    th2.innerText = "Email";

    th3.setAttribute("scope", "col");
    th3.innerText = "Member Since";

    th4.setAttribute("scope", "col");
    th4.innerText = "Status";

    th5.setAttribute("scope", "col");
    th5.innerHTML = `<span class="sr-only">Actions</span>`;

    tr.insertAdjacentElement("beforeend", th1);
    tr.insertAdjacentElement("beforeend", th2);
    tr.insertAdjacentElement("beforeend", th3);
    tr.insertAdjacentElement("beforeend", th4);
    tr.insertAdjacentElement("beforeend", th5);

    head.insertAdjacentElement("beforeend", tr);
    return head;
  }

  createTableBody(members) {
    this._tbody = document.createElement("tbody");
    const rowEls = this.createTableBodyContents(members);

    rowEls.forEach(el => {
      this._tbody.insertAdjacentElement("beforeend", el);
    });

    return this._tbody;
  }

  createTableBodyContents(members) {
    let result = [];

    if (members.length <= 0) {
      const tr = document.createElement("tr");
      const td = document.createElement("td");

      td.setAttribute("scope", "row");
      td.innerText = `No member records`;

      tr.insertAdjacentElement("beforeend", td);
      return [tr];
    }

    members.forEach(member => {
      const tr = document.createElement("tr");
      const th1 = document.createElement("th");
      const td2 = document.createElement("td");
      const td3 = document.createElement("td");
      const td4 = document.createElement("td");
      const td5 = document.createElement("td");
      const buttons = this.createActionButtons(member);

      th1.setAttribute("scope", "row");
      th1.innerText = `${member.firstName} ${member.lastName}`;

      td2.innerText = member.email;
      td3.innerText = dayjs(member.CreatedAt).format("MMM D, YYYY");
      td4.innerText = member.memberStatus.status;

      buttons.forEach(button => {
        td5.insertAdjacentElement("beforeend", button);
      });

      tr.insertAdjacentElement("beforeend", th1);
      tr.insertAdjacentElement("beforeend", td2);
      tr.insertAdjacentElement("beforeend", td3);
      tr.insertAdjacentElement("beforeend", td4);
      tr.insertAdjacentElement("beforeend", td5);

      result.push(tr);
    });

    return result;
  }

  createActionButtons(member) {
    const button1 = document.createElement("button");

    if (member.memberStatus.id === PendingApproval) {
      button1.classList.add("action-button");
      button1.setAttribute("alt", `Approve ${member.firstName} ${member.lastName}`);
      button1.setAttribute("title", `Approve ${member.firstName} ${member.lastName}`);
      button1.innerHTML = `<i data-feather="user-check"></i>`;
    }

    if (member.memberStatus.id === Active) {
      button1.classList.add("delete-button");
      button1.setAttribute("alt", `Inactivate ${member.firstName} ${member.lastName}`);
      button1.setAttribute("title", `Inactivate ${member.firstName} ${member.lastName}`);
      button1.innerHTML = `<i data-feather="user-minus"></i>`;
    }

    if (member.memberStatus.id === Inactive) {
      button1.classList.add("action-button");
      button1.setAttribute("alt", `Activate ${member.firstName} ${member.lastName}`);
      button1.setAttribute("title", `Activate ${member.firstName} ${member.lastName}`);
      button1.innerHTML = `<i data-feather="user-check"></i>`;
    }

    button1.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation();

      this.onActionButtonClick(member);
    });

    return [button1];
  }

  async getMembers() {
    const options = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    };

    const response = await frame.fetcher(`/admin/api/members?page=${this._page}`, options);
    const result = await response.json();
    return result;
  }

  async onActionButtonClick(member) {
    let result;
    console.log(member);

    switch (member.memberStatus.id) {
      case PendingApproval:
        result = await this.activateMember(member);

        if (!result.success) {
          console.log(result);
          window.alert.error(result.message);
        } else {
          window.alert.success("Member approved!");
        }

        break;

      case Active:
        break;

      case Inactive:
        result = await this.activateMember(member);

        if (!result.success) {
          console.log(result);
          window.alert.error(result.message);
        } else {
          window.alert.success("Member approved!");
        }

        break;
    }

    this.rerenderBody();
  }

  async inactivateMember(member) {

  }

  async activateMember(member) {
    const data = new URLSearchParams();
    data.set("id", member.ID);

    const options = {
      method: "PUT",
      body: data,
    }

    const response = await fetch(`/admin/api/member/activate`, options);
    const result = await response.json();
    return result;
  }

  async rerenderBody() {
    const members = await this.getMembers();

    this._tbody.innerHTML = "";
    const rowEls = this.createTableBodyContents(members);

    rowEls.forEach(el => {
      this._tbody.insertAdjacentElement("beforeend", el);
    });

    feather.replace();
  }
}

customElements.define("members-table", MembersTable);
