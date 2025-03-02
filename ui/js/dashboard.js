$(document).ready(function () {
    const token = localStorage.getItem("jwt");
    console.log(token);

    if (!token) {
        window.location.href = "index.html";
        return;
    }

    // Проверка токена перед загрузкой страницы
    $.ajax({
        url: "http://api.127.0.0.1.nip.io/protected/read",
        method: "GET",
        headers: { "Authorization": "Bearer " + token },
        success: function (data) {
            $("#username").text(JSON.parse(data).username);
            $("#role").text(JSON.parse(data).role);
            if (JSON.parse(data).role === "write") {
                $("#writeData").removeClass("d-none");
            }
            $("#content").removeClass("d-none"); // Показываем контент
        },
        error: function () {
            localStorage.removeItem("jwt");
            window.location.href = "index.html";
        }
    });

    // Получение данных
    $("#getData").click(function () {
        $.ajax({
            url: "http://api.127.0.0.1.nip.io/protected/read",
            method: "GET",
            headers: { "Authorization": "Bearer " + token },
            success: function (data) {
                $("#data").text(JSON.stringify(data));
            }
        });
    });

    // Запись данных (только для write)
    $("#writeData").click(function () {
        $.ajax({
            url: "http://api.127.0.0.1.nip.io/protected/write",
            method: "POST",
            headers: { "Authorization": "Bearer " + token },
            success: function (data) {
                $("#writeMessage").text(JSON.parse(data).message || "Ошибка!");
            }
        });
    });

    // Выход
    $("#logout").click(function () {
        localStorage.removeItem("jwt");
        window.location.href = "index.html";
    });
});
