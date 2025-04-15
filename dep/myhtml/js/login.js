
const API_HOST = 'http://10.23.50.159:2015'; // 替换为你的服务器地址

const CodeAck = {
    E_OK: 0,
	E_ERROR: 1,
	E_EXISTS: 2,
    E_NOT_EXISTS: 3,
	E_PWD_WRONG: 4,
}

function on_click_login() {
    var name = document.getElementById("name").value;
    var pwd = document.getElementById("pwd").value;
    var data = {};
    data["name"] = name;
    data["pwd"] = pwd;
    console.log(data);
    do_login(data)
}

function on_click_register() {
    var name = document.getElementById("name").value;
    var pwd = document.getElementById("pwd").value;
    var pwd_again = document.getElementById("pwd2").value;
    if (name == "" || pwd == "") {
        alert("用户名和密码不能为空!");
        return;
    }
    if (pwd == pwd_again) {
        var data = {};
        data["name"] = name;
        data["pwd"] = pwd;
        console.log(data);
        do_register(data);
    } else {
        alert("密码不一致!");
    }
}

// 发送 POST 请求
function do_register(data) {
    fetch(`${API_HOST}/api/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('网络响应不正常');
        }
        return response.json();
    })
    .then(data => {
        console.log('服务器响应:', data);
        if (data.code == CodeAck.E_OK) {
            alert("注册成功!");
            window.location.href = "login.html";
        } else {
            alert("注册失败! 错误码: " + data.code);
        }
    })
    .catch(error => {
            console.error('发送数据时出错:', error);
    });
}

function on_after_login(user) {
    var save_str = JSON.stringify(user);
    localStorage.setItem("userInfo", save_str);
}

function do_login(data) {
    fetch(`${API_HOST}/api/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('网络响应不正常');
        }
        return response.json();
    })
    .then(res => {
        console.log('服务器响应登录:', res);
        if (res.code == CodeAck.E_OK) {
            on_after_login(res.data);
            alert("登录成功!");
            window.location.href = "home.html";
        } else {
            alert("登录失败! 错误码: " + res.code);
        }
    })
    .catch(error => {
            console.error('发送数据时出错:', error);
    });
}



