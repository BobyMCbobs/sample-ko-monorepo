function httpGet(url) {
  var xmlHttp = new XMLHttpRequest()
  xmlHttp.open("GET", url, false )
  xmlHttp.send(null)
  return xmlHttp.responseText
}

function getNumber() {
  return httpGet("/api/number")
}

document.addEventListener("DOMContentLoaded", () => {
  let magicNumber = document.querySelector('.magic-number')
  let number = getNumber()
  magicNumber.innerText = number

})
