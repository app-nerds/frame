document.addEventListener("DOMContentLoaded", () => {
  document.querySelector("#colorPicker").addEventListener("color-selected", onColorSelected);
  document.querySelector("#close").addEventListener("click", onCloseClick);

  function onColorSelected(e) {
    console.log(e);
  }

  function onCloseClick() {
    window.location = "/admin/roles/manage";
  }
});
