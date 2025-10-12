<template>
    <footer class="footer-gradient">
      <h1>{{ footerText }}</h1>
    </footer>
</template>

<script>
export default {
  name: 'AppHeader',
  computed: {
    footerText() {
      return this.$route.meta.footerText || 'PhoSho'
    }
  }
}
</script>

<style scoped>
.footer-gradient {
  width: 100%;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.15);
  font-family: 'Trebuchet MS', 'Lucida Sans', Arial;
  background: linear-gradient(to right, #ff4b4b, #ff7b7b, #ffeaea);
}

.footer-gradient h1 {
  color: white;
  font-size: 2.5rem;
  font-weight: 800;
  text-shadow: 2px 2px 2px rgba(0, 0, 0, 0.4);
}
</style>
