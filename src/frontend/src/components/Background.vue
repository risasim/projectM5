<template>
  <div class="neon-background"></div>
</template>

<script>
export default {
  name: 'NeonBackground'
}
</script>

<style scoped>
.neon-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  overflow: hidden;
  background-color: rgb(80, 4, 4);

  /* cool laser background no longer neon jus red :( */
  background-image:
    repeating-linear-gradient(50deg,transparent,transparent 80px,rgb(152, 8, 8) 80px,rgba(255, 0, 0, 0.5) 82px,transparent 82px),
    repeating-linear-gradient(60deg, transparent, transparent 80px, rgba(243, 38, 38, 0.167) 80px, rgba(255, 0, 0, 0.15) 160px),
    repeating-linear-gradient(-60deg, transparent, transparent 80px, rgba(255, 0, 4, 0.022) 80px, rgba(255, 0, 8, 0.06) 160px),
    repeating-linear-gradient(-60deg, transparent, transparent 80px, rgba(255, 0, 4, 0.104) 80px, rgba(255, 0, 8, 0.06) 160px),
    repeating-linear-gradient(20deg,transparent,transparent 80px,rgb(170, 9, 9) 80px,rgba(128, 8, 8, 0.985) 82px,transparent 82px),
    repeating-linear-gradient(-50deg,transparent,transparent 80px,rgba(255, 0, 0, 0.5) 80px,rgba(122, 4, 4, 0.896) 82px,transparent 82px),
    repeating-linear-gradient(-20deg,transparent,transparent 80px,rgba(255, 0, 0, 0.5) 80px,rgba(255, 0, 0, 0.5) 82px,transparent 82px),
    repeating-linear-gradient(20deg,transparent,transparent 80px,rgb(135, 12, 12) 80px,rgba(98, 2, 2, 0.893) 82px,transparent 82px);
  background-blend-mode: screen;
  filter: blur(3px) brightness(1.2);
}
</style>
