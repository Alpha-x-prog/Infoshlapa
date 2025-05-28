<template>
<div class="app">
    <router-view></router-view>
</div>
</template>

<script>
export default {
    name: 'App'
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

</style>