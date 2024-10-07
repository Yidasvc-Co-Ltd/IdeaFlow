<template>
    <div class="main">
        <div class="global-bar-top">
            <div class="global-menu">
                <div class="global-menu-item" style="font-weight:800;font-style:italic">ideaFlow</div>
                <div class="global-menu-item" v-for="menu_item in global.menu" style="cursor: pointer;">{{
                    menu_item.hint }}
                </div>
            </div>
        </div>
        <div class="global-bar-left">
            <div class="global-leftbar-item" v-for="(leftbar_item, index) in global.leftbar">
                <img :src="leftbar_item.icon" width="28px" height="28px" style="cursor: pointer;"
                    @click="show_jpg(index)" v-bind:path="index === 0 ? '文件管理器' : (index === 1 ? '服务器管理器' : '')" />
            </div>

        </div>

        <div class="global-bar-bottom">
            <div style="background-color: rgb(114, 161, 169); width: 100%;height: 100%;"></div>
        </div>
        <div class="global-context">
            <div class="document-bar-top">
                <div class="document-path">{{ document_info.path }}</div>
            </div>
            <div class="show_photo">
                <div class="operate">
                    <span title="后退">&#9668</span>
                    <!-- 这里先当作测试show_model1 的按钮 -->
                    <span title="前进" style="margin-right: 400px;" @click="show_model2">&#9658</span>
                    <!-- 这里先当作测试show_model2 的按钮 -->

                </div>
                <!-- 
                    这里是固定的一开始有一个可以变换的文档，作为rex代码，这里是现在的div上去填充
                    同时后面存在model1和model2
                    model1对应静态图片，静态图片的内容就只是一个div，里面展示图片，修改方法：v-if
                    model2对应动态图片，动态图片的内容是三个框，分别展示不同的内容，修改方法：v-if
                -->
                <!-- 这里展示图片 -->
                <!-- 如果觉得空，可以一开始插入图片，当查看自己的图片时，可以将已有图片替换掉 -->
                <div class="model_all">
                    <!-- model1  -->
                    <img :src="selectedImage" alt="Selected Image" v-if="selectedImage" />

                    <!-- model2 -->
                    <div v-if="dyn_img">
                        <div class="top" style="vertical-align: text-top">
                            <div class="save" @click="updateImage(img.name)">保存</div>
                            <div class="save">
                                <select :v-model="select_server" @change="printSelectedName($event.target.value)"
                                    class="select">
                                    <option v-for="(item, index) in hostsArray" :value="item.value"
                                        :selected="index === 0">
                                        {{ item.name }}</option>
                                </select>
                            </div>
                        </div>
                        <div class="model2_top" contenteditable @input="updateCodeGenTask" title="CodeGenTask">{{
                    procedure.codeGenTask }}
                        </div>
                        <div class="model2_top" contenteditable @input="updatecodeFig" title="CodeFig">{{
                    procedure.codeFig }}</div>
                        <div class="model2_middle" @click="runCode()">运行Py代码</div>
                        <div class="model2_bottom_left" title="实验任务">{{ procedure.tagQueue }}</div>
                        <div class="model2_bottom_right" title="图片">
                        </div>
                    </div>

                </div>

            </div>
            <div class="show_PDF">
                <div class="compile_button">
                    <button class="button" @click="Compile_PDF">重新编译</button>
                    <img src="../assets/导出.jpg" alt="" title="导出PDF" @click="download_PDF">
                </div>
                <div class="compile_PDF"><img src="../assets/PDF.jpg" alt="PDF image"></div>

            </div>
        </div>
        <div class="global-bar-middle" v-if="showMiddle1" id="middle1">
            <p
                style="font-size: 16px;font-weight: bold; text-align: center; margin: 0 0 10px 0;border-bottom: 2px solid rgb(185, 185, 185);">
                文件管理</p>
            <input type="file" ref="fileInput" style="display: none" @change="handleFileChange" />
            <button @click="triggerFileInput" style="margin-bottom: 10px; cursor: pointer;width: 50%;">上传</button>
            <button style="margin-bottom: 10px; cursor: pointer;width: 50%; " @click="creat_dyn">新建</button>

            <div v-for="(image, index) in uploadedImages" :key="index" style="background-color: rgb(131, 226, 255);">
                <div class="photo_border" @contextmenu.prevent="showImageContextMenu($event, image)">
                    <span class="image_name" :data-name="image.name" :path="image.name" :title="image.name"
                        @dblclick="displayImage(image.name)">{{ image.name }}</span>
                </div>
            </div>
            <div v-for="(image, index) in uploadedDynImages" :key="index" style="background-color: rgb(131, 226, 255);">
                <div class="photo_border" @contextmenu.prevent="showImageContextMenu($event, image)">
                    <span class="image_name" :data-name="image.name" :path="image.name" :title="image.name"
                        @dblclick="displayDynImage(image)">{{ image.name }}</span>
                </div>
            </div>
            <transition name="fade" appear>
                <div v-if="isImageContextMenuVisible" class="context-menu"
                    :style="{ top: imageContextMenuPosition.y + 'px', left: imageContextMenuPosition.x + 'px' }">
                    <ul>
                        <li @click="renameDyn()">重命名</li>
                        <li @click="deleteImage()">删除</li>
                        <li @click="addImage()">添加</li>
                    </ul>
                </div>
            </transition>
        </div>
        <div class="global-bar-middle" v-if="showMiddle2" id="middle2">
            <p
                style="font-size: 16px;font-weight: bold; text-align: center; margin: 0 0 10px 0;border-bottom: 2px solid rgb(185, 185, 185);">
                HOST</p>
            <button @click="createServer()" style="cursor: pointer; width: 100%;margin-bottom: 5px;">新增host</button>
            <dialog v-if="showForm" class="server-form-overlay">
                <div class="server-form">
                    <div class="button-container">
                        <div class="add_server">
                            添加服务器
                        </div>
                    </div>
                    <div class="adapt-form">
                        <form v-if="currentModule === 'module'" @submit.prevent="submitForm" class="form">
                            <div class="form-group">
                                <label for="name">名称</label>
                                <input type="text" id="name" v-model="formFields.name" required>
                            </div>
                            <div class="form-group">
                                <label for="address">地址</label>
                                <input type="text" id="address" v-model="formFields.address" required>
                            </div>
                            <div class="form-group">
                                <label for="port">端口</label>
                                <input type="number" id="port" v-model="formFields.port" required>
                            </div>
                            <div class="form-group">
                                <label for="note">备注</label>
                                <textarea id="note" v-model="formFields.note"></textarea>
                            </div>
                            <div class="form-group">
                                <label for="username">用户名</label>
                                <input type="text" id="username" v-model="formFields.username">
                            </div>
                            <div class="form-group">
                                <label for="password">密码</label>
                                <input type="password" id="password" v-model="formFields.password">
                            </div>
                            <div class="form-group">
                                <label for="proxy_jump">跳板机</label>
                                <input type="text" id="proxy__jump" v-model="formFields.proxy_jump">
                            </div>
                            <div class="form-actions">
                                <button type="submit" @click="create_Server()">新建</button>
                                <button type="submit" @click="update_Server()">保存</button>
                                <button type="button" @click="closeForm">取消</button>
                            </div>

                        </form>
                    </div>
                </div>
            </dialog>
            <div v-for="( hostItem, index ) in  hostsArray " :key="index">
                <div class="photo_border" :data-name="hostItem.name" @dblclick="displayHost(hostItem.name)"
                    @contextmenu.prevent="showServertMenu($event, hostItem)" :path="hostItem.name"
                    :title="hostItem.name">
                    <span class="photo_name">{{ hostItem.name }}</span>
                    <img src="../assets/设置.png" alt="" @click="Modify_configuration(hostItem)">
                </div>
            </div>
            <transition name="fade" appear>
                <div v-if="isServerMenuVisible" class="context-menu"
                    :style="{ top: imageContextMenuPosition.y + 'px', left: imageContextMenuPosition.x + 'px' }">
                    <ul>
                        <li @click="showServerForm()">修改</li>
                        <li @click="deleteServer()">删除</li>
                    </ul>
                </div>
            </transition>

        </div>
    </div>
