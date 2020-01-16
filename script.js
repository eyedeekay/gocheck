
document.addEventListener("click", clickEvent => {
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
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].removeClass("hidden");
  }
}

function showDeadOnly() {
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].removeClass("hidden");
  }
}

function showUnknownOnly() {
  var selects = document.getElementsByTagName("false");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("true");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].addClass("hidden");
  }
  var selects = document.getElementsByTagName("unknown");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].removeClass("hidden");
  }
}
