
function showError(str, code) {
  errorElement = document.getElementById("error");
  errorElement.innerText = str;

  if (code !== undefined) {
    const pre = document.createElement("pre");
    pre.innerText = code
    errorElement.appendChild(pre)
  }

  errorElement.style.display = "block";

}

function hide(obj) {
  obj.style.display = "none";
}
