$(document).ready(function () {
    $("#open-log-container").click(function () {
        alert("Clicked")
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
