
function showAliveOnly(){
  var selects = document.getElementsByTagName("select");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].className += " hidden";
  }
}

function showDeadOnly(){
  var selects = document.getElementsByTagName("select");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].className += " hidden";
  }
}

function showUnknownOnly(){
  var selects = document.getElementsByTagName("select");
  for(var i =0, il = selects.length;i<il;i++){
     selects[i].className += " hidden";
  }
}