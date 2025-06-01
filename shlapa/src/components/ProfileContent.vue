<template>
    <div class="profile-content">
        <div class="profile-content-info">
            <h1>Профиль</h1>
            <p>Профиль позволяет сохранять новости, которые вам интересны. Вы можете просматривать их в любое время. Также вы можете удалить новости из профиля.</p>
            <p>Вы так же получите доступ к новостям из любимых Telegram каналов.</p>
        </div>
        <div class="profile-form-login">
            <b-form @submit="onSubmit" @reset="onReset" v-if="show">
                <b-form-group
                    id="input-group-1"
                    label="Email:"
                    label-for="input-1"
                >
                    <b-form-input
                        id="input-1"
                        v-model="form.email"
                        type="email"
                        placeholder="Введите email"
                        :state="emailValidation"
                        required
                    ></b-form-input>
                </b-form-group>

                <b-form-group id="input-group-2" label="Пароль:" label-for="input-2">
                    <b-form-input
                        id="input-2"
                        v-model="form.password"
                        type="password"
                        pattern="[0-9]{5,10}"
                        placeholder="Введите пароль"
                        :state="passwordValidation"
                        required
                    ></b-form-input>
                    <b-form-text id="password-help-block">
                    Ваш пароль должен состоять из 5-10 цифр.
                    </b-form-text>
                </b-form-group>

                <b-button type="submit" variant="primary" class="profile-form-login-button" :disabled="!formValid">
                    Авторизоваться
                </b-button>
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
            },
            show: true
        }
    },
    computed: {
        emailValidation() {
            const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/
            return emailRegex.test(this.form.email)
        },
        passwordValidation() {
            return this.form.password.length >= 5 && 
                   this.form.password.length <= 10 && 
                   /^[0-9]+$/.test(this.form.password)
        },
        formValid() {
            return this.emailValidation && 
                   this.passwordValidation && 
                   this.form.email.length > 0 && 
                   this.form.password.length > 0
        }
    },
    methods: {
        onSubmit(event) {
            event.preventDefault()
            if (this.formValid) {
                alert(JSON.stringify(this.form))
            }
        },
        onReset(event) {
            event.preventDefault()
            this.form.email = ''
            this.form.password = ''
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
