import frame from "/admin-static/js/frame.min.js";

window.spinner = frame.spinner();
window.alert = frame.alert();
window.confirm = frame.confirm();

document.addEventListener("DOMContentLoaded", () => {
  feather.replace();
});