</template>

<script>
import { createServer } from "vite";
import global from "../components/global.vue"
import axios from 'axios'
// import { VueTex } from 'vuetex'

export default {
    data() {
        return {
            global: global,
            document_info: null,
            showMiddle1: false,
            showMiddle2: false,
            uploadedDynImages: [],
            uploadedImages: [],
            selectedImage: null,
            img: null,
            dyn_img: false,
            currentDyn: '',
            hostsArray: [],
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
            isImageContextMenuVisible: false,
            isServerMenuVisible: false,
            imageContextMenuPosition: { x: 0, y: 0 },
            procedure: {
                codeGenTask: 'codeGenTask',
                codeFig: 'codeFig',
                tagQueue: '实验任务',
            },
            text: '展示图片'
        }
    },
    mounted() {
        // this.show_jpg(0);
        // const data = {
        //     "documentID": this.document_info.documentID,
        //     "userID": "qqq",
        //     "operate": "DynFigures_query_all",
        //     "operate_type": "Operate_dynFigures"
        // };
        // //这个是文件管理
        // axios.post('/api/backend', data)
        //     .then(response => {
        //         console.log('file_response.data:', response.data);
        //         console.log('文件初始化成功');
        //         this.uploadedDynImages = response.data.data
        //     })
        //     .catch(error => {
        //         console.error(error);
        //     });

        // //这个是服务器管理   后面有个refreshServer也要改
        // const data_server = {
        //     "userID": "qqq",
        //     "operate": "Servers_query_all",
        //     "operate_type": "Operate_servers"
        // };
        // axios.post('/api/backend', data_server)
        //     .then(response => {
        //         console.log('server_response.data:', response.data.data);
        //         this.hostsArray = response.data.data
        //         console.log('this.hostsArray', this.hostsArray);
        //         console.log('服务器初始化成功');
        //     })
        //     .catch(error => {
        //         console.error(error);
        //     });
    },
    methods: {
        show_jpg(index) {
            // 重置所有图片的样式
            const leftbarItems = document.querySelectorAll('.global-leftbar-item img');
            leftbarItems.forEach((item, idx) => {
                item.style.filter = 'brightness(100%)'; // 重置图片颜色
            });

            // 根据 index 改变对应图片的样式
            const targetImage = leftbarItems[index];
            targetImage.style.filter = 'brightness(150%)'; // 修改图片颜色为加深

            // 如果点击的是同一个图片，则将图片颜色恢复成原样
            if (index === this.prevIndex) {
                targetImage.style.filter = 'brightness(100%)'; // 恢复图片颜色
                this.prevIndex = -1; // 将 prevIndex 重置为 -1
            } else {
                this.prevIndex = index; // 更新 prevIndex 的值
            }

            //判断展示什么内容
            if (index === 0) {
                this.showMiddle1 = !this.showMiddle1;
                if (this.showMiddle1) {
                    this.showMiddle2 = false;
                }
            } else if (index === 1) {
                this.showMiddle2 = !this.showMiddle2;
                if (this.showMiddle2) {
                    this.showMiddle1 = false;
                }
            }

            if (this.showMiddle1 || this.showMiddle2) {
                this.updateContextPosition(131);
            } else {
                this.updateContextPosition(41);
            }
        },
        updateContextPosition(left) {
            const context = document.querySelector('.global-context');
            if (context) {
                context.style.left = left + 'px';
            }
        },
        triggerFileInput() {
            // 点击按钮时触发文件上传框
            this.$refs.fileInput.click();
        },
        handleFileChange(event) {
            const file = event.target.files[0];
            // 检查是否已存在相同的图片名称
            const existingImage = this.uploadedDynImages.find(image => image.name === file.name);
            if (existingImage) {
                // 如果存在相同的图片名称，取消上传操作并弹出提醒框
                alert('请不要提交相同的图片名称');
                return;
            }
            const imageUrl = URL.createObjectURL(file);
            // 将上传的图片对象保存到数组中
            this.uploadedImages.push({ name: file.name, url: imageUrl });
            // 清除文件输入，允许多次上传同一文件
            event.target.value = '';
        },
        displayImage(imageName) {
            //dyn_img用于展示model，currentDyn用来展示自己对应的model
            this.dyn_img = false
            this.currentDyn = ''

            // 恢复所有块的默认颜色
            console.log(imageName);
            const photoBorders = document.querySelectorAll('.image_name');
            photoBorders.forEach(border => {
                border.style.backgroundColor = 'rgb(131, 226, 255)';
            });
            const selectedImage = this.uploadedImages.find(image => image.name === imageName);
            if (this.selectedImage === selectedImage.url) {
                this.selectedImage = ''; // 清空选中的图片路径
                return; // 结束函数执行
            }
            // 根据图片名字找到对应的图片对象
            if (selectedImage) {
                this.selectedImage = selectedImage.url; // 设置选中的图片路径
                this.selectedPhotoName = imageName; // 设置选中的图片名称

                // 将选中的图片名称对应的块进行变色
                const selectedPhotoBorder = document.querySelector(`.photo_border span[data-name="${imageName}"]`);
                if (selectedPhotoBorder) {
                    selectedPhotoBorder.style.backgroundColor = '#2590d2';
                }
            }
        },
        displayDynImage(img) {
            //静态图片消失
            this.selectedImage = ''
            console.log(img);
            //判断是否双击的是同一个动态图片
            if (this.currentDyn == img) {
                this.dyn_img = !this.dyn_img
            } else {
                this.dyn_img = true
            }
            // 恢复所有块的默认颜色
            console.log(img.name);

            // 将选中的图片名称对应的块进行变色
            const selectedPhotoBorder = document.querySelector(`.photo_border span[data-name="${img.name}"]`);
            if (this.currentDyn == img) {
                if (selectedPhotoBorder.style.backgroundColor == 'rgb(131, 226, 255)') {
                    const photoBorders = document.querySelectorAll('.image_name');
                    photoBorders.forEach(border => {
                        border.style.backgroundColor = 'rgb(131, 226, 255)';
                    });
                    selectedPhotoBorder.style.backgroundColor = '#2590d2';
                }
                else {
                    const photoBorders = document.querySelectorAll('.image_name');
                    photoBorders.forEach(border => {
                        border.style.backgroundColor = 'rgb(131, 226, 255)';
                    });
                    selectedPhotoBorder.style.backgroundColor = 'rgb(131, 226, 255)'
                }
            }
            else {
                const photoBorders = document.querySelectorAll('.image_name');
                photoBorders.forEach(border => {
                    border.style.backgroundColor = 'rgb(131, 226, 255)';
                });
                selectedPhotoBorder.style.backgroundColor = '#2590d2';
            }
            this.currentDyn = img
            this.img = img

            //上面是颜色的调整，下面是请求
            const data = {
                "dynFigureID": img.dynFigureID,
                "documentID": this.document_info.documentID,
                "userID": "qqq",
                "operate": "DynFigures_query",
                "operate_type": "Operate_dynFigures"
            };
            axios.post('/api/backend', data)
                .then(response => {
                    console.log(response.data);
                    this.procedure.codeGenTask = response.data.codeGenTask
                    this.procedure.codeFig = response.data.codeFig
                    this.procedure.tagQueue = response.data.tagQueue

                })
                .catch(error => {
                    console.error(error);
                });

        },
        showServerForm() {
            this.showForm = true;

        },

        closeForm() {
            this.showForm = false;
            console.log(this.formFields);
        },
        submitForm() {
            console.log('Form submitted', this.formFields);
            this.showForm = false;
        },
        Modify_configuration(hostItem) {
            console.log(hostItem);
        },
        printSelectedName(name) {
            const selectedHost = this.hostsArray.find(item => item.name === name);
            if (selectedHost) {
                console.log(selectedHost);
                //这里发送服务器的所有请求
            }
        },
        Compile_PDF() {
            console.log('compile_PDF');
        },
        download_PDF() {
            console.log('download_PDF');
        },
        show_model2() {
            this.dyn_img = !this.dyn_img
            this.selectedImage = ''
            const photoBorders = document.querySelectorAll('.image_name');
            photoBorders.forEach(border => {
                border.style.backgroundColor = 'rgb(131, 226, 255)';
            });
        },
        Refreshdyn() {
            const data = {
                "documentID": this.document_info.documentID,
                "userID": "qqq",
                "operate": "DynFigures_query_all",
                "operate_type": "Operate_dynFigures"
            };
            axios.post('/api/backend', data)
                .then(response => {
                    console.log('response.data:', response.data);
                    console.log('初始化成功');
                    this.uploadedDynImages = response.data.data
                })
                .catch(error => {
                    console.error(error);
                });
        },
        async creat_dyn() {
            const name = prompt("请输入新任务的标题", "New Task");
            // 如果用户点击了取消按钮或者没有输入标题，则不进行任何操作
            if (name === null || name === "") {
                return;
            }

            const data = {
                "documentID": this.document_info.documentID,
                "userID": "qqq",
                "name": name,
                "operate": "DynFigures_create",
                "operate_type": "Operate_dynFigures"
            };
            try {
                await axios.post('/api/backend', data)
                await this.Refreshdyn();
            }
            catch (error) {
                console.error(error);
            }
        },
        showImageContextMenu(event, image) {
            event.preventDefault(); // 阻止默认右键菜单
            this.img = image;
            this.isImageContextMenuVisible = true;
            this.imageContextMenuPosition = { x: event.clientX, y: event.clientY };
            document.addEventListener('click', this.hideContextMenu);
            window.addEventListener('click', () => {
                this.isImageContextMenuVisible = false;
            }, { once: true });
        },
        hideContextMenu() {
            this.isImageContextMenuVisible = false;
            document.removeEventListener('click', this.hideContextMenu);
        },
        showServertMenu(event, hostItem) {
            event.preventDefault(); // 阻止默认右键菜单
            this.server = hostItem;
            this.isServerMenuVisible = true;
            this.imageContextMenuPosition = { x: event.clientX, y: event.clientY };
            document.addEventListener('click', this.hideServertMenu);
            window.addEventListener('click', () => {
                this.isServerMenuVisible = false;
            }, { once: true });
        },
        hideServerMenu() {
            this.isServerMenuVisible = false;
            document.removeEventListener('click', this.hideServerMenu);
        },

        renameDyn() {
            this.hideContextMenu();
            console.log('重命名图片');
            const name = prompt("请输入新文档的标题", this.img.name);
            if (name === null || name === "") {
                return;
            }
            this.updateImage(name)
        },
        async updateImage(name) {
            const data = {
                "dynFigureID": this.img.dynFigureID,
                "documentID": this.document_info.documentID,
                "userID": "qqq",
                "name": name,
                "currentTag": this.img.currentTag,
                "codeGenTask": this.img.codeGenTask,
                "codeFig": this.img.codeFig,
                "tagQueue": this.img.tagQueue,
                "operate": "DynFigures_update",
                "operate_type": "Operate_dynFigures"
            };
            try {
                console.log('保存成功');
                await axios.post('/api/backend', data);
                await this.Refreshdyn();
            }
            catch (error) {
                console.error(error)
            };
        },
        async deleteImage() {
            // 删除图片逻辑
            this.hideContextMenu();
            this.dyn_img = false
            const data = {
                "dynFigureID": this.img.dynFigureID,
                "documentID": this.document_info.documentID,
                "userID": "qqq",
                "operate": "DynFigures_delete",
                "operate_type": "Operate_dynFigures"
            };
            try {
                await axios.post('/api/backend', data);
                await this.Refreshdyn();
            }
            catch (error) {
                console.error(error);
            }
            console.log(this.img);
        },
        runCode() {
            console.log('codeGenTask为', this.img.codeGenTask);
            console.log('codeFig为', this.img.codeFig);
        },
        addImage() {
            console.log('将图片添加到tex中');
        },
        RefreshServer() {
            const data_server = {
                "userID": "qqq",
                "operate": "Servers_query_all",
                "operate_type": "Operate_servers"
            };
            axios.post('/api/backend', data_server)
                .then(response => {
                    console.log('response.data:', response.data);
                    this.hostsArray = response.data.data
                    console.log(this.hostsArray);
                    console.log('服务器初始化成功');
                })
                .catch(error => {
                    console.error(error);
                });
        },
        async createServer() {
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
                await this.RefreshServer();
            }
            catch (error) {
                console.error(error);
            }
        },

        async updateServer() {
            const data = {
                "name": this.formFields.name,
                "userID": "qqq",
                "port": this.formFields.port,
                "ssh_user": "",
                "ip": this.formFields.address,
                "auth_method": "",
                "passwd": "",
                "key": "",
                "jumpServerName": "",
                "jumpServerUserID": "",
                "login_command": "",
                "operate": "Servers_create",
                "operate_type": "Operate_servers"
            };
            console.log(this.formFields.name);
            try {
                await axios.post('/api/backend', data)
                console.log('cg');
                await this.RefreshServer();
            }
            catch (error) {
                console.error(error);
            }
        },
        async deleteServer() {
            const data = {
                "name": this.formFields.name,
                "userID": "qqq",
                "operate": "Servers_delete",
                "operate_type": "Operate_servers"
            };
            console.log(this.formFields.name);
            try {
                await axios.post('/api/backend', data)
                await this.RefreshServer();
            }
            catch (error) {
                console.error(error);
            }
        },
        updateCodeGenTask(event) {
            this.img.codeGenTask = event.target.innerText
        },
        updatecodeFig(event) {
            this.img.codeFig = event.target.innerText
        }





    },
    created() {
        this.document_info = this.$cookies.get("fastpaper-document-info");
        this.updateContextPosition(41);
    },
}
</script>

