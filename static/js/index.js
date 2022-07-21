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

    map = L.map('map', {crs: L.CRS.Simple, minZoom: -10}).setView([0, 0], 0)

    $.ajax({
        method: "GET",
        url: "/api/v1/regions",
        dataType: "json",
        contentType: "application/json",
        success: function (data) {
            scale = 32 * 16
            for (var i = 0; i < data.length; i++) {
                var region = data[i];

                var bounds = [
                    [
                        -(region.Z) * scale,
                        (region.X + 1) * scale,
                    ],
                    [
                        -(region.Z + 1) * scale,
                        region.X * scale,
                    ]
                ];

                L.imageOverlay(`/api/v1/chunk/${region.X}/${region.Z}/map.png`, bounds).addTo(map);
            }
        },
    })

    $.ajax({
        method: "GET",
        url: "/api/v1/players",
        dataType: "json",
        contentType: "application/json",
        success: function (data) {
            for (var i = 0; i < data.length; i++) {
                var player = data[i];

                var headIcon = L.icon({
                    iconUrl: `https://crafatar.com/renders/head/${player.UUID}`,
                    iconSize:     [30, 30],
                });

                L.
                    marker([-player.Z, player.X], {icon: headIcon}).
                    addTo(map).
                    bindPopup(`Name: ${player.Name}<br/>UUID: ${player.UUID}<br/>Position: [x: ${player.X.toFixed(2)}, y: ${player.Y.toFixed(2)}, x: ${player.Z.toFixed(2)}]`);;
            }
        }
    })
});

function getColoredMsg(message, type) {
    if (type == "stderr")
        return `<span class='stderr-msg'>${message}</span>`
    if (type == "System")
        return `<span class='system-msg'>${message}</span>`
    return message
}
