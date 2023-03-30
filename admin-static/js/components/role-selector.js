import frame from "../frame.min.js";

export default class RoleSelector extends HTMLElement {
  constructor() {
    super();

    this.selected = this.getAttribute("selected") || "0";
    this.name = this.getAttribute("name") || "role";

    this.selectEl = document.createElement("select");
    this.selectEl.name = this.name;
    this.selectEl.addEventListener("change", this.onItemSelected.bind(this));
    this.appendChild(this.selectEl);
  }

  async connectedCallback() {
    const roles = await this.getRoles();
    this.attachSelectOptions(roles);
  }

  async getRoles() {
    const options = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    };

    const response = await frame.fetcher(`/admin/api/member/role`, options, window.spinner);
    const result = await response.json();

    if (!response.ok) {
      console.log(result);
      throw new Error(result.message);
    }

    return result;
  }

  attachSelectOptions(roles) {
    this.selectEl.options.length = 0;

    roles.forEach(item => {
      const selected = window.parseInt(this.selected) === item.id;
      this.selectEl.options.add(new Option(item.role, item.id, selected, selected));
    });
  }

  onItemSelected(e) {
    this.dispatchEvent(new CustomEvent("role-selected", { data: e.target.value }));
  }
}

customElements.define("role-selector", RoleSelector);
