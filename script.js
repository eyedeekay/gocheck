
document.addEventListener("click", clickEvent => {
  console.log("Showing a class of site");
  if (clickEvent.target.id === "alive") {
    showAliveOnly()
  }
  if (clickEvent.target.id === "dead") {
    showDeadOnly()
  }
  if (clickEvent.target.id === "unknown") {
    showUnknownOnly()
  }
});

function showAliveOnly() {
  console.log("Showing alive only");
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("removing class hidden")
     selects[i].removeClass("hidden");
  }
}

function showDeadOnly() {
  console.log("Showing dead only");
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("removing class hidden")
     selects[i].removeClass("hidden");
  }
}

function showUnknownOnly() {
  console.log("Showing unknown only");
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     console.log("adding class hidden")
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].removeClass("hidden");
     console.log("removing class hidden")
  }
}
