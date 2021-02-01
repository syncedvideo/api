<template>
  <div class="container mx-auto">
    <button
      @click="createRoomHandler()"
      :disabled="state.createRoomLoading"
      class="px-5 py-3 font-semibold bg-white focus:outline-none rounded disabled:opacity-50 text-gray-500"
    >
      Create room
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from 'vue'
import router from '@/router'
import { createRoom, auth } from '@/api'

interface State {
  createRoomLoading: boolean
}

export default defineComponent({
  name: 'HomeView',

  setup() {
    const state: State = reactive({
      createRoomLoading: false
    })

    async function createRoomHandler() {
      try {
        state.createRoomLoading = true
        await auth()
        const response = await createRoom()
        router.push({ name: 'ShowRoom', params: { id: response.data.id } })
      } catch (err) {
        state.createRoomLoading = false
        throw err
      }
    }

    return { state, createRoomHandler }
  }
})
</script>
