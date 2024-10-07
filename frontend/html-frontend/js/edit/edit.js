function render_edit() {
    let data = cget("edit");
    let global_bar_left = ``;
    for (let i = 0; i < data.global.leftbar.length; i++) {
        let leftbar_item = data.global.leftbar[i];
        if (leftbar_item.brighter == false) {
            global_bar_left += `
        <div class="global-leftbar-item">
            <img src="${leftbar_item.icon}" width="28px" height="28px" style="cursor: pointer;filter:brightness(100%)"
                onclick="javascript:show_jpg(${i});" path="${i === 0 ? '文件管理器' : (i === 1 ? '服务器管理器' : '')}" />
        </div>`;
        }
        else {
            global_bar_left += `
        <div class="global-leftbar-item">
            <img src="${leftbar_item.icon}" width="28px" height="28px" style="cursor: pointer;filter:brightness(150%)"
                onclick="javascript:show_jpg(${i});" path="${i === 0 ? '文件管理器' : (i === 1 ? '服务器管理器' : '')}" />
        </div>`;
        }

    }
    document.querySelector("#global-bar-left").innerHTML = global_bar_left;

    document.querySelector("#document-path").innerHTML = data.document_info.path;

    let model_all = ``;
    if (data.selectedImage) {
        model_all += `<img src="${data.selectedImage}" alt="Selected Image" />`;
    }
    if (data.dyn_img) {
        model_all += `
        <div>
            <div class="top" style="vertical-align: text-top">
                <div class="save" onclick='javascript:updateImage("${data.img.name}")'>保存</div>
                <div class="save">
                    <select value="${data.select_server}" onchange="javascript:changeSelect(event.target.value)" class="select">`
        for (let index = 0; index < data.hostsArray.length; index++) {
            let item = data.hostsArray[index];
            model_all += `<option value="${item.value}" selected="${index === 0}">${item.name}</option>`;
        }
        model_all += `</select>
                </div>
            </div>
            <div class="model2_top" contenteditable oninput="javascript:updateCodeGenTask(event)" title="CodeGenTask">
                ${data.procedure.codeGenTask}
            </div>
            <div class="model2_top" contenteditable oninput="javascript:updatecodeFig(event)" title="CodeFig">
                ${data.procedure.codeFig}
            </div>
            <div class="model2_middle" onclick="javascript:runCode()">运行Py代码</div>
            <div class="model2_bottom_left" title="实验任务">${data.procedure.tagQueue}</div>
            <div class="model2_bottom_right" title="图片">
            </div>
        </div>`;
    }
    document.querySelector("#model_all").innerHTML = model_all;

    let middle1_container = ``;
    if (data.showMiddle1) {
        middle1_container += `
        <div id="middle1" class="global-bar-middle">
            <p id="file_manage">文件管理</p>
            <input type="file" id="fileInput" style="display: none" onchange="javascript:handleFileChange(event)" />
            <button onclick="javascript:triggerFileInput()" style="margin-bottom:2px;padding:0;border:0; cursor: pointer;width:calc(50% - 1px);">上传</button><button style="margin:0;padding:0;border-left:2px solid rgb(209,209,209); cursor: pointer;width:calc(50% - 1px); " onclick="javascript:creat_dyn()">新建</button>`;
        for (let index = 0; index < data.uploadedImages.length; index++) {
            let image = data.uploadedImages[index];
            middle1_container += `
            <div style="background-color: rgb(131, 226, 255);">
                <div class="photo_border" oncontextmenu="showImageContextMenu(event, ${index}, 0)">
                    <span class="image_name" data-name="${image.name}" path="${image.name}" title="${image.name}"
                        ondblclick='javascript:displayImage("${image.name}")'>${image.name}</span>
                </div>
            </div>`;
        }
        if (data.uploadedDynImages) {
            for (let index = 0; index < data.uploadedDynImages.length; index++) {
                let image = data.uploadedDynImages[index];
                let div_color;
                if (data.DynImages_color[index]) {
                    div_color = data.origin_div_color
                }
                else {
                    div_color = data.deepem_div_color
                }
                middle1_container += div_color + `
                    <div class="photo_border" oncontextmenu="showImageContextMenu(event, ${index}, 1)">
                        <span class="image_name" data-name="${image.name}" path="${image.name}" title="${image.name}"
                            ondblclick='javascript:displayDynImage(${index});'>${image.name}</span>
                    </div>
                </div>`;
            }
        }
        middle1_container += `
            <transition name="fade" appear>`
        if (data.isImageContextMenuVisible.value) {
            console.log('show');
            middle1_container += `
            <div class="context-menu" style="{ top: ${data.imageContextMenuPosition.y}px, left: ${data.imageContextMenuPosition.x}px }">
                <ul>
                    <li onclick="javascript:renameDyn()">重命名</li>
                    <li onclick="javascript:deleteImage()">删除</li>
                    <li onclick="javascript:addImage()">添加</li>
                </ul>
            </div>`;
        }
        middle1_container += `
            </transition>
        </div>`;
    }
    document.querySelector("#middle1_container").innerHTML = middle1_container;

    let middle2_container = ``;
    if (data.showMiddle2) {
        middle2_container += `
        <div class="global-bar-middle" id="middle2">
            <p class="p_host">HOST</p>
            <button onclick="javascript:createServer();"
                style="cursor: pointer; width: 100%;margin-bottom: 2px;">新增host</button>`
        if (data.showForm) {
            middle2_container += `
            <dialog class="server-form-overlay">
                <div class="server-form">
                    <div class="button-container">
                        <div class="add_server">
                            添加服务器
                        </div>
                    </div>
                    <div class="adapt-form">`
            if (data.currentModule == "module") {
                middle2_container += `
                        <form onsubmit="javascript:submitForm(event)" class="form">
                            <div class="form-group">
                                <label for="name">名称</label>
                                <input type="text" id="name" value="${data.formFields.name}" onchange='javascript:updateFormFields("name", "${data.formFields.name}")' required>
                            </div>
                            <div class="form-group">
                                <label for="address">地址</label>
                                <input type="text" id="address" value="${data.formFields.address}" onchange='javascript:updateFormFields("address", "${data.formFields.address}")' required>
                            </div>
                            <div class="form-group">
                                <label for="port">端口</label>
                                <input type="number" id="port" value="${data.formFields.port}" onchange='javascript:updateFormFields("port", "${data.formFields.port}")' required>
                            </div>
                            <div class="form-group">
                                <label for="note">备注</label>
                                <textarea id="note" value="${data.formFields.note}" onchange='javascript:updateFormFields("note", "${data.formFields.note}")'></textarea>
                            </div>
                            <div class="form-group">
                                <label for="username">用户名</label>
                                <input type="text" id="username" value="${data.formFields.username}" onchange='javascript:updateFormFields("username", "${data.formFields.username}")'>
                            </div>
                            <div class="form-group">
                                <label for="password">密码</label>
                                <input type="password" id="password" value="${data.formFields.password}" onchange='javascript:updateFormFields("password", "${data.formFields.password}")'>
                            </div>
                            <div class="form-group">
                                <label for="proxy_jump">跳板机</label>
                                <input type="text" id="proxy__jump" value="${data.formFields.proxy_jump}" onchange='javascript:updateFormFields("proxy_jump", "${data.formFields.proxy_jump}")'>
                            </div>
                            <div class="form-actions">
                                <button type="submit" onclick="javascript:create_Server()">新建</button>
                                <button type="submit" onclick="javascript:update_Server()">保存</button>
                                <button type="button" onclick="javascript:closeForm(event)">取消</button>
                            </div>
                        </form>`;
            }
            middle2_container += `
                    </div>
                </div>
            </dialog>`;
        }

        if (data.hostsArray) {
            for (let index = 0; index < data.hostsArray.length; index++) {
                let hostItem = data.hostsArray[index];
                middle2_container += `
                <div>
                    <div class="photo_border" style="height:27px;weight:100%"  data-name="${hostItem.name}"
                        ondblclick='javascript:displayHost("${hostItem.name}");'
                        oncontextmenu='javascript:showServertMenu(event, ${JSON.stringify(hostItem)})' path="${hostItem.name}"
                        title="${hostItem.name}">
                        <span class="photo_name">${hostItem.name}</span><img src="src/设置.png" alt="" onclick='javascript:Modify_configuration(${JSON.stringify(hostItem)})'>
                    </div>
                </div>`;
            }
        }
        middle2_container += `
            <transition name="fade" appear>`
        if (data.isServerMenuVisible) {
            middle2_container += `
            <div class="context-menu"
                style="{ top: ${data.imageContextMenuPosition.y}px, left: ${data.imageContextMenuPosition.x}px }">
                <ul>
                    <li onclick="javascript:showServerForm()">修改</li>
                    <li onclick="javascript:deleteServer()">删除</li>
                </ul>
            </div>`;
        }
        middle2_container += `
            </transition>
        </div>`;
    }
    document.querySelector("#middle2_container").innerHTML = middle2_container;
}

