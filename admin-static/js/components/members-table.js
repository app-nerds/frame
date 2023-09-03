import { fetcher, PopupMenu } from "../frame.min.js";
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
    const th0 = document.createElement("th");
    const th1 = document.createElement("th");
    const th2 = document.createElement("th");
    const th3 = document.createElement("th");
    const th4 = document.createElement("th");
    const th5 = document.createElement("th");

    th0.setAttribute("scope", "col");
    th0.innerText = "Role";

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

    tr.insertAdjacentElement("beforeend", th0);
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
      td.setAttribute("colspan", "6");
      td.innerText = `No member records`;

      tr.insertAdjacentElement("beforeend", td);
      return [tr];
    }

    members.forEach(member => {
      const tr = document.createElement("tr");
      const td0 = document.createElement("td");
      const th1 = document.createElement("th");
      const td2 = document.createElement("td");
      const td3 = document.createElement("td");
      const td4 = document.createElement("td");
      const td5 = document.createElement("td");
      const buttons = this.createActionButtons(member);

      td0.innerHTML = `<span class="member-table-role-block" style="background-color: ${member.role.color};" title="Role: ${member.role.role}"><span class="sr-only">Role: ${member.role.role}</span></span>`;
      th1.setAttribute("scope", "row");
      th1.innerText = `${member.firstName} ${member.lastName}`;
      td2.innerText = member.email;
      td3.innerText = dayjs(member.CreatedAt).format("MMM D, YYYY");
      td4.innerText = member.memberStatus.status;

      buttons.forEach(button => {
        td5.insertAdjacentElement("beforeend", button);
      });

      tr.insertAdjacentElement("beforeend", td0);
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
    const buttonID = `member-action-${member.id}`;

    const button = document.createElement("button");
    button.id = buttonID;
    button.classList.add("action-button");
    button.setAttribute("alt", "Action Menu");
    button.setAttribute("title", "Action Menu");
    button.innerHTML = `<i class="icon--mdi icon--mdi--menu"></i>`;

    const popup = document.createElement("popup-menu");
    popup.setAttribute("trigger", `#${buttonID}`);

    let menuItems = [
      { id: `member-edit-button-${member.id}`, text: `Edit`, icon: "icon--mdi icon--mdi--pencil", handler: () => { this.onEditMemberClick(member.id); } },
    ];

    if (member.memberStatus.id === PendingApproval) {
      menuItems.push({ id: `member-status-button-${member.id}`, text: `Approve`, icon: "icon--mdi icon--mdi--check", handler: () => { this.onActionButtonClick(member); } });
    }

    if (member.memberStatus.id === Inactive) {
      menuItems.push({ id: `member-status-button-${member.id}`, text: `Inactivate`, icon: "icon--mdi icon--mdi--minus", handler: () => { this.onActionButtonClick(member); } });
    }

    if (member.memberStatus.id === Inactive) {
      menuItems.push({ id: `member-status-button-${member.id}`, text: `Activate`, icon: "icon--mdi icon--mdi--check", handler: () => { this.onActionButtonClick(member); } });
    }

    menuItems.push({ id: `member-delete-button-${member.id}`, text: `Delete`, icon: "icon--mdi icon--mdi--delete", handler: () => { this.onDeleteButtonClick(member); } });

    menuItems.forEach(data => {
      const menuItem = document.createElement("popup-menu-item");
      menuItem.setAttribute("id", data.id);
      menuItem.setAttribute("text", data.text);
      menuItem.setAttribute("icon", data.icon);

      popup.insertAdjacentElement("beforeend", menuItem);
    });

    popup.addEventListener("menu-item-click", e => {
      if (e.detail.id === `member-edit-button-${member.id}`) {
        this.onEditMemberClick.call(this, member.id);
      } else if (e.detail.id === `member-delete-button-${member.id}`) {
        this.onDeleteButtonClick.call(this, member);
      } else {
        this.onActionButtonClick.call(this, member)
      }
    });

    return [button, popup];
  }

  async getMembers() {
    const options = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    };

    const response = await fetcher(`/admin/api/members?page=${this._page}`, options);
    const result = await response.json();
    return result;
  }

  onEditMemberClick(memberID) {
    window.location = `/admin/members/edit/${memberID}`;
  }

  async onActionButtonClick(member) {
    let result;

    switch (member.memberStatus.id) {
      case PendingApproval:
        result = await this.activateMember(member);

        if (!result.success) {
          window.alert.error(result.message);
        } else {
          window.alert.success("Member approved!");
        }

        break;

      case Active:
        // TODO: Add inactivate
        break;

      case Inactive:
        result = await this.activateMember(member);

        if (!result.success) {
          window.alert.error(result.message);
        } else {
          window.alert.success("Member approved!");
        }

        break;
    }

    this.rerenderBody();
  }

  async onDeleteButtonClick(member) {
    const confirmation = await window.confirm.yesNo("Are you sure you wish to delete this member?");

    if (!confirmation) {
      return;
    }

    const options = {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    };

    const response = await fetcher(`/admin/api/member/delete/${member.id}`, options, window.spinner);
    const result = await response.json();

    if (!response.ok) {
      window.alert.error(result.message);
      return;
    }

    window.alert.success("Member deleted.");
    this.rerenderBody();
  }

  async activateMember(member) {
    const data = new URLSearchParams();
    data.set("id", member.id);

    const options = {
      method: "PUT",
      body: data,
    }

    const response = await fetcher(`/admin/api/member/activate`, options, window.spinner);
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
