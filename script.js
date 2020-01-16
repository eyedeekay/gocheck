document.addEventListener("click", clickEvent => {
  console.log("Showing a class of site");
  if (clickEvent.target.id === "alive") {
    showAliveOnly();
  }
  if (clickEvent.target.id === "dead") {
    showDeadOnly();
  }
  if (clickEvent.target.id === "unknown") {
    showUnknownOnly();
  }
});

function showAliveOnly() {
  console.log("Showing alive only");
  var selects = document.getElementsByClassName("false");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("unknown");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("true");
  for (let select in selects) {
    console.log("removing class hidden");
    select.removeClass("hidden");
  }
}

function showDeadOnly() {
  console.log("Showing dead only");
  var selects = document.getElementsByClassName("true");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("unknown");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("false");
  for (let select in selects) {
    console.log("removing class hidden");
    select.removeClass("hidden");
  }
}

function showUnknownOnly() {
  console.log("Showing unknown only");
  var selects = document.getElementsByClassName("false");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("true");
  for (let select in selects) {
    console.log("adding class hidden");
    select.addClass("hidden");
  }
  var selects = document.getElementsByClassName("unknown");
  for (let select in selects) {
    select.removeClass("hidden");
    console.log("removing class hidden");
  }
}