function updateFormFields(key, value) {
    let _ = cget("edit");
    _.formFields[key] = value;
    cset("edit", _);
    render_edit();
}

function show_jpg(index) {
    // 重置所有图片的样式
    let _ = cget("edit");
    // const leftbarItems = document.querySelectorAll('.global-leftbar-item img');
    // leftbarItems.forEach((item, idx) => {
    //     item.style.filter = 'brightness(100%)'; // 重置图片颜色
    // });
    for (let i = 0; i < _.global.leftbar.length; i++) { _.global.leftbar[i].brighter = false }
    // 根据 index 改变对应图片的样式
    // const targetImage = leftbarItems[index];
    // targetImage.style.filter = 'brightness(150%)'; // 修改图片颜色为加深
    _.global.leftbar[index].brighter = true;
    // 如果点击的是同一个图片，则将图片颜色恢复成原样
    if (index === _.prevIndex) {
        // targetImage.style.filter = 'brightness(100%)'; // 恢复图片颜色
        _.global.leftbar[index].brighter = false;
        _.prevIndex = -1; // 将 prevIndex 重置为 -1
    } else {
        _.prevIndex = index; // 更新 prevIndex 的值
    }
    //判断展示什么内容
    if (index === 0) {
        _.showMiddle1 = !_.showMiddle1;
        if (_.showMiddle1) {
            _.showMiddle2 = false;
        }
    } else if (index === 1) {
        _.showMiddle2 = !_.showMiddle2;
        if (_.showMiddle2) {
            _.showMiddle1 = false;
        }
    }
    cset("edit", _);
    render_edit();
    if (_.showMiddle1 || _.showMiddle2) {
        updateContextPosition(131);
    } else {
        updateContextPosition(41);
    }
}