<style scoped lang="scss">
.main {
    width: 100%;
    height: 100vh;
    background: linear-gradient(to right, #57bfff, #feb47b);
    user-select: none;
}

.global-bar-top {
    position: fixed;
    height: 25px;
    padding: 1px 2px;
    width: 100vw;
    border-bottom: 1px solid rgb(209, 209, 209);
}

.global-menu {
    white-space: nowrap;
    overflow-x: scroll;
    overflow-y: hidden;
}

.global-menu::-webkit-scrollbar {
    display: none;
}

.global-menu-item {
    color: black;
    height: 25px;
    line-height: 25px;
    display: inline-block;
    margin: 0 1px;
    padding: 0 9px;
    font-size: 15px;
    user-select: none;
    cursor: default;
}

.global-menu-item:hover {
    background-color: #EEEEEE;
    border-radius: 5px;
}

.global-bar-left {
    position: fixed;
    top: 27px;
    width: 40px;
    height: calc(100vh - 27px - 27px);
    border-right: 1px solid rgb(209, 209, 209);
}

.global-leftbar-item {
    height: 28px;
    width: 28px;
    padding: 4px;
    margin: 2px;
    user-select: none;
    cursor: default;
}

.global-leftbar-item:hover {
    background-color: #EEEEEE;
    border-radius: 5px;
}


.global-bar-bottom {
    position: fixed;
    top: calc(100vh - 27px);
    width: 100vw;
    height: 25px;
    border: 1px 2px;
    border-top: 1px solid rgb(209, 209, 209);
    z-index: 2;
}

.global-bar-middle {
    position: fixed;
    top: 27px;
    width: 90px;
    left: 41px;
    height: calc(100vh - 27px - 27px);
    border-right: 1px solid rgb(209, 209, 209);
    overflow: hidden;
    /* 隐藏溢出的内容 */
    white-space: nowrap;
    /* 不换行 */
    text-overflow: ellipsis;
    /* 溢出时显示省略号 */
    z-index: 999;
}

.global-context {
    position: fixed;
    top: 27px;
    left: 131px;
    width: calc(100vw - 41px - 0px);
    height: calc(100vh - 27px - 27px);
    // background-color: azure;
    z-index: 1;
}

.document-bar-top {
    width: 100%;
    height: 40px;
    border-bottom: 1px solid rgb(209, 209, 209);
}

.document-path {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 80%;
    padding: 0 10px;
    line-height: 40px;
}

.show_photo {
    position: absolute;
    width: 800px;
    height: calc(100% - 40px);
    overflow: auto;
    border-right: 2px solid rgb(209, 209, 209);
    border-top: 2px solid rgb(209, 209, 209);

    img {
        width: 100%;
        max-height: 100%;
        object-fit: contain;
        opacity: 0.9;
        display: block;
        margin: auto;
    }

    .operate {
        margin-bottom: 2px;

        span {
            display: inline-block;
            width: 50px;
            margin-left: 2px;
            font-size: 25px;
            text-align: center;
            background-color: rgb(84, 167, 215);
            cursor: pointer;
            opacity: 0.6;
            border-radius: 6px;
        }
    }
}




.show_PDF {
    position: relative;
    left: 800px;
    width: 610px;
    border: rgb(209, 209, 209) 2px solid;


}


.compile_PDF {
    margin: 10px;
    border: rgb(209, 209, 209) 2px solid;
    overflow: auto;
    height: 665px;
    overflow-y: auto;

    img {
        display: block;
        width: 100%;
        opacity: 0.7;
    }
}

.compile_button {
    height: 70px;
    margin: 10px;
    border: rgb(209, 209, 209) 2px solid;

    .button {
        display: inline-block;
        border-radius: 15px;
        /* 调整按钮的边框圆角 */
        background-color: #ffcd8f;
        /* 修改背景颜色 */
        border: 3px rgb(219, 245, 255) solid;
        color: #000000;
        font-size: 16px;
        font-weight: 400;
        width: 100px;
        height: 60px;
        cursor: pointer;
        margin: 5px;
    }

    img {
        width: 60px;
        height: 60px;
        opacity: 0.6;
        cursor: pointer;
        vertical-align: middle;
        border-radius: 10px;


    }
}

span {
    user-select: none;
}

.photo_border {
    position: relative;
    border-bottom: 1px solid rgb(93, 93, 93);
    margin: 0;
    padding: 0;
    overflow: hidden;

    .image_name {
        display: inline-block;
        width: 100%;
        height: 30px;
        vertical-align: middle;
        background-color: rgb(131, 226, 255);
    }
}


.photo_border:hover {
    cursor: pointer;
}

.photo_border span {
    display: inline-block;
    width: 60px;
    height: 30px;
    vertical-align: middle;
    background-color: rgb(131, 226, 255);
}

.photo_border img {
    vertical-align: middle;
    position: relative;
    left: 0px;
    width: 30px;
    height: 100%;
    opacity: 0.7;
}

.overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
}

