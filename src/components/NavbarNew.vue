<template>
    <div class="navbar-container">
    <b-nav tabs  >
        <b-nav-item active>Активная</b-nav-item>
        <b-nav-item>Tg новости</b-nav-item>
        <b-nav-item-dropdown
        id="my-nav-dropdown"
        text="Категории"
        toggle-class="nav-link-custom"
        right
        >
        <b-dropdown-item>Все</b-dropdown-item>
        <b-dropdown-divider></b-dropdown-divider>
        <b-dropdown-item>бизнес</b-dropdown-item>
        <b-dropdown-item>развлечения</b-dropdown-item>
        <b-dropdown-item>здоровье</b-dropdown-item>
        <b-dropdown-item>наука</b-dropdown-item>
        <b-dropdown-item>спорт</b-dropdown-item>
        <b-dropdown-item>технологии</b-dropdown-item>
        </b-nav-item-dropdown>
    </b-nav>

    <b-nav pills>
        <b-nav-form @submit.stop.prevent="alert('Форма отправлена')">
        <b-form-input aria-label="Ввод" class="mr-1"></b-form-input>
        <b-button type="submit">Хорошо</b-button>
        </b-nav-form>
    </b-nav>
    <nav>
        <router-link to="/">Инфоshlapa</router-link> |
        <router-link to="/profile">Профиль</router-link>
    </nav>


</div>
</template>

<script>
export default {
    name: 'NavbarNew'
}
</script>
<style>
.navbar-container {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
</style>