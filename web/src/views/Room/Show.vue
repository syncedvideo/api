<template>
  <div>
    <app-room v-model="state.room" />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive } from 'vue'
import AppRoom from '@/components/AppRoom'
import * as api from '@/api'
import { useRoute } from 'vue-router'

interface RoomState {
  loading: boolean
  room?: api.RoomDto
}

export default defineComponent({
  name: 'RoomView',

  components: {
    AppRoom
  },

  setup() {
    const state: RoomState = reactive({ loading: true, room: undefined })
    async function getRoom() {
      const route = useRoute()
      const roomId = route.params.id.toString()
      if (!roomId) {
        return
      }
      try {
        state.loading = true
        const res = await api.getRoom(roomId)
        state.room = res.data
      } catch (error) {
        console.log(error)
      } finally {
        state.loading = false
      }
    }

    onMounted(() => {
      getRoom()
    })

    return { state }
  }
})
</script>
