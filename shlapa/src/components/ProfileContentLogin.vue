<template>
    <div class="profile-content">
        <div class="profile-content-info">
            <h1>Профиль</h1>
            <p>Профиль позволяет сохранять новости, которые вам интересны. Вы можете просматривать их в любое время. Также вы можете удалить новости из профиля.</p>
            <p>Вы так же получите доступ к новостям из любимых Telegram каналов.</p>
        </div>
        <div class="profile-form-login">
            <b-alert
                v-model="showAlert"
                :variant="alertVariant"
                dismissible
            >
                {{ alertMessage }}
            </b-alert>

            <div v-if="isLoggedIn" class="logged-in-container">
                <h3>Вы вошли как: {{ userEmail }}</h3>
                <b-button variant="danger" @click="logout" class="mt-3">Выйти</b-button>
            </div>

            <b-form v-else @submit="onSubmit" @reset="onReset" v-if="show">
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

                <div class="d-flex">
                    <b-button type="submit" variant="primary" class="profile-form-login-button" :disabled="!formValid">
                        {{ isRegistering ? 'Зарегистрироваться' : 'Войти' }}
                    </b-button>
                    <b-button variant="link" @click="toggleMode">
                        {{ isRegistering ? 'Уже есть аккаунт? Войти' : 'Создать аккаунт' }}
                    </b-button>
                </div>
            </b-form>
        </div>
    </div>
</template>

<script>
import axios from '@/utils/axios'

export default {
    data() {
        return {
            form: {
                email: '',
                password: '',
            },
            show: true,
            isRegistering: false,
            showAlert: false,
            alertVariant: 'success',
            alertMessage: '',
            userEmail: ''
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
        },
        isLoggedIn() {
            return !!localStorage.getItem('token')
        }
    },
    methods: {
        async onSubmit(event) {
            event.preventDefault()
            if (this.formValid) {
                try {
                    const endpoint = this.isRegistering ? '/api/register' : '/api/login'
                    const { data } = await axios.post(endpoint, this.form)
                    
                    if (data.success) {
                        this.showAlert = true
                        this.alertVariant = 'success'
                        this.alertMessage = data.message
                        
                        if (data.token) {
                            localStorage.setItem('token', data.token)
                            localStorage.setItem('user', JSON.stringify(data.user))
                            this.userEmail = data.user.email
                            setTimeout(() => {
                                window.location.reload()
                            }, 1000)
                        }
                    } else {
                        this.showAlert = true
                        this.alertVariant = 'danger'
                        this.alertMessage = data.error || 'Произошла ошибка'
                    }
                } catch (error) {
                    this.showAlert = true
                    this.alertVariant = 'danger'
                    this.alertMessage = error.response?.data?.error || 'Ошибка подключения к серверу'
                    console.error('Error:', error.response?.data || error.message)
                }
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
        },
        toggleMode() {
            this.isRegistering = !this.isRegistering
            this.showAlert = false
        },
        logout() {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            this.userEmail = ''
            this.showAlert = true
            this.alertVariant = 'success'
            this.alertMessage = 'Вы успешно вышли из системы'
            setTimeout(() => {
                                window.location.reload()
            }, 1000)
        }
    },
    mounted() {
        const user = localStorage.getItem('user')
        if (user) {
            this.userEmail = JSON.parse(user).email
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
    align-items: start;
}
.profile-content-info {
    width: 45%;
    margin-right: 10%;
}
.profile-form-login {
    width: 55%;
}
.profile-form-login-button {
    margin-right: 10px;
}
.logged-in-container {
    text-align: center;
    padding: 20px;
    background: #f8f9fa;
    border-radius: 8px;
}
</style>
