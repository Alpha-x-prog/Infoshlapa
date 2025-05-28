<template>
    <div class="main-page">
        <Header/>

        <section class="main-container">
            <NavbarTG/>
            <TimeBar/>
            <AiWidget/>
        </section>

        <Footer/>
    </div>
</template>
    
<script setup>
import Header from "@/components/Header";
import Footer from "@/components/Footer";

import NavbarTG from "@/components/NavbarTG";
import TimeBar from "@/components/TimeBar";
import AiWidget from "@/components/AiWidget.vue";
</script>

<style>
.main-page {
    width: 100%;
}
.main-container{
    width: 65%;
    min-height: 85vh;
    margin: -60px auto 0 auto;
    padding: 20px 60px;
    border-top-left-radius: 15px;
    border-top-right-radius: 15px;
    background-color: rgba(243, 243, 243, 1);
    display: flex;
    flex-direction: column;
    align-items: center;
}
</style>



