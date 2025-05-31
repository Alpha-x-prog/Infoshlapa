<template>
    <div class="profile-content">
        <div class="profile-content-info">
            <h1>Telegram новости</h1>
            <p>Введите ссылку на telegram канал и получайте новости в реальном времени</p>
            
        </div>
        
            <b-nav-form @submit.stop.prevent="alert('Форма отправлена')">
                <b-form-input aria-label="Ввод" class="mr-1" placeholder="Введите ссылку на канал"></b-form-input>
                <b-button type="submit">Загрузить</b-button>
            </b-nav-form>
        
        </div>
            
</template>

<script setup>

</script>

<style>
.profile-content {
    width: 100%;
    margin-top: 5%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
}
.profile-content-info {
    width: 50%;
    /*margin-right: 5%;*/
}
.profile-form-login {
    width: 50%; 
}
.profile-form-login-button {
    margin-right: 10px;
}
</style>
