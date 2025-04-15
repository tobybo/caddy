
const API_HOST = 'http://10.23.50.159:2015'; // 替换为你的服务器地址

document.getElementById('upload-button').addEventListener('click', async () => {
    const fileInput = document.getElementById('file-input');
    const titleInput = document.getElementById('title-input');
    const descriptionInput = document.getElementById('description-input');
    const messageDiv = document.getElementById('message');

    const file = fileInput.files[0];
    const description = descriptionInput.value.trim();
    const title = titleInput.value.trim();

    // 检查文件是否存在和有效
    if (!file || !description) {
        messageDiv.textContent = '请确保选择图片文件并输入描述信息。';
        return;
    }

    // 检查文件类型
    if (!file.type.startsWith('image/')) {
        messageDiv.textContent = '只能上传图片文件。';
        return;
    }

    // 创建 FormData 对象
    const formData = new FormData();
    formData.append('image', file);
    formData.append('description', description);
    formData.append('title', title);

    try {
        const response = await fetch(`${API_HOST}/upload`, { // 替换为你的服务器上传接口
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            messageDiv.textContent = '上传成功！';
            descriptionInput.value = ''; // 清空输入框
            titleInput.value = ''; // 清空输入框
            fileInput.value = ''; // 清空文件选择
            window.location.href = "home.html";
        } else {
            messageDiv.textContent = '上传失败，请重试。';
        }
    } catch (error) {
        messageDiv.textContent = '发生错误：' + error.message;
    }
});

