<template>
    <div class="navbar-container">
        <nav>
            <router-link to="/" class="island-moments-font">Инфоshlapa</router-link> 
            <router-link to="/profile">Профиль</router-link>
        </nav>
        <b-nav pills>
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