<script setup>
import HelloWorld from './components/HelloWorld.vue'
import TheWelcome from './components/TheWelcome.vue'

function addTwoNumbersHandler(){
  let result = window.addTwoNumbers(1, 2)
  console.log("Two numbers, 1 & 2, added: ", result)
}

async function pingSlaveHandler(){
  console.log("Slave pinged...")
  let res = await window.pingSlave()
  console.log("slave's response: ", res)
}

async function pingServiceProviderHandler(){
  console.log("Ping Service Provider")
  let res = await window.pingServiceProvider()
  console.log("SP's response: ", res)
}

</script>

<template>
  <header>
    <img alt="Vue logo" class="logo" src="./assets/logo.svg" width="125" height="125" />

    <div class="wrapper">
      <HelloWorld msg="WASM Buttons" />
      <button @click="addTwoNumbersHandler"> Add 1 & 2 </button>
      <button @click="pingSlaveHandler"> Ping Slave </button>
      <button @click="pingServiceProviderHandler"> Ping Service Provider </button>
      <button > buttonTheThird </button>
    </div>
  </header>

  <main>
    <TheWelcome />
  </main>
</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
