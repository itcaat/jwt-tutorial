$(document).ready(function () {
    $("#loginForm").submit(function (event) {
        event.preventDefault();
        const username = $("#username").val();
        const password = $("#password").val();

        $.ajax({
            url: "http://api.127.0.0.1.nip.io/login",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({ username, password }),
            success: function (data) {
                localStorage.setItem("jwt", JSON.parse(data).token);
                window.location.href = "dashboard.html";
            },
            error: function () {
                $("#errorMessage").removeClass("d-none");
            }
        });
    });
});