function updateContextPosition(left) {
    const context = document.querySelector('.global-context');
    if (context) {
        context.style.left = left + 'px';
    }
}

function triggerFileInput() {
    // 点击按钮时触发文件上传框
    document.querySelector("#fileInput").click();
}

function handleFileChange(event) {
    const file = event.target.files[0];
    let _ = cget("edit");
    // 检查是否已存在相同的图片名称
    const existingImage = _.uploadedImages.find(image => image.name === file.name);
    if (existingImage) {
        // 如果存在相同的图片名称，取消上传操作并弹出提醒框
        alert('请不要提交相同的图片名称');
        return;
    }
    const imageUrl = URL.createObjectURL(file);
    // 将上传的图片对象保存到数组中
    _.uploadedImages.push({ name: file.name, url: imageUrl });

    console.log('_.uploadedImages', _.uploadedImages);
    cset("edit", _);
    setTimeout(function () {
        render_edit();
    }, 1000);


    // 清除文件输入，允许多次上传同一文件
    event.target.value = '';
    console.log(file);
}

function displayImage(imageName) {
    //dyn_img用于展示model，currentDyn用来展示自己对应的model
    let _ = cget("edit")
    _.dyn_img = false
    _.currentDyn = ''
    // 恢复所有块的默认颜色
    console.log(imageName);
    const photoBorders = document.querySelectorAll('.image_name');
    photoBorders.forEach(border => {
        border.style.backgroundColor = 'rgb(131, 226, 255)';
    });
    const selectedImage = _.uploadedImages.find(image => image.name === imageName);
    if (_.selectedImage === selectedImage.url) {
        _.selectedImage = ''; // 清空选中的图片路径
        return; // 结束函数执行
    }
    // 根据图片名字找到对应的图片对象
    if (selectedImage) {
        _.selectedImage = selectedImage.url; // 设置选中的图片路径
        _.selectedPhotoName = imageName; // 设置选中的图片名称
        cset("edit", _);
        render_edit();
        // 将选中的图片名称对应的块进行变色
        const selectedPhotoBorder = document.querySelector(`.photo_border span[data-name="${imageName}"]`);
        if (selectedPhotoBorder) {
            selectedPhotoBorder.style.backgroundColor = '#2590d2';
        }
    }
}

