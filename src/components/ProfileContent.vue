<template>
    <div class="profile-content">
        <div class="profile-content-info">
            <h1>Профиль</h1>
            <p>Профиль позволяет сохранять новости, которые вам интересны. Вы можете просматривать их в любое время. Также вы можете удалить новости из профиля.</p>
        <p>Вы так же получит доступ к новостям из любимых Telegram каналов.</p>
        </div>
        <div class="profile-form-login">
            <b-form @submit="onSubmit" @reset="onReset" v-if="show">
            <b-form-group
                id="input-group-1"
                label="Email:"
                label-for="input-1"
                description="Мы никогда не передадим вашу электронную почту кому-либо еще."
            >
                <b-form-input
                id="input-1"
                v-model="form.email"
                type="email"
                placeholder="Введите email"
                required
                ></b-form-input>
            </b-form-group>

            <b-form-group id="input-group-2" label="Пароль:" label-for="input-2">
                <b-form-input
                id="input-2"
                v-model="form.password"
                type="password"
                placeholder="Введите пароль"
                required
                ></b-form-input>
            </b-form-group>

            <b-button type="submit" variant="primary" class="profile-form-login-button">Авторизовваться</b-button>
            <b-button type="reset" variant="outline-primary">Очистить</b-button>
            </b-form>

        </div>
    </div>
</template>

<script>
export default {
        data() {
        return {
            form: {
            email: '',
            password: '',
            checked: []
            },
            show: true
        }
        },
        methods: {
        onSubmit(event) {
            event.preventDefault()
            alert(JSON.stringify(this.form))
        },
        onReset(event) {
            event.preventDefault()
            // Сбросить значения нашей формы
            this.form.email = ''
            this.form.password = ''
            this.form.checked = []
            // Уловка для сброса/очистки состояния проверки формы в собственном браузере
            this.show = false
            this.$nextTick(() => {
            this.show = true
            })
        }
        }
    }
</script>

<style>
.profile-content {
    width: 90%;
    margin-top: 5%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
}
.profile-content-info {
    width: 50%;
    margin-right: 10%;
}
.profile-form-login {
    width: 50%;
    margin-right: 5%;
    
}
.profile-form-login-button {
    margin-right: 10px;
}



</style>