.form-container {
    background-color: white;
    padding: 20px;
    border-radius: 5px;
}

.server-form-overlay {
    position: fixed;
    top: 50px;
    bottom: 20px;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
}

.server-form {
    width: 100%;
    background-color: #f4f4f4;
    border-radius: 10px;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
    margin: 10px auto;

    /* 水平居中 */
    padding: 20px;
    /* 添加内边距 */
    max-width: 500px;
    max-height: 700px;
    /* 设置最大宽度 */

}

.adapt-form {
    height: 600px;
    overflow-y: auto;
}

.form-group {
    margin-bottom: 1.5rem;
}

.form-group label {
    display: block;
    margin-bottom: .4rem;
    font-weight: bold;
}

.form-group input,
.form-group textarea {
    width: calc(100% - 2rem);
    /* 减去padding的宽度 */
    height: 30px;
    padding: .2rem;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.form-group textarea {
    min-height: 100px;
}


.form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 1.5rem;
}

.form-actions button {
    padding: .8rem 1.5rem;
    margin-left: 1rem;
    font-size: 1rem;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}

.form-actions button:first-child {
    background-color: #4CAF50;
    color: #fff;
}

.form-actions button:last-child {
    background-color: #ccc;
    color: #fff;
}

.form {
    margin: 10px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.button-container {
    display: flex;
}




.active {
    background-color: #ccc;
    /* 你希望的深色背景 */
}

.add_server {
    font-size: 20px;
    text-align: center;
    font-weight: bold;
}

.model_all {
    width: 100%;
    height: calc(100vh - 105px);
    margin-top: 6px;
}

.save {
    background-color: rgb(238, 176, 61);
    width: 50%;
    vertical-align: middle;
    display: inline-block;
    border-top: 2px solid rgb(209, 209, 209);
    height: 31px;
    border-radius: 2px;
    cursor: pointer;
    border-radius: 6px;
}

.select {
    width: 100%;
    height: 100%;
    background-color: aqua;
    cursor: pointer;
    padding: 0;
    border-radius: 6px;
}

.model2_top {
    width: 49.5%;
    height: 420px;
    display: inline-block;
    background-color: skyblue;
    border: 2px solid rgb(209, 209, 209);
    margin-bottom: 2px;
}

.model2_middle {
    background-color: #ffcd8f;
    width: 100%;
    height: 50px;
    margin-bottom: 2px;
    border-top: 2px solid rgb(209, 209, 209);
    border-bottom: 2px solid rgb(209, 209, 209);
    cursor: pointer;
    border-radius: 10px;

}

.model2_bottom_left {
    display: inline-block;
    vertical-align: middle;
    height: 205px;
    width: 49%;
    background-color: aqua;
    border-right: 2px solid rgb(209, 209, 209);
    border-top: 2px solid rgb(209, 209, 209);
}

.model2_bottom_right {
    display: inline-block;
    vertical-align: middle;
    height: 205px;
    width: 49.8%;
    background-color: aquamarine;
    border-top: 2px solid rgb(209, 209, 209);
    background-image: url('../assets/图片.jpg');
    background-size: 100% 100%;
    opacity: 0.8;
}

.context-menu {
    position: fixed;
    background-color: #afaf1c;
    border: 2px solid #36a525;
    border-radius: 5px;
    box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
    z-index: 1000;
}

.context-menu ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.context-menu li {
    padding: 10px;
    cursor: pointer;
}

.context-menu li:hover {
    background-color: rgb(135, 197, 53);
}
</style>