<template>
  <div v-if="!loading && room">
    <app-room v-model="room" />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, Ref, ref } from 'vue'
import AppRoom from '@/components/AppRoom'
import * as roomApi from '@/api/room'
import { useRoute } from 'vue-router'
import { RoomDto } from '@/api'

export default defineComponent({
  name: 'RoomView',

  components: {
    AppRoom
  },

  setup() {
    const room: Ref<RoomDto | undefined> = ref(undefined)
    const loading = ref(false)

    async function loadRoom() {
      const route = useRoute()
      const roomId = route.params.id.toString()
      if (!roomId) {
        return
      }
      try {
        loading.value = true
        const res = await roomApi.getRoom(roomId)
        room.value = res.data
      } catch (error) {
        console.log(error)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      loadRoom()
    })

    return { loading, room }
  }
})
</script>
