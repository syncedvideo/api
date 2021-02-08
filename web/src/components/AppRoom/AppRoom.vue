<template>
  <div>
    <div class="fixed z-10 inset-x-0 top-0 w-full bg-gray-800 shadow-lg">
      <div class="container mx-auto">
        <div class="flex items-center justify-between h-14">
          <div class="flex items-center">
            <div>synced.video</div>
            <div class="pl-10">
              <app-search :room="state.room" />
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
        {{ state }}
        <div class="flex-1 mr-10"></div>
        <div
          class="max-w-md w-full h-full flex-shrink flex flex-col justify-between"
        >
          <app-chat />
        </div>
      </div>
      <div></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, PropType } from 'vue'
import AppSearch from './AppSearch'
import AppChat from './AppChat'
import AppRoomSettings from './AppRoomSettings'
import * as roomApi from '@/api/room'
import { RoomDto } from '@/api'
import { RoomWebSocket, state } from '.'

export default defineComponent({
  name: 'AppRoom',

  components: {
    AppSearch,
    AppChat,
    AppRoomSettings
  },

  props: {
    modelValue: {
      type: Object as PropType<RoomDto>,
      required: true
    }
  },

  setup(props) {
    onMounted(() => {
      state.room = props.modelValue
      new RoomWebSocket(roomApi.getRoomWebSocketUrl(props.modelValue.id), state)
    })
    return { state }
  }
})
</script>
