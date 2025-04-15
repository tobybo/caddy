
const API_HOST = 'http://10.23.50.159:2015'; // 替换为你的服务器地址

var PERMISSION_READ = 1;
var PERMISSION_WRITE = 2;

async function fetchImageList() {
    try {
        const response = await fetch(`${API_HOST}/images`);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const images = await response.json();
        console.log("images:");
        console.log(images);

        const imageContainer = document.getElementById('image-container');
        imageContainer.innerHTML = ''; // 清空容器

        images.forEach(imageInfo => {
            const imgElement = document.createElement('img');
            imgElement.className = 'image-item';
            imgElement.src = imageInfo.image_path; // 设置图片源
            imgElement.alt = imageInfo.description; // 设置描述为 alt 属性
            console.log("图片信息");
            console.log(imageInfo);
            console.log(imgElement);

            const descriptionElement = document.createElement('p');
            descriptionElement.textContent = imageInfo.description; // 设置描述文本

            const divElement = document.createElement('div');
            divElement.appendChild(imgElement);
            divElement.appendChild(descriptionElement);
            imageContainer.appendChild(divElement); // 将图片和描述添加到容器
        });

        var can_upload = false;
        var user_info = localStorage.getItem('userInfo');
        if (user_info) {
            console.log("user_info:");
            console.log(user_info);
            var data = JSON.parse(user_info);
            can_upload = (data.permission & PERMISSION_WRITE) == PERMISSION_WRITE;
        }

		// 添加上传框
		const uploadBox = document.getElementById('upload-box');
		uploadBox.addEventListener('click', () => {
			// 这里可以添加跳转到上传界面的逻辑
			//alert('进入上传图片和描述的界面');
            if (!can_upload) {
                alert("对不起，您没有上传权限!");
                return;
            }
            window.location.href = "upload.html";
			// 例如，可以使用 window.location.href 跳转到上传页面
			// window.location.href = 'upload.html';
		});

		// // 处理上下滑动
		// const galleryContainer = document.querySelector('.image-gallery-container');
		// galleryContainer.addEventListener('wheel', (event) => {
		//     event.preventDefault(); // 防止默认滚动行为
		//     galleryContainer.scrollBy({
		//         top: event.deltaY, // 根据鼠标滚动的方向和距离滚动
		//         left: 0,
		//         behavior: 'smooth' // 平滑滚动
		//     });
		// });
    } catch (error) {
        console.error('Failed to fetch images:', error);
    }
}

function set_login_btn() {
    var login_btn = document.getElementById("loginBtn");
    var user_info = localStorage.getItem('userInfo');
    if (user_info) {
        var data = JSON.parse(user_info);
        login_btn.innerText = data.name;
    } else {
        login_btn.innerText = "登录";
    }
}

function on_click_login() {
    window.location.href = "login.html";
}

// 调用函数以获取图片列表
fetchImageList();
set_login_btn();

