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
import { ChatMessageDto } from '@/api'
import { computed, ComputedRef, defineComponent, ref, watch } from 'vue'
import AppChatMessage from './AppChatMessage.vue'
import * as roomApi from '@/api/room'
import { state } from '..'

export default defineComponent({
  name: 'Chat',

  components: {
    AppChatMessage
  },

  setup() {
    const messages: ComputedRef<ChatMessageDto[]> = computed(() => {
      return state.chatMessages
    })

    const connections: ComputedRef<undefined> = computed(() => {
      return undefined
      // return room.state.value.data.connectionHub.connections
    })

    const currentUser: ComputedRef<undefined> = computed(() => {
      // return room.state.value.user
      return undefined
    })

    const newMessage = ref('')
    function messageHandler() {
      if (state.room && newMessage.value) {
        roomApi.sendChatMessage(state.room.id, newMessage.value)
        newMessage.value = ''
      }
    }

    watch(
      () => messages,
      () => {
        console.log('scroll down')
      }
    )

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
