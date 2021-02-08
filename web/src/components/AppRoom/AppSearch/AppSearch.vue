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
            ref="queryInputElem"
            v-model.trim="query"
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
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, nextTick, PropType, Ref, ref, watch } from 'vue'
import AppIcon from '@/components/AppIcon'
import { mdiVideoPlus, mdiLoading, mdiMagnify } from '@mdi/js'
import * as api from '@/api'
import * as roomApi from '@/api/room'

export default defineComponent({
  components: { AppIcon },

  name: 'AppSearch',

  props: {
    room: {
      type: Object as PropType<api.RoomDto>,
      required: true
    }
  },

  setup(props) {
    const showSearch = ref(false)
    const query = ref('')
    const queryInputElem: Ref<null | HTMLInputElement> = ref(null)
    const videoSearch: Ref<undefined> = ref(undefined)
    const searchLoading = ref(false)

    function openSearch() {
      query.value = ''
      showSearch.value = true
      nextTick(() => {
        if (queryInputElem.value) {
          queryInputElem.value.focus()
        }
      })
    }

    function closeSearch() {
      showSearch.value = false
      videoSearch.value = undefined
      searchLoading.value = false
    }

    async function searchHandler() {
      if (query.value && !searchLoading.value) {
        try {
          searchLoading.value = true
          const response = await roomApi.getVideoInfo(
            props.room.id,
            query.value
          )
          console.log(response.data)
        } catch (err) {
          console.error(err)
        } finally {
          searchLoading.value = false
        }
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
      queryInputElem,
      openSearch,
      closeSearch,
      showSearch,
      query,
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
