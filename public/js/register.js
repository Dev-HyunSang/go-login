'use strict';

$("#form").keypress(function(e) {
    if (e.keyCode === 13) {
        Register();
    }
});

function Register() {
    const FirstName = document.getElementById("FirstName").value;
    const LastName = document.getElementById("LastName").value;
    const Email = document.getElementById("Email").value;
    const Passwrod = document.getElementById("Password").value;

    var data = {
        "first_name": FirstName,
        "last_name": LastName,
        "email": Email,
        "password": Passwrod
    }

    $.ajax({
        url: '/register/new', 
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function (completion) {
            console.log(completion);
            alert('회원가입 완료!');
            location.href="/login";
        },
        error: function (err){
            console.error(err);
        }
    });
}