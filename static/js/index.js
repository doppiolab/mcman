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
        prompt: 'mcman> '
    });

    webSocketConnect = new WebSocket("ws://" + document.location.host + '/ws/terminal');
    webSocketConnect.onopen = function () {
        terminalObject.echo(getColoredMsg("[System] Connection opened."), { raw: true });
    };
    webSocketConnect.onclose = function (event) {
        terminalObject.echo(getColoredMsg("[System] Connection closed."), { raw: true });
    };
    webSocketConnect.onmessage = function (event) {
        payload = JSON.parse(event.data)
        terminalObject.echo(getColoredMsg(payload.msg), { raw: true });
    };
});

function getColoredMsg(message) {
    if (message.includes("[stdout]"))
        return message.replace("[stdout]", "<span class='system-stdout'>[stdout]</span>")
    if (message.includes("[stderr]"))
        return message.replace("[stderr]", "<span class='system-stderr'>[stderr]</span>")
    if (message.includes("[System]"))
        return message.replace("[System]", "<span class='system-msg'>[System]</span>")
    return message
}
