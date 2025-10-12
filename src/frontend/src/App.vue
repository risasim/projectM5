<script setup>
import AppFooter from './components/Header.vue'
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <AppFooter />   
    <main class="flex-1 flex justify-center items-center">
      <router-view />  
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: Arial, Helvetica, sans-serif;
}
</style>
