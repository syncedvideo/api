<template>
  <div>
    <div class="font-semibold" :style="{ color: message.user.color }">
      <span v-if="message.user.name">{{ message.user.name }}</span>
      <span v-else>{{ message.user.id }}</span>
    </div>
    <div>{{ message.text }}</div>
  </div>
</template>

<script lang="ts">
import { ChatMessageDto } from '@/api'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ChatMessage',

  props: {
    message: {
      type: Object as () => ChatMessageDto,
      required: true
    }
  },

  setup() {
    //
  }
})
</script>
