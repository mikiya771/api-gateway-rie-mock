<script setup>
import axios from 'axios'
import {onMounted, ref} from "vue";

const results = ref([])
onMounted( () => {
  axios.get('/api').then(res => {results.value = res.data.msg})

})

</script>

<template>
  <h1>
    <suspense>
      <template #default>{{results}}</template>
      <template #fallback>Loading...</template>
    </suspense>
  </h1>

</template>

<style scoped>
.read-the-docs {
  color: #888;
}
</style>
