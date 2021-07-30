package models

var Qrocde = `<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>二维码</title>
    <style>
        * {
            margin: 0;
            padding: 0;
        }
        .middle {
            position: fixed;
            top: 0;
            right: 0;
            bottom: 0;
            left: 0;
            margin: auto;
        }
        #button {
            text-align: center;
            height: 2rem;
            width: 3rem;
            font-size: 1rem;
            line-height: 2rem;
            opacity: 0;
            color: white;
        }
    </style>
</head>

<body>
    <img id="qrcode" height="170em" class="middle" src=""
    />
    <div id="button" class="middle" onclick="refresh()">刷新</div>

    <script>
        var timer;
        var qrcode = document.getElementById("qrcode")
        var button = document.getElementById("button")
        function showQrcode() {
            var xmlhttp;
            if (window.XMLHttpRequest) {
                xmlhttp = new XMLHttpRequest();
            } else {
                xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
            }
            xmlhttp.onreadystatechange = function() {
                if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
                    qrcode.src = xmlhttp.responseText
                }
            }
            xmlhttp.open("GET", "/api/login/qrcode", true);
            xmlhttp.send();
        }
        showQrcode()
        function login() {
            var xmlhttp;
            if (window.XMLHttpRequest) {
                xmlhttp = new XMLHttpRequest();
            } else {
                xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
            }
            xmlhttp.onreadystatechange = function() {
                if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
                    var results = xmlhttp.responseText
                    console.log(results)
                    switch (results) {
                        case "授权登录未确认":
                            break;
                        case "登录":
                            window.location.href = "/admin"
                            break;
                        case "成功":
                            button.innerText = "成功"
                            button.style.opacity = 1
                            button.style.backgroundColor = "#1E9FFF"
                            qrcode.style.opacity = 0.4
                            clearInterval(timer)
                            break;
                        default:
                            showQrcode()
                            break;
                    }
                }
            }
            xmlhttp.open("GET", "/api/login/query", true);
            xmlhttp.send();
        }

        function polling() {
            timer = setInterval(() => {
                login()
            }, 1500);
        }
        polling()

        function refresh() {
            if (button.style.opacity == 1) {
                showQrcode()
                button.style.opacity = 0
                qrcode.style.opacity = 1
                polling()
            } else {
                login()
            }
        }
    </script>
</body>


</html>`
