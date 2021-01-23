<template>
  <div>
    <div>Queue</div>
    <div v-if="videos.length">
      <app-queue-item v-for="video of videos" :key="video.id" :video="video" />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, ComputedRef, defineComponent } from 'vue'
import AppQueueItem from './AppQueueItem.vue'
import { VideoDto } from '@/api'
import * as room from '../'

export default defineComponent({
  name: 'Queue',

  components: {
    AppQueueItem
  },

  setup() {
    const videos: ComputedRef<VideoDto[]> = computed(() => {
      return room.state.value.data.player.queue.videos
    })
    return { videos }
  }
})
</script>
