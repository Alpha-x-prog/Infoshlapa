<template>
    <div class="navbar-container">
        <nav>
            <router-link to="/" class="island-moments-font">Новости</router-link> 
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