function displayDynImage(img_index) {
    //静态图片消失
    let _ = cget("edit")
    _.selectedImage = ''
    let img = _.uploadedDynImages[img_index];
    //判断是否双击的是同一个动态图片
    if (_.currentDyn == img) {
        _.dyn_img = !_.dyn_img
    } else {
        _.dyn_img = true
    }
    // 恢复所有块的默认颜色 
    // 将选中的图片名称对应的块进行变色
    const selectedPhotoBorder = document.querySelector(`.photo_border span[data-name="${img.name}"]`);
    console.log('selectedPhotoBorder', selectedPhotoBorder);
    if (_.currentDyn == img) {
        if (_.DynImages_color[img_index] == 0) {
            // const photoBorders = document.querySelectorAll('.image_name');
            // photoBorders.forEach(border => {
            //     border.style.backgroundColor = 'rgb(131, 226, 255)';
            // });
            for (let i = 0; i < _.DynImages_color.length; i++) {
                _.DynImages_color[i] = 0;
            }
            // selectedPhotoBorder.style.backgroundColor = '#2590d2';
            _.DynImages_color[img_index] = 1;
        }
        else {
            // const photoBorders = document.querySelectorAll('.image_name');
            // photoBorders.forEach(border => {
            //     border.style.backgroundColor = 'rgb(131, 226, 255)';
            // });
            // selectedPhotoBorder.style.backgroundColor = 'rgb(131, 226, 255)'
            for (let i = 0; i < _.DynImages_color.length; i++) {
                _.DynImages_color[i] = 0;
            }
        }
    }
    else {
        // const photoBorders = document.querySelectorAll('.image_name');
        // photoBorders.forEach(border => {
        //     border.style.backgroundColor = 'rgb(131, 226, 255)';
        // });
        for (let i = 0; i < _.DynImages_color.length; i++) {
            _.DynImages_color[i] = 0;
        }
        // selectedPhotoBorder.style.backgroundColor = '#2590d2';
        _.DynImages_color[img_index] = 1;
    }
    console.log('recent1', _.DynImages_color);
    _.currentDyn = img
    _.img = img
    //上面是颜色的调整，下面是请求
    const data = {
        "dynFigureID": img.dynFigureID,
        "documentID": _.document_info.documentID,
        "userID": "qqq",
        "operate": "DynFigures_query",
        "operate_type": "Operate_dynFigures"
    };
    axios.post('/api/backend', data)
        .then(response => {
            console.log(response.data);
            _.procedure.codeGenTask = response.data.codeGenTask
            _.procedure.codeFig = response.data.codeFig
            _.procedure.tagQueue = response.data.tagQueue
            // console.log('recent2', _.DynImages_color);
        })
        .catch(error => {
            console.error(error);
        });
    console.log('recent3', _.DynImages_color);
    cset("edit", _);

    console.log('dyn_img1', _.dyn_img);
    let r = cget("edit")
    console.log('dyn_img2', r.dyn_img);

    render_edit();
}

