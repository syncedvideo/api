<template>
  <div v-if="video" class="bg-black rounded px-3 py-2">
    <div class="player">
      <div class="pointer-events-none" ref="ytPlayerElement"></div>
    </div>

    <div>
      <progress
        ref="progressElement"
        @click="handleSeek"
        :value="currentTime"
        :max="video.duration"
        class="player-progress"
      ></progress>
    </div>

    <div class="flex items-center leading-none">
      <div>
        <button
          v-if="player.playing"
          @click="room.actions.player.pause()"
          class="focus:outline-none hover:text-white"
        >
          <app-icon :value="mdiPause" class="w-7 h-7"></app-icon>
        </button>
        <button
          v-else
          @click="room.actions.player.play()"
          class="focus:outline-none text-primary-500"
        >
          <app-icon :value="mdiPlay" class="w-7 h-7"></app-icon>
        </button>
      </div>
      <div class="mx-2">
        <button
          @click="room.actions.player.skip()"
          class="focus:outline-none hover:text-white"
        >
          <app-icon :value="mdiSkipNext" class="w-7 h-7"></app-icon>
        </button>
      </div>
      <div>
        <input v-model.number="volume" ref="volumeRangeElement" type="range" />
      </div>
      <div class="ml-3 text-sm">
        <span>{{ room.toDurationString(currentTime) }}</span>
        <span class="mx-1">/</span>
        <span>{{ room.toDurationString(video.duration) }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  onMounted,
  onUnmounted,
  Ref,
  ref,
  watch
} from 'vue'
import * as room from '../'
import YouTubePlayer from 'yt-player'
import { UserDto, VideoDto, PlayerDto } from '@/api'
import { mdiPlay, mdiPause, mdiSkipNext, mdiVolumeHigh } from '@mdi/js'
import AppIcon from '@/components/AppIcon'

export default defineComponent({
  name: 'AppPlayer',

  components: {
    AppIcon
  },

  setup() {
    const ytPlayerElement: Ref<HTMLDivElement | null> = ref(null)
    const progressElement: Ref<HTMLDivElement | null> = ref(null)
    const volumeRangeElement: Ref<HTMLDivElement | null> = ref(null)
    const currentTime = ref(0)
    const volume = ref(0)

    const player: ComputedRef<PlayerDto> = computed(() => {
      return room.state.value.data.player
    })

    const video: ComputedRef<VideoDto | undefined> = computed(() => {
      return player.value.video
    })

    const user: ComputedRef<UserDto> = computed(() => {
      return room.state.value.user
    })

    function initPlayer() {
      console.log('initPlayer')
      if (!ytPlayerElement.value || !video.value) {
        return
      }
      if (room.ytPlayer.value) {
        room.ytPlayer.value.destroy()
        room.ytPlayer.value = undefined
      }
      room.ytPlayer.value = new YouTubePlayer(ytPlayerElement.value, {
        controls: false,
        captions: false,
        related: false,
        keyboard: false,
        modestBranding: true,
        autoplay: true
      })

      room.ytPlayer.value.on('paused', () => {
        console.log('paused')
      })

      room.ytPlayer.value.on('buffering', () => {
        console.log('buffer')
        room.actions.user.setBuffering(true)
      })

      room.ytPlayer.value.on('playing', () => {
        console.log('playing')
        room.actions.user.setBuffering(false)
      })

      room.ytPlayer.value.on('timeupdate', n => {
        console.log('timeupdate', n)
        currentTime.value = Math.round(n)
      })

      room.ytPlayer.value.setVolume(0) // required for autoplay
      room.ytPlayer.value.load(
        video.value.providerId,
        player.value.playing,
        player.value.time
      )
    }

    // handle new video
    watch(video, (curr, prev) => {
      if (curr?.id !== prev?.id) {
        initPlayer()
      }
    })

    // handle play and pause
    watch(player, (curr, prev) => {
      if (!room.ytPlayer.value) {
        return
      }
      if (curr.playing && !prev.playing) {
        room.ytPlayer.value.play()
        return
      }
      if (!curr.playing && prev.playing) {
        room.ytPlayer.value.pause()
      }
    })

    onMounted(() => {
      initPlayer()
    })

    onUnmounted(() => {
      if (room.ytPlayer.value) {
        room.ytPlayer.value.destroy()
      }
    })

    function handleSeek(e: MouseEvent) {
      const element = e.target as HTMLElement
      const clickedPercent = (100 / element.offsetWidth) * e.offsetX
      const max = element.getAttribute('max')
      if (room.ytPlayer.value && max) {
        room.actions.player.seek(
          Math.round((parseFloat(max) / 100) * clickedPercent)
        )
        room.actions.player.play()
      }
    }

    watch(volume, volume => {
      room.ytPlayer.value?.setVolume(volume)
    })

    return {
      player,
      ytPlayerElement,
      progressElement,
      volumeRangeElement,
      room,
      user,
      handleSeek,
      video,
      currentTime,
      volume,
      mdiPlay,
      mdiPause,
      mdiSkipNext,
      mdiVolumeHigh
    }
  }
})
</script>

<style lang="scss">
.player {
  float: none;
  clear: both;
  width: 100%;
  position: relative;
  padding-bottom: 56.25%;
  padding-top: 25px;
  height: 0;
}
.player iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
.player-progress {
  @apply w-full;
  @apply h-1;
  @apply cursor-pointer;
  &[value]::-webkit-progress-bar {
    @apply bg-gray-700;
  }
  &[value]::-webkit-progress-value {
    @apply bg-primary-500;
  }
}
</style>
