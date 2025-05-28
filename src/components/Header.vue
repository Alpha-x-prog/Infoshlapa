<template>
    <div class="header-top">
        <div class="logo-container">
        <!--<img src="" alt="logo">-->
        <span class="island-moments-font">ИнфоShlapa</span>
        </div>
    </div>
</template>

<script>
</script>

<style scoped>
    .header-top{
        width: 100%;
        height: 170px;
        display: flex;
        justify-content: start;
        align-items: center;
        background-color: rgba(135, 151, 158, 1);
        z-index: 0;
    }
    .logo-container{
        margin-left: 40px;
        display: flex;
        align-items: center;
        user-select: none;
    }
    .logo-container img{
        width: 60px;
        padding: 0 10px;
    }
    .logo-container span{
        font-size: 29pt;
        color: aliceblue;
    }
    @media (max-width: 700px) {
        .logo-container span{
            font-size: 19pt;
        }
    }
</style>