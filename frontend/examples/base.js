import nerdjslibrary from "./nerdjslibrary.min.js";

document.addEventListener("DOMContentLoaded", () => {
  const alert = nerdjslibrary.alert();
  const confirm = nerdjslibrary.confirm();
  const spinner = nerdjslibrary.spinner();

  feather.replace();

  /* 
   * Click button to show alert 
   */
  document.getElementById("showAlert").addEventListener("click", () => {
    alert.info("This is an information alert!");
  });

  /*
   * Click to show confirm 
   */
  document.getElementById("showConfirm").addEventListener("click", async () => {
    const result = await confirm.yesNo("Are you sure?");
    alert.info(`You chose "${result}"`);
  });

  /*
   * Click to show spinner 
   */
  document.getElementById("showSpinner").addEventListener("click", () => {
    spinner.show();

    setTimeout(() => {
      spinner.hide();
    }, 3000);
  });

  /*
   * Popup menu events
   */
  document.getElementById("menuItem1").addEventListener("click", () => {
    alert.info("Menu item 1");
  });

  document.getElementById("menuItem2").addEventListener("click", () => {
    alert.error("Menu item 1");
  });

  document.getElementById("menuItem3").addEventListener("click", () => {
    alert.warn("Menu item 1");
  });
});

