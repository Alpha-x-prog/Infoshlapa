<template>
    <div class="navbar-container">
        <nav class="island-moments-font">
            <router-link to="/" class="island-moments-font">ИнфоShlapa</router-link> 
            <router-link to="/tgnews">Тг новости</router-link>
            <router-link to="/profile">Профиль</router-link>
        </nav>
        <b-nav pills >
            <b-nav-item-dropdown
            id="my-nav-dropdown"
            text="Категории"
            toggle-class="nav-link-custom"
            right
            >
            <b-dropdown-item :class="{ active: currentCategory === 'all' }" href="#"
                @click.prevent="setCategory('all')">Все</b-dropdown-item>

            <b-dropdown-divider></b-dropdown-divider>
            
            <b-dropdown-item :class="{ active: currentCategory === 'top' }" href="#"
                @click.prevent="setCategory('top')">топ</b-dropdown-item>

            <b-dropdown-item :class="{ active: currentCategory === 'business' }" href="#"
                @click.prevent="setCategory('business')">бизнес</b-dropdown-item>
            
            <b-dropdown-item :class="{ active: currentCategory === 'health' }" href="#"
                @click.prevent="setCategory('health')">здоровье</b-dropdown-item>

            <b-dropdown-item :class="{ active: currentCategory === 'science' }" href="#"
                @click.prevent="setCategory('science')">наука</b-dropdown-item>

            <b-dropdown-item :class="{ active: currentCategory === 'sports' }" href="#"
                @click.prevent="setCategory('sports')">спорт</b-dropdown-item>

            <b-dropdown-item :class="{ active: currentCategory === 'technology' }" href="#"
                @click.prevent="setCategory('technology')">технологии</b-dropdown-item>
            </b-nav-item-dropdown>
        </b-nav>
    </div>
</template>

<script>
export default {
    name: 'NavbarMain',
    data() {
        return {
            currentCategory: 'all',
        };
    },
    methods: {
        setCategory(category) {
            this.currentCategory = category;
            if (this.isMobile) {
                this.menuVisible = false;
            }
            this.$emit('category-changed', category);
        },
    }
};
</script>
<style>
.navbar-container {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
</style>