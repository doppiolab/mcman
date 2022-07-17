$(document).ready(function () {
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

    $('#log-container').terminal(function (command) {
        this.echo(command)
    }, {
        greetings: `Welcome to the minecraft server mananger (mcman).
Let's start with writing \"help\"!

GitHub Repository: https://github.com/doppiolab/mcman
`,
        height: '100%',
        width: '100%',
        prompt: 'mcman> '
    });
});
