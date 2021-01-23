<template>
  <div class="flex py-3 rounded">
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
        <img class="w-24 rounded" :src="video.thumbnail" :alt="video.title" />
      </div>
    </div>
    <div class="flex-1">
      <div>
        <span class="text-lg font-semibold leading-none">
          {{ video.title }}
        </span>
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
      <div class="pt-1">
        <div v-if="video.description" style="width:400px">
          <div class="truncate text-sm text-gray-400">
            {{ video.description }}
          </div>
        </div>
        <div v-else class="italic text-sm text-gray-400">
          No description
        </div>
      </div>
      <div class="pt-1 flex items-center">
        <div>
          <button
            @click="removeHandler"
            class="inline-flex items-center text-gray-400 text-sm focus:outline-none hover:text-red-500"
          >
            <app-icon class="w-4 h-4" :value="mdiDelete"></app-icon>
          </button>
        </div>
        <div class="ml-3">
          <button
            @click="voteHandler"
            class="inline-flex items-center focus:outline-none text-sm"
            :class="
              hasVoted
                ? ['text-primary-500']
                : ['text-gray-400 hover:text-primary-700']
            "
          >
            <app-icon class="w-4 h-4 mr-2" :value="mdiThumbUp"></app-icon>
            <span class="font-semibold">{{ voteCount }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { VideoDto } from '@/api'
import { computed, ComputedRef, defineComponent } from 'vue'
import { abbreviateNumber } from 'js-abbreviation-number'
import { mdiDelete, mdiThumbUp } from '@mdi/js'
import AppIcon from '@/components/AppIcon'
import * as room from '../'

export default defineComponent({
  name: 'AppQueueItem',

  components: {
    AppIcon
  },

  props: {
    video: {
      type: Object as () => VideoDto,
      required: true
    }
  },

  setup(props) {
    const durationString = computed(() => {
      return room.toDurationString(props.video.duration)
    })

    const voteCount = computed(() => {
      return Object.keys(props.video.votes).length
    })

    const hasVoted: ComputedRef<boolean> = computed(() => {
      return props.video.votes[room.state.value.user.id] !== undefined
    })

    function voteHandler() {
      room.actions.queue.vote(props.video)
    }

    function removeHandler() {
      room.actions.queue.remove(props.video)
    }

    return {
      durationString,
      abbreviateNumber,
      voteHandler,
      removeHandler,
      voteCount,
      hasVoted,

      mdiDelete,
      mdiThumbUp
    }
  }
})
</script>
