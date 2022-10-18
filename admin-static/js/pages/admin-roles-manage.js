document.addEventListener("DOMContentLoaded", () => {
  const popupMenuItems = document.querySelectorAll(".role-popup-menu-item");

  for (let i = 0; i < popupMenuItems.length; i++) {
    popupMenuItems[i].addEventListener("click", onPopupMenuItemClicked);
  }

  function onPopupMenuItemClicked(e) {
    const [_, id] = e.target.id.split("_");
    window.location = `/admin/roles/edit/${id}`;
  }
});
