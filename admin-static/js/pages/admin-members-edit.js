import RoleSelector from "../components/role-selector.js";

document.addEventListener("DOMContentLoaded", () => {
  document.querySelector("#cancel").addEventListener("click", onCancelClick);

  /*
   * Event handlers
   */
  function onCancelClick() {
    window.location = "/admin/members/manage";
  }
});
