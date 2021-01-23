<template>
  <div class="flex flex-col justify-between" style="height:700px">
    <div class="overflow-auto">
      <div v-for="msg of messages" :key="msg.id">
        <app-chat-message :message="msg" />
      </div>
    </div>
    <input
      v-model.trim="newMessage"
      @keydown.enter="messageHandler()"
      type="text"
      class="bg-gray-700 px-3 py-4 outline-none rounded-lg w-full"
      placeholder="Send a message"
    />
  </div>
</template>

<script lang="ts">
import { ChatMessageDto, ConnectionMap, UserDto } from '@/api'
import { computed, ComputedRef, defineComponent, ref } from 'vue'
import AppChatMessage from './AppChatMessage.vue'
import * as room from '../'

export default defineComponent({
  name: 'Chat',

  components: {
    AppChatMessage
  },

  setup() {
    const messages: ComputedRef<ChatMessageDto[]> = computed(() => {
      return room.state.value.data.chat.messages
    })

    const connections: ComputedRef<ConnectionMap> = computed(() => {
      return room.state.value.data.connectionHub.connections
    })

    const currentUser: ComputedRef<UserDto> = computed(() => {
      return room.state.value.user
    })

    const newMessage = ref('')
    function messageHandler() {
      if (newMessage.value) {
        room.sendChatMessage(newMessage.value)
        newMessage.value = ''
      }
    }

    return {
      messages,
      connections,
      currentUser,
      newMessage,
      messageHandler
    }
  }
})
</script>
