<template>
  <div
    class="flex cursor-pointer hover:opacity-100 opacity-80 rounded"
    @click="clickHandler"
  >
    <div class="flex-shrink pr-5">
      <div class="relative">
        <div class="absolute bottom-1 right-1">
          <span
            class="inline-block text-sm leading-none p-1 rounded"
            style="background:rgba(0,0,0,.9)"
          >
            {{ durationString }}
          </span>
        </div>
        <img class="w-20 rounded" :src="video.thumbnail" :alt="video.title" />
      </div>
    </div>
    <div class="flex-1">
      <div>
        <span class="font-semibold leading-none" v-html="video.title"></span>
      </div>
      <div class="flex text-sm text-gray-400">
        <div>{{ abbreviateNumber(video.statistics.viewCount) }} views</div>
        <div class="mx-3">
          {{ abbreviateNumber(video.statistics.likeCount) }} likes
        </div>
        <div>
          {{ abbreviateNumber(video.statistics.dislikeCount) }} dislikes
        </div>
      </div>
      <div v-if="video.description" class="pt-1" style="width:400px">
        <div class="truncate text-sm text-gray-400">
          {{ video.description }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { VideoDto } from '@/api'
import { computed, defineComponent } from 'vue'
import { addQueueVideo, toDurationString } from '../'
import { abbreviateNumber } from 'js-abbreviation-number'

export default defineComponent({
  name: 'AppSearchItem',

  props: {
    video: {
      type: Object as () => VideoDto,
      required: true
    }
  },

  setup(props) {
    const durationString = computed(() => {
      return toDurationString(props.video.duration)
    })

    function clickHandler() {
      addQueueVideo(props.video)
    }

    return { clickHandler, durationString, abbreviateNumber }
  }
})
</script>
