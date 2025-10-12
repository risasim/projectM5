<script setup>
import AppHeader from './components/Header.vue'
import NeonBackground from './components/Background.vue'
</script>

<template>
  <div class="app-container">
    <!-- Background sits behind everything -->
    <NeonBackground />

    <!-- Header is visible on all pages -->
    <AppHeader />

    <!-- Page content changes depending on the route -->
    <main class="page-content">
      <router-view />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: Arial, Helvetica, sans-serif;
  overflow-x: hidden;
}

/* app-container ensures layout stability */
.app-container {
  position: relative;
  min-height: 100vh;
}

/* content area below header */
.page-content {
  display: flex;
  justify-content: center;
  align-items: center;
  padding-top: 90px; /* same height as header so it doesn't overlap */
  z-index: 1;
  position: relative;
}
</style>