function showServerForm() {
    let _ = cget("edit")
    _.showForm = true;
    cset("edit", _);
    render_edit();
}


function closeForm() {
    let _ = cget("edit")
    _.showForm = false;
    console.log(_.formFields);
    cset("edit", _);
    render_edit();
}

function submitForm(event) {
    event.preventDefault();
    let _ = cget("edit")
    console.log('Form submitted', _.formFields);
    _.showForm = false;
    cset("edit", _);
    render_edit();
}

function Modify_configuration(hostItem) {
    console.log(hostItem);
}

function changeSelect(name) {
    let _ = cget("edit");
    _.select_server = name;
    cset("edit", _);
    render_edit();
    const selectedHost = _.hostsArray.find(item => item.name === name);
    if (selectedHost) {
        console.log(selectedHost);
        //这里发送服务器的所有请求
    }
}

function Compile_PDF() {
    console.log('compile_PDF');
}

function download_PDF() {
    console.log('download_PDF');
}

function show_model2() {
    let _ = cget("edit")
    _.dyn_img = !_.dyn_img
    _.selectedImage = ''
    const photoBorders = document.querySelectorAll('.image_name');
    photoBorders.forEach(border => {
        border.style.backgroundColor = 'rgb(131, 226, 255)';
    });
    cset("edit", _);
    render_edit();
}

function Refreshdyn() {
    let _ = cget("edit")
    const data = {
        "documentID": _.document_info.documentID,
        "userID": "qqq",
        "operate": "DynFigures_query_all",
        "operate_type": "Operate_dynFigures"
    };
    axios.post('/api/backend', data)
        .then(response => {
            console.log('response.data:', response.data);
            console.log('初始化成功');
            _.uploadedDynImages = response.data.data
        })
        .catch(error => {
            console.error(error);
        });
    cset("edit", _);
    render_edit();
}

function creat_dyn() {
    let _ = cget("edit")
    const name = prompt("请输入新任务的标题", "New Task");
    // 如果用户点击了取消按钮或者没有输入标题，则不进行任何操作
    if (name === null || name === "") {
        return;
    }

    const data = {
        "documentID": _.document_info.documentID,
        "userID": "qqq",
        "name": name,
        "operate": "DynFigures_create",
        "operate_type": "Operate_dynFigures"
    };

    axios.post('/api/backend', data)
        .then((response) => {
            console.log(_.uploadedDynImages);
            render_edit();


        })

}

function showImageContextMenu(event, image_index, isDynfig) {
    //菜单
    console.log('菜单');
    event.preventDefault();
    let _ = cget("edit")
    _.img = isDynfig ? _.uploadedDynImages[image_index] : _.uploadedImages[image_index];
    _.isImageContextMenuVisible.value = 1;
    console.log('_.isImageContextMenuVisible', _.isImageContextMenuVisible);
    _.imageContextMenuPosition = { x: event.clientX, y: event.clientY };
    // document.addEventListener('click', _.hideContextMenu);
    // window.addEventListener('click', () => {
    //     _.isImageContextMenuVisible.value = false;
    //     cset("edit", _);
    //     render_edit();
    // }, { once: true });
    cset("edit", _);
    render_edit();
    let r = cget("edit")
    console.log('_.isImageContextMenuVisible', r.isImageContextMenuVisible);

}

