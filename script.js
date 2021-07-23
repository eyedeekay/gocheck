document.addEventListener("click", (clickEvent) => {
  if (clickEvent.target.id === "alive") {
    showAliveOnly();
  }
  if (clickEvent.target.id === "dead") {
    showDeadOnly();
  }
  if (clickEvent.target.id === "untested") {
    showUnknownOnly();
  }
});

function hide(selected) {
  //if (selected.classList != undefined)
  selected.classList.add("hidden");
  selected.classList.remove("show");
}

function show(selected) {
  //if (selected.classList != undefined)
  selected.classList.remove("hidden");
  selected.classList.add("show");
}

function hideGroup(group) {
  var unselects = document.getElementsByClassName(group);
  console.log("  adding class hidden"); //, unselects);
  for (let select of unselects) {
    hide(select);
  }
}

function showGroup(group) {
  var unselects = document.getElementsByClassName(group);
  console.log("  removing class hidden"); //, unselects);
  for (let select of unselects) {
    show(select);
  }
}

function group(hide, otherhide, show) {
  hideGroup(hide);
  hideGroup(otherhide);
  showGroup(show);
}

function showAliveOnly() {
  console.log("Showing only alive");
  group("false", "unknown", "true");
}

function showDeadOnly() {
  console.log("Showing only dead");
  group("true", "unknown", "false");
}

function showUnknownOnly() {
  console.log("Showing only unknown");
  group("false", "true", "unknown");
}

//showAliveOnly();
