$(document).ready(function () {
    var webSocketConnect;
    var isContainerOpened = false;

    $("#open-log-container").click(function () {
        var targetValue;
        var buttonText;
        if (isContainerOpened) {
            targetValue = "-75vw";
            buttonText = "/"
        } else {
            targetValue = "0";
            buttonText = "X"
        }
        isContainerOpened = !isContainerOpened;

        $("#open-log-container").text(buttonText)

        $('#log-container').animate({
            right: targetValue
        }, {
            duration: 500,
            specialEasing: {
                right: "swing"
            }
        });
    })

    var terminalObject = $('#log-container').terminal(function (command) {
        console.log(command)

        if (!webSocketConnect) {
            return false;
        }
        if (!command) {
            return false;
        }
        webSocketConnect.send(JSON.stringify({msg: command}));
        return false;
    }, {
        greetings: `Welcome to the minecraft server mananger (mcman).
Let's start with writing \"help\"!

GitHub Repository: https://github.com/doppiolab/mcman
`,
        height: '100%',
        width: '100%',
        prompt: '[[g;#00ff00;>]mcman ➜] '
    });

    webSocketConnect = new WebSocket("ws://" + document.location.host + '/ws/terminal');
    webSocketConnect.onopen = function () {
        terminalObject.echo(getColoredMsg("[System] Connection opened.", "System"), { raw: true });
    };
    webSocketConnect.onclose = function (event) {
        terminalObject.echo(getColoredMsg("[System] Connection closed.", "System"), { raw: true });
    };
    webSocketConnect.onmessage = function (event) {
        var payload = JSON.parse(event.data)
        var msg = payload.msg.replaceAll("<", "&lt;").replaceAll(">", "&gt;")
        terminalObject.echo(getColoredMsg(msg, payload.type), { raw: true });
    };

    // load map
    $.ajax({
        method: "POST",
        url: "/api/v1/map",
        dataType: "json",
        contentType: "application/json",
        data : JSON.stringify( { x: 0, z: 0 }),
        success: function (data) {
            // map = L.map('map').setView([0, 0], 13)
        },
    })
});

function getColoredMsg(message, type) {
    if (type == "stderr")
        return `<span class='stderr-msg'>${message}</span>`
    if (type == "System")
        return `<span class='system-msg'>${message}</span>`
    return message
}
