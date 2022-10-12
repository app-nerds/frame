document.addEventListener("DOMContentLoaded", () => {
  document.querySelector("#cancel").addEventListener("click", () => {
    window.location = "/member/profile";
  });

  document.querySelector("#uploadForm").addEventListener("submit", (e) => {
    e.preventDefault();
    const el = document.querySelector("#imageFile");

    console.log(el.files);

    if (el.files.length <= 0) {
      window.alert.error("Please choose an image to upload!");
      return false;
    }

    e.target.submit();
    return true;
  });
});
