<template>
  <router-view />
</template>

<script setup>
// App.vue 作为路由容器
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-sans);
  font-size: 16px;
  line-height: 1.5;
  color: var(--claude-text-primary);
  background-color: var(--claude-parchment);
}

#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>