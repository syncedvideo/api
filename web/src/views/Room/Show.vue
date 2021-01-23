<template>
  <div>
    <app-room />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import AppRoom from '@/components/AppRoom'

export default defineComponent({
  name: 'RoomView',
  components: {
    AppRoom
  }
})
</script>