function hideContextMenu() {
    let _ = cget("edit")
    console.log('触发');
    _.isImageContextMenuVisible = false;
    document.removeEventListener('click', _.hideContextMenu);
    cset("edit", _);
    render_edit();
}

function showServertMenu(event, hostItem) {
    let _ = cget("edit")
    event.preventDefault(); // 阻止默认右键菜单
    _.server = hostItem;
    _.isServerMenuVisible = true;
    _.imageContextMenuPosition = { x: event.clientX, y: event.clientY };
    document.addEventListener('click', _.hideServertMenu);
    window.addEventListener('click', () => {
        _.isServerMenuVisible = false;
    }, { once: true });
    cset("edit", _);
    render_edit();
}

function hideServerMenu() {
    let _ = cget("edit")
    _.isServerMenuVisible = false;
    document.removeEventListener('click', _.hideServerMenu);
    cset("edit", _);
    render_edit();
}

function renameDyn() {
    let _ = cget("edit")
    hideContextMenu();
    console.log('重命名图片');
    const name = prompt("请输入新文档的标题", _.img.name);
    if (name === null || name === "") {
        return;
    }
    updateImage(name)
}

async function updateImage(name) {
    let _ = cget("edit")
    const data = {
        "dynFigureID": _.img.dynFigureID,
        "documentID": _.document_info.documentID,
        "userID": "qqq",
        "name": name,
        "currentTag": _.img.currentTag,
        "codeGenTask": _.img.codeGenTask,
        "codeFig": _.img.codeFig,
        "tagQueue": _.img.tagQueue,
        "operate": "DynFigures_update",
        "operate_type": "Operate_dynFigures"
    };
    try {
        console.log('保存成功');
        await axios.post('/api/backend', data);
        await Refreshdyn();
    }
    catch (error) {
        console.error(error)
    };
}

async function deleteImage() {
    // 删除图片逻辑
    let _ = cget("edit")
    hideContextMenu();
    _.dyn_img = false
    const data = {
        "dynFigureID": _.img.dynFigureID,
        "documentID": _.document_info.documentID,
        "userID": "qqq",
        "operate": "DynFigures_delete",
        "operate_type": "Operate_dynFigures"
    };
    try {
        await axios.post('/api/backend', data);
        await Refreshdyn();
    }
    catch (error) {
        console.error(error);
    }
    console.log(_.img);
    cset("edit", _);
    render_edit();
}

function runCode() {
    let _ = cget("edit")
    console.log('codeGenTask为', _.img.codeGenTask);
    console.log('codeFig为', _.img.codeFig);
}

function addImage() {
    console.log('将图片添加到tex中');
}

function RefreshServer() {
    let _ = cget("edit")
    const data_server = {
        "userID": "qqq",
        "operate": "Servers_query_all",
        "operate_type": "Operate_servers"
    };
    axios.post('/api/backend', data_server)
        .then(response => {
            console.log('response.data:', response.data);
            _.hostsArray = response.data.data
            console.log(_.hostsArray);
            console.log('服务器初始化成功');
        })
        .catch(error => {
            console.error(error);
        });
    cset("edit", _);
    render_edit();
}

async function createServer() {
    const name = prompt("请输入服务器名称", "name");
    if (name === null || name === "") {
        return;
    }
    const data = {
        "userID": "qqq",
        "name": name,
        "operate": "Servers_createe",
        "operate_type": "Operate_servers"
    };
    try {
        await axios.post('/api/backend', data)
        await RefreshServer();
    }
    catch (error) {
        console.error(error);
    }
}

async function updateServer() {
    let _ = cget("edit")
    const data = {
        "name": _.formFields.name,
        "userID": "qqq",
        "port": _.formFields.port,
        "ssh_user": "",
        "ip": _.formFields.address,
        "auth_method": "",
        "passwd": "",
        "key": "",
        "jumpServerName": "",
        "jumpServerUserID": "",
        "login_command": "",
        "operate": "Servers_create",
        "operate_type": "Operate_servers"
    };
    console.log(_.formFields.name);
    try {
        await axios.post('/api/backend', data)
        console.log('cg');
        await RefreshServer();
    }
    catch (error) {
        console.error(error);
    }
}

