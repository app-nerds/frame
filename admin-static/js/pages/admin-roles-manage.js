document.addEventListener("DOMContentLoaded", () => {
  const popupMenuEls = document.querySelectorAll(`popup-menu[id*="role-menu"]`);

  popupMenuEls.forEach(el => {
    el.addEventListener("menu-item-click", onPopupMenuItemClicked);
  })

  function onPopupMenuItemClicked(e) {
    const [_, id] = e.detail.id.split("_");
    window.location = `/admin/roles/edit/${id}`;
  }
});
