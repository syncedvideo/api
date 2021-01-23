<template>
  <svg class="inline-block" viewBox="0 0 24 24">
    <path :d="value" fill="currentColor"></path>
  </svg>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'AppIcon',

  props: {
    value: {
      type: String,
      required: true
    }
  },

  setup() {
    //
  }
})
</script>