async function deleteServer() {
    let _ = cget("edit")
    const data = {
        "name": _.formFields.name,
        "userID": "qqq",
        "operate": "Servers_delete",
        "operate_type": "Operate_servers"
    };
    console.log(_.formFields.name);
    try {
        await axios.post('/api/backend', data)
        await RefreshServer();
    }
    catch (error) {
        console.error(error);
    }
}

function updateCodeGenTask(event) {
    let _ = cget("edit")
    _.img.codeGenTask = event.target.innerText
    cset("edit", _);
    render_edit();
}

function updatecodeFig(event) {
    let _ = cget("edit")
    _.img.codeFig = event.target.innerText
    cset("edit", _);
    render_edit();
}

function created() {
    let _ = cget("edit")
    _.document_info = cget("fastpaper-document-info");
    updateContextPosition(41);
    cset("edit", _);
    render_edit();
}

function edit_mounted() {
    let document_info = cget("fastpaper-document-info");
    cset("edit", {
        global: {
            leftbar: [
                { key: "file", hint: "文件", icon: "src/favicon.ico", component: null, brighter: false },
                { key: "server", hint: "服务器管理", icon: "src/favicon.ico", component: null, brighter: false }
            ]
        },
        origin_div_color: '<div style="background-color: rgb(131, 226, 255);">',
        deepem_div_color: '<div style="background-color: #2590d2;">',
        document_info: document_info,
        showMiddle1: false,
        showMiddle2: false,
        uploadedDynImages: [],
        DynImages_color: [],
        uploadedImages: [],
        Images_color: [],
        selectedImage: null,
        img: null,
        dyn_img: false,
        currentDyn: '',
        hostsArray: [],
        hostsArray_color: [],
        showForm: false,
        currentModule: 'module',
        formFields: {
            name: '',
            address: '',
            port: null,
            note: '',
            username: '',
            password: '',
            proxy_jump: ''
        },
        select_server: '',
        server: '',
        isImageContextMenuVisible: { value: false },
        isServerMenuVisible: false,
        imageContextMenuPosition: { x: 0, y: 0 },
        procedure: {
            codeGenTask: 'codeGenTask',
            codeFig: 'codeFig',
            tagQueue: '实验任务',
        },
        text: '展示图片'
    });
    render_edit()
    show_jpg(0);
    let _ = cget("edit");
    const data = {
        "documentID": document_info.documentID,
        "userID": document_info.userID,
        "operate": "DynFigures_query_all",
        "operate_type": "Operate_dynFigures"
    };
    //这个是文件管理
    axios.post('/api/backend', data)
        .then(response => {
            let _ = cget("edit");
            _.uploadedDynImages = response.data.data
            console.log('_.uploadedDynImages', _.uploadedDynImages);
            if (_.uploadedDynImages) {
                for (let i = 0; i < _.uploadedDynImages.length; i++) {
                    _.DynImages_color.push(0)
                }
            }
            cset("edit", _);
            render_edit();
        })
        .catch(error => {
            console.error(error);
        });

    //这个是服务器管理   后面有个refreshServer也要改
    const data_server = {
        "userID": "qqq",
        "operate": "Servers_query_all",
        "operate_type": "Operate_servers"
    };
    axios.post('/api/backend', data_server)
        .then(response => {
            console.log('server_response.data:', response.data.data);
            let _ = cget("edit");
            _.hostsArray = response.data.data
            if (_.hostsArray) {
                for (let i = 0; i < _.hostsArray.length; i++) {
                    _.hostsArray_color.push(0);
                }
            }
            cset("edit", _);
            render_edit();
        })
        .catch(error => {
            console.error(error);
        });
    render_edit();
}
