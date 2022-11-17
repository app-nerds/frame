import { alert } from "../src/dialogs/alert.js";
import { confirm } from "../src/dialogs/confirm.js";
import { spinner } from "../src/spinner/spinner.js";
import { PopupMenu, PopupMenuItem } from "../src/menus/popup-menu.js";
import ColorPicker from "../src/color-picker/color-picker.js";
import MessageBar from "../src/message-bar/message-bar.js";

document.addEventListener("DOMContentLoaded", () => {
  const alerter = alert();
  const confirmer = confirm();
  const spinnerer = spinner();

  feather.replace();

  /* 
   * Click button to show alert 
   */
  document.getElementById("showInfoAlert").addEventListener("click", () => {
    alerter.info("This is an information alert!");
  });

  document.getElementById("showSuccessAlert").addEventListener("click", () => {
    alerter.success("This is a success alert!");
  });

  document.getElementById("showWarnAlert").addEventListener("click", () => {
    alerter.warn("This is a warning alert!");
  });

  document.getElementById("showErrorAlert").addEventListener("click", () => {
    alerter.error("This is an error alert!");
  });

  /*
   * Click to show confirm 
   */
  document.getElementById("showConfirm").addEventListener("click", async () => {
    const result = await confirmer.yesNo("Are you sure?");
    alerter.info(`You chose "${result}"`);
  });

  /*
   * Click to show spinner 
   */
  document.getElementById("showSpinner").addEventListener("click", () => {
    spinnerer.show();

    setTimeout(() => {
      spinnerer.hide();
    }, 3000);
  });

  /*
   * Popup menu events
   */
  document.getElementById("menuItem1").addEventListener("click", () => {
    alerter.info("Menu item 1");
  });

  document.getElementById("menuItem2").addEventListener("click", () => {
    alerter.error("Menu item 1");
  });

  document.getElementById("menuItem3").addEventListener("click", () => {
    alerter.warn("Menu item 1");
  });
});

