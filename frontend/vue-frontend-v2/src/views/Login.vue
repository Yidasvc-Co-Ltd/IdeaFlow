<template>
    <!-- 背景图片 -->
    <div class="background-image"></div>

    <div class="login-box">
        <!-- Login标题 -->
        <h4 class="login-title" style="font-size: 30px;
    color:#fefefe;
    text-shadow:0px 1px 0px #c0c0c0,
       0px 2px 0px #b0b0b0,
       0px 3px 0px #a0a0a0,
       0px 4px 0px #909090,
       0px 5px 10px rgba(0, 0, 0, .9);">Login</h4>

        <!-- 用户名输入框 -->
        <div class="form-group">
            <img src="../assets/账户.png" alt="" class='small_img'>
            <input type="text" id="username" v-model="formData.username" placeholder="Enter your username">
        </div>

        <!-- 密码输入框 -->
        <div class="form-group">
            <img src="../assets/密码.png" alt="" class='small_img'>
            <input type="password" id="password" v-model="formData.password" placeholder="Enter your password">
        </div>

        <!-- 登录按钮 -->
        <button @click="login" class="login-btn">Login</button>

        <!-- 注册链接 -->
        <p>Don't have an account? <a href="#">Register</a></p>
    </div>
</template>
  
<script>
export default {
    data() {
        return {
            formData: {
                username: '',
                password: '',
            }
        };
    },
    methods: {
        login() {
            console.log(this.formData);
            // this.$router.push("/home");
            window.location.href = "/home";
        }
    }
};
</script>
  
<style scoped>
.background-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image: url('../src/assets/背景图.jpg');
    /* 使用您的背景图片路径 */
    background-size: cover;
    background-position: center;
}

.login-box {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 400px;
    padding: 20px;
    border-radius: 5px;
    background-color: rgba(255, 255, 255, 0.5);
    /* 登录框背景颜色及透明度 */
}

.small_img {
    width: 37px;
    /* 调整图像的宽度 */
    display: inline-block;
    vertical-align: middle;
    padding-right: 5px;
    opacity: 0.2;
    border-radius: 2px;

}

.form-group {
    margin-bottom: 20px;
}

input[type="text"],
input[type="password"] {
    width: calc(100% - 70px);
    /* 减去图像的宽度和两侧的空白间距 */
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    background-color: rgba(255, 255, 255, 0.7);
    vertical-align: middle;
}

.login-title {
    font-family: 'Arial Black', sans-serif;
    /* 设置艺术字体 */
    text-align: center;
}

.form-group {
    margin-bottom: 20px;
}

label {
    display: block;
    margin-bottom: 5px;
}


.login-btn {
    width: 100%;
    padding: 10px;
    border: none;
    border-radius: 5px;
    background-color: #ff7474;
    color: #fff;
    cursor: pointer;
    border-radius: 10px;
}

.login-btn:hover {
    background-color: #f79825;
}

p {
    text-align: center;
}
</style>