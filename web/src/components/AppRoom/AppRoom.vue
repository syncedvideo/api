<template>
  <div>{{ modelValue }}</div>
  <div v-if="room.state.value.data">
    <div class="fixed z-10 inset-x-0 top-0 w-full bg-gray-800 shadow-lg">
      <div class="container mx-auto">
        <div class="flex items-center justify-between h-14">
          <div class="flex items-center">
            <div>synced.video</div>
            <div class="pl-10">
              <app-search />
            </div>
          </div>
          <div>
            <app-room-settings />
          </div>
        </div>
      </div>
    </div>

    <div class="container mx-auto mt-16 pt-8">
      <div class="flex">
        <div class="flex-1 mr-10">
          <app-player />
        </div>
        <div
          class="max-w-md w-full h-full flex-shrink flex flex-col justify-between"
        >
          <app-chat />
        </div>
      </div>
      <div>
        <app-queue />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted } from 'vue'
import AppPlayer from './AppPlayer'
import AppSearch from './AppSearch'
import AppQueue from './AppQueue'
import AppChat from './AppChat'
import AppRoomSettings from './AppRoomSettings'
import * as room from './'
import { RoomDto } from '@/api'

export default defineComponent({
  name: 'AppRoom',

  components: {
    AppPlayer,
    AppSearch,
    AppQueue,
    AppChat,
    AppRoomSettings
  },

  props: {
    modelValue: {
      type: Object as () => RoomDto,
      required: true
    }
  },

  setup() {
    onMounted(() => {
      room.connect()
    })
    onUnmounted(() => {
      room.disconnect()
    })
    return { room }
  }
})
</script>
