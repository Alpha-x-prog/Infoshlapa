<template>
<div class="app">
    <router-view></router-view>
    <div  @click="scrollToTop" class="up-button">
        <p>Вверх</p>
    </div>
</div>
</template>

<script>
export default {
    name: 'App',
    methods: {
        scrollToTop() {
            window.scrollTo(0, 0);
        }
    } 
}
</script>

<style>
/*Сброс дефолтных стилей*/
*{
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
/*Модный шрифт*/
@import url('https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Island+Moments&display=swap');
.island-moments-font{
    font-family: "Island Moments", cursive !important;
    font-weight: 400;
    font-style: normal;
}

nav {
    padding: 1rem;
}

nav a {
    text-decoration: none;
    font-size: 1.2rem;
    color: #2c3e50;
    margin: 0 0.5rem;
}
nav a:hover{
    color: #2d8a60c9;
}
nav a.router-link-active {
    color: #42b983;
}

body {
    margin: 0;
    padding: 0;
    font-family: 'Noto Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

#app {
    font-family: 'Noto Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
}
.up-button {
    position: fixed;
    text-align: center;
    background-color: #2d8a60c9;
    padding: 10px;
    border-radius: 5px;
    cursor: pointer;
    bottom: 20px;
    left: 20px;
    z-index: 1000;
}
.up-button p{
    font-size: 1.0rem;
    color: #fafafa;
    font-weight: 400;
    padding: 0;
    margin: 0;
}
</style>