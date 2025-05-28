<template>
    <section class="categories">
        <a class="category-link island-moments-font" :class="{ active: currentCategory === 'top' }" href="#"
            @click.prevent="setCategory('top')">
        ИнфоShlapa
        </a>

        <button class="burger-button" @click="toggleMenu" v-if="isMobile">
            <i class="fas fa-bars"></i>
            77
        </button>

        <nav class="nav" :class="{ 'burger-menu': isMobile, 'show': menuVisible }">
            <a class="category-link" :class="{ active: currentCategory === 'top' }" href="#"
                @click.prevent="setCategory('top')" @click="toggleMenu">Топ</a>

            <a class="category-link" :class="{ active: currentCategory === 'politics' }" href="#"
                @click.prevent="setCategory('politics')" @click="toggleMenu">политика</a>

            <a class="category-link" :class="{ active: currentCategory === 'health' }" href="#"
                @click.prevent="setCategory('health')" @click="toggleMenu">здоровье</a>

            <a class="category-link" :class="{ active: currentCategory === 'sports' }" href="#"
                @click.prevent="setCategory('sports')" @click="toggleMenu">спорт</a>

            <a class="category-link" :class="{ active: currentCategory === 'business' }" href="#"
                @click.prevent="setCategory('business')" @click="toggleMenu">бизнес</a>

            <a class="category-link" :class="{ active: currentCategory === 'science' }" href="#"
                @click.prevent="setCategory('science')"  @click="toggleMenu">наука</a>

            <a class="category-link" :class="{ active: currentCategory === 'food' }" href="#"
                @click.prevent="setCategory('food')"  @click="toggleMenu">еда</a>
        </nav>
    </section>
</template>

<script>
export default {
    name: 'Navbar',
    data() {
        return {
            currentCategory: 'top',
            isMobile: false,
            menuVisible: false
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
        checkScreenSize() {
            this.isMobile = window.innerWidth <= 700;
            if (!this.isMobile) {
                this.menuVisible = false;
            }
        },
        toggleMenu() {
            this.menuVisible = !this.menuVisible;
        }
    },
    mounted() {
        this.checkScreenSize();
        window.addEventListener('resize', this.checkScreenSize);
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.checkScreenSize);
    }
};
</script>

<style scoped>
    .categories{
        position: sticky;
        top: 0px;
        width: 100%;
        margin-top: 10px;
        display: flex;
        justify-content: space-between;
        align-items: start;
        /*flex-wrap: wrap;*/
        background-color: rgba(243, 243, 243, 0.906);
        /*background-color: rgba(95, 151, 214, 0.815);*/
        padding: 10px;
        border-bottom-right-radius: 10px;
        border-bottom-left-radius: 10px;
        z-index: 5;
    }
    .nav{
        width: 100%;
        display: flex;
        justify-content: space-around;
        flex-wrap: wrap;
        
    }
    .categories a{
        margin: 0 10px;
        text-decoration: none;

        font-size: 1.2vw;
        font-family: Arial, Helvetica, sans-serif;
        font-weight: 500;
        color: rgb(32, 35, 37);
        position: relative;
        overflow: hidden;
    }
    .categories .active{
        color: rgb(144, 144, 144);
        text-shadow: 1px 1px 1px rgba(0, 0, 0, 0.055);
        cursor: default;
    }
    .nav a::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 50%;     
        width: 0;        
        height: 2px;
        background-color: rgb(144, 144, 144);
        transform: translateX(-50%); 
        transition: width 0.2s ease-in-out;
    }
    .nav a:hover::after,
    .categories .active::after {
        width: 100%; 
    }
    .categories a:hover{
        color: rgb(144, 144, 144);
    }

    @media(max-width: 1390px){
        .categories a{
            font-size: 1vw;
        }
    }
    @media (max-width: 700px) {
        .categories {
            flex-wrap: wrap;
            padding: 10px 20px;
        }
        
        .burger-button {
            display: block;
        }
        
        .nav {
            display: none;
            position: absolute;
            top: 100%;
            left: 0;
            flex-direction: column;
            width: 100%;
            background-color: rgba(243, 243, 243, 0.98);
            padding: 10px;
            box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
            z-index: 10;
        }

        .nav.show {
            display: flex;
        }

        .nav a {
            margin: 5px 0;
            display: block;
            text-align: center;
            padding: 12px 0;
            font-size: 16px !important;
        }

        .category-link {
            font-size: 1rem !important;
        }

        .island-moments-font {
            font-size: 1.5rem !important;
        }
    }

    .burger-button {
        display: none;
        background: none;
        border: none;
        font-size: 24px;
        cursor: pointer;
        padding: 5px 10px;
        color: rgb(32, 35, 37);
        transition: color 0.3s ease;
    }

    .burger-button:hover {
        color: rgb(144, 144, 144);
    }
</style>