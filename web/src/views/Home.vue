<template>
  <div class="container mx-auto">
    <button
      @click="createRoom"
      :disabled="loading"
      class="px-5 py-3 font-semibold bg-white focus:outline-none rounded disabled:opacity-50 text-gray-500"
    >
      Create room
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import router from '@/router'
import * as userApi from '@/api/user'
import * as roomApi from '@/api/room'

export default defineComponent({
  name: 'HomeView',

  setup() {
    const loading = ref(false)

    async function createRoom() {
      try {
        loading.value = true
        await userApi.auth()
        const response = await roomApi.createRoom()
        router.push({ name: 'ShowRoom', params: { id: response.data.id } })
      } catch (err) {
        loading.value = false
        throw err
      }
    }

    return { loading, createRoom }
  }
})
</script>
