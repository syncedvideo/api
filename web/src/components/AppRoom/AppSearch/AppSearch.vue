<template>
  <div>
    <div>
      <button
        @click="openSearch"
        class="flex items-center text-primary-500 hover:text-primary-700 focus:outline-none select-none leading-none"
      >
        <app-icon :value="mdiVideoPlus" class="w-6 h-6 mr-2" />
        <span>Add video</span>
      </button>
    </div>

    <div
      v-if="showSearch"
      class="fixed inset-0 w-full h-full z-50 overflow-auto"
      style="background:rgba(0,0,0,.9)"
      @click.self="closeSearch"
    >
      <div
        class="relative z-10 mx-auto top-5 w-full max-w-2xl bg-gray-900 rounded-lg overflow-hidden"
      >
        <div class="relative flex items-center">
          <app-icon
            :value="mdiMagnify"
            class="z-10 absolute left-4 w-6 h-6 pointer-events-none text-primary-500"
          />
          <input
            ref="searchInput"
            v-model.trim="search"
            @keydown.enter="searchHandler"
            class="block w-full leading-none py-4 px-5 pl-14 bg-transparent outline-none text-lg disabled:opacity-50 bg-gray-700 text-primary-500"
            placeholder="Search YouTube"
            type="text"
            :disabled="searchLoading"
          />
          <app-icon
            v-if="searchLoading"
            :value="mdiLoading"
            class="animate-spin z-10 absolute right-4 w-6 h-6 pointer-events-none text-primary-500"
          />
        </div>

        <div v-if="videoSearch" class="overflow-auto px-5 mt-5">
          <app-search-item
            v-for="video of videoSearch.videos"
            :key="video.id"
            :video="video"
            class="mb-5"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, nextTick, Ref, ref, watch } from 'vue'
import { search as searchApi, VideoSearchResponseDto } from '@/api'
import AppSearchItem from './AppSearchItem.vue'
import AppIcon from '@/components/AppIcon'
import { mdiVideoPlus, mdiLoading, mdiMagnify } from '@mdi/js'

export default defineComponent({
  components: { AppSearchItem, AppIcon },

  name: 'AppSearch',

  setup() {
    const showSearch = ref(false)
    const search = ref('')
    const searchInput: Ref<null | HTMLInputElement> = ref(null)
    const videoSearch: Ref<VideoSearchResponseDto | undefined> = ref(undefined)
    const searchLoading = ref(false)

    function openSearch() {
      search.value = ''
      showSearch.value = true
      nextTick(() => {
        if (searchInput.value) {
          searchInput.value.focus()
        }
      })
    }

    function closeSearch() {
      showSearch.value = false
      videoSearch.value = undefined
      searchLoading.value = false
    }

    async function searchHandler() {
      if (
        search.value &&
        videoSearch.value?.query !== search.value &&
        !searchLoading.value
      ) {
        searchLoading.value = true
        const response = await searchApi(search.value)
        videoSearch.value = response.data
        searchLoading.value = false
      }
    }

    watch(showSearch, showSearch => {
      const htmlElem = document.querySelector('html')
      if (!htmlElem) {
        return
      }
      if (showSearch) {
        htmlElem.classList.add('overflow-hidden')
      } else {
        htmlElem.classList.remove('overflow-hidden')
      }
    })

    return {
      searchInput,
      openSearch,
      closeSearch,
      showSearch,
      search,
      videoSearch,
      searchLoading,
      searchHandler,
      mdiVideoPlus,
      mdiLoading,
      mdiMagnify
    }
  }
})
</script>